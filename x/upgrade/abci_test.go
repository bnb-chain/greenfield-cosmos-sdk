package upgrade_test

import (
	"os"
	"testing"
	"time"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/log"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cosmos/cosmos-sdk/core/appmodule"

	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

type TestSuite struct {
	suite.Suite

	module  appmodule.HasBeginBlocker
	keeper  *keeper.Keeper
	handler govtypesv1beta1.Handler
	ctx     sdk.Context
	baseApp *baseapp.BaseApp
	encCfg  moduletestutil.TestEncodingConfig
}

var s TestSuite

func setupTest(t *testing.T, height int64) *TestSuite {
	s.encCfg = moduletestutil.MakeTestEncodingConfig(upgrade.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))

	s.baseApp = baseapp.NewBaseApp(
		"upgrade",
		log.NewNopLogger(),
		testCtx.DB,
		s.encCfg.TxConfig.TxDecoder(),
	)

	s.keeper = keeper.NewKeeper(key, s.encCfg.Codec, t.TempDir(), nil)
	s.keeper.SetVersionSetter(s.baseApp)

	s.ctx = testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: time.Now(), Height: height})

	s.module = upgrade.NewAppModule(s.keeper)
	return &s
}

func TestRequireName(t *testing.T) {
	s := setupTest(t, 10)
	err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{})
	require.Error(t, err)
}

func TestRequireFutureBlock(t *testing.T) {
	s := setupTest(t, 10)
	err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "test", Height: s.ctx.BlockHeight() - 1})
	require.Error(t, err)
}
func TestDoHeightUpgrade(t *testing.T) {
	s := setupTest(t, 10)
	t.Log("Verify can schedule an upgrade")

	err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "test", Height: s.ctx.BlockHeight() + 1})
	require.NoError(t, err)

	VerifyDoUpgrade(t)
}

func TestCanOverwriteScheduleUpgrade(t *testing.T) {
	s := setupTest(t, 10)
	t.Log("Can overwrite plan")
	err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "test", Height: s.ctx.BlockHeight() + 10})
	require.NoError(t, err)
	err = s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "test", Height: s.ctx.BlockHeight() + 1})
	require.NoError(t, err)

	VerifyDoUpgrade(t)
}

func VerifyDoUpgrade(t *testing.T) {
	t.Log("Verify that a panic happens at the upgrade height")
	newCtx := s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 1).WithBlockTime(time.Now())

	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx)
	})

	t.Log("Verify that the upgrade can be successfully applied with a handler")
	s.keeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return vm, nil
	})
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx)
	})

	VerifyCleared(t, newCtx)
}

func VerifyDoUpgradeWithCtx(t *testing.T, newCtx sdk.Context, proposalName string) {
	t.Log("Verify that a panic happens at the upgrade height")
	require.Panics(t, func() {
		s.module.BeginBlock(newCtx)
	})
	t.Log("Verify that the upgrade can be successfully applied with a handler")
	s.keeper.SetUpgradeHandler(proposalName, func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return vm, nil
	})
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx)
	})

	VerifyCleared(t, newCtx)
}

func TestHaltIfTooNew(t *testing.T) {
	s := setupTest(t, 10)
	t.Log("Verify that we don't panic with registered plan not in database at all")
	var called int
	s.keeper.SetUpgradeHandler("future", func(_ sdk.Context, _ types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		called++
		return vm, nil
	})
	s.keeper.SetUpgradeInitializer("future", func() error { return nil })

	newCtx := s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 1).WithBlockTime(time.Now())
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx)
	})
	require.Equal(t, 0, called)

	t.Log("Verify we panic if we have a registered handler ahead of time")
	err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "future", Height: s.ctx.BlockHeight() + 3})
	require.NoError(t, err)
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx)
	})
	require.Equal(t, 0, called)

	t.Log("Verify we no longer panic if the plan is on time")

	futCtx := s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 3).WithBlockTime(time.Now())
	require.NotPanics(t, func() {
		s.module.BeginBlock(futCtx)
	})
	require.Equal(t, 1, called)

	VerifyCleared(t, futCtx)
}

func VerifyCleared(t *testing.T, newCtx sdk.Context) {
	t.Log("Verify that the upgrade plan has been cleared")
	plan, _ := s.keeper.GetUpgradePlan(newCtx)
	var expected []*types.Plan
	require.Equal(t, plan, expected)
}

func TestCanClear(t *testing.T) {
	s := setupTest(t, 10)
	t.Log("Verify upgrade is scheduled")
	err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "test", Height: s.ctx.BlockHeight() + 100})
	require.NoError(t, err)

	s.keeper.ClearUpgradePlan(s.ctx)
	VerifyCleared(t, s.ctx)
}

func TestPlanStringer(t *testing.T) {
	require.Equal(t, "name:\"test\" height:100 ", (&types.Plan{Name: "test", Height: 100, Info: ""}).String())
	require.Equal(t, `name:"test" height:100 `, (&types.Plan{Name: "test", Height: 100, Info: ""}).String())
}

func VerifyNotDone(t *testing.T, newCtx sdk.Context, name string) {
	t.Log("Verify that upgrade was not done")
	height := s.keeper.GetDoneHeight(newCtx, name)
	require.Zero(t, height)
}

func VerifyDone(t *testing.T, newCtx sdk.Context, name string) {
	t.Log("Verify that the upgrade plan has been executed")
	height := s.keeper.GetDoneHeight(newCtx, name)
	require.NotZero(t, height)
}

func TestUpgrade(t *testing.T) {
	s := setupTest(t, 10)
	newCtx := s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 1).WithBlockTime(time.Now())
	err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "test", Height: s.ctx.BlockHeight() + 1})
	s.keeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return vm, nil
	})
	s.keeper.SetUpgradeInitializer("test", func() error { return nil })
	require.NoError(t, err)
	t.Log("Verify if upgrade happens without skip upgrade")
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx)
	})

	VerifyDoUpgrade(t)
	VerifyDone(t, s.ctx, "test")
}

func TestDumpUpgradeInfoToFile(t *testing.T) {
	s := setupTest(t, 10)
	require := require.New(t)

	// require no error when the upgrade info file does not exist
	_, err := s.keeper.ReadUpgradeInfoFromDisk()
	require.NoError(err)

	planHeight := s.ctx.BlockHeight() + 1
	plan := types.Plan{
		Name:   "test",
		Height: 0, // this should be overwritten by DumpUpgradeInfoToFile
	}
	t.Log("verify if upgrade height is dumped to file")
	err = s.keeper.DumpUpgradeInfoToDisk(planHeight, plan)
	require.Nil(err)

	upgradeInfo, err := s.keeper.ReadUpgradeInfoFromDisk()
	require.NoError(err)

	t.Log("Verify upgrade height from file matches ")
	require.Equal(upgradeInfo.Height, planHeight)
	require.Equal(upgradeInfo.Name, plan.Name)

	// clear the test file
	upgradeInfoFilePath, err := s.keeper.GetUpgradeInfoPath()
	require.Nil(err)
	err = os.Remove(upgradeInfoFilePath)
	require.Nil(err)
}

// TODO: add testcase to for `no upgrade handler is present for last applied upgrade`.
func TestBinaryVersion(t *testing.T) {
	s := setupTest(t, 10)

	testCases := []struct {
		name        string
		preRun      func() sdk.Context
		expectPanic bool
	}{
		{
			"test not panic: no scheduled upgrade or applied upgrade is present",
			func() sdk.Context {
				return s.ctx
			},
			false,
		},
		{
			"test not panic: upgrade handler is present for last applied upgrade",
			func() sdk.Context {
				s.keeper.SetUpgradeHandler("test0", func(_ sdk.Context, _ types.Plan, vm module.VersionMap) (module.VersionMap, error) {
					return vm, nil
				})

				err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "test0", Height: s.ctx.BlockHeight() + 2})
				require.NoError(t, err)

				newCtx := s.ctx.WithBlockHeight(12)
				s.keeper.ApplyUpgrade(newCtx, types.Plan{
					Name:   "test0",
					Height: 12,
				})

				return newCtx
			},
			false,
		},
	}

	for _, tc := range testCases {
		ctx := tc.preRun()

		require.NotPanics(t, func() {
			s.module.BeginBlock(ctx)
		})
	}
}
