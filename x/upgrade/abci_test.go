package upgrade_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type TestSuite struct {
	module  module.BeginBlockAppModule
	keeper  keeper.Keeper
	querier sdk.Querier
	ctx     sdk.Context
}

var s TestSuite

func setupTest(t *testing.T, height int64) TestSuite {
	db := dbm.NewMemDB()
	app := simapp.NewSimappWithCustomOptions(t, false, simapp.SetupOptions{
		Logger:         log.NewNopLogger(),
		DB:             db,
		InvCheckPeriod: 0,
		HomePath:       simapp.DefaultNodeHome,
		EncConfig:      simapp.MakeTestEncodingConfig(),
		AppOpts:        simapp.EmptyAppOptions{},
	})

	s.keeper = app.UpgradeKeeper
	s.ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height, Time: time.Now()})

	s.module = upgrade.NewAppModule(s.keeper)
	s.querier = s.module.LegacyQuerierHandler(app.LegacyAmino())
	return s
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

	req := abci.RequestBeginBlock{Header: newCtx.BlockHeader()}
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx, req)
	})

	t.Log("Verify that the upgrade can be successfully applied with a handler")
	s.keeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return vm, nil
	})
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx, req)
	})

	VerifyCleared(t, newCtx)
}

func VerifyDoUpgradeWithCtx(t *testing.T, newCtx sdk.Context, proposalName string) {
	t.Log("Verify that a panic happens at the upgrade height")
	req := abci.RequestBeginBlock{Header: newCtx.BlockHeader()}

	t.Log("Verify that the upgrade can be successfully applied with a handler")
	s.keeper.SetUpgradeHandler(proposalName, func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return vm, nil
	})
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx, req)
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
	req := abci.RequestBeginBlock{Header: newCtx.BlockHeader()}
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx, req)
	})
	require.Equal(t, 0, called)

	t.Log("Verify we panic if we have a registered handler ahead of time")
	err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "future", Height: s.ctx.BlockHeight() + 3})
	require.NoError(t, err)
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx, req)
	})
	require.Equal(t, 0, called)

	t.Log("Verify we no longer panic if the plan is on time")

	futCtx := s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 3).WithBlockTime(time.Now())
	req = abci.RequestBeginBlock{Header: futCtx.BlockHeader()}
	require.NotPanics(t, func() {
		s.module.BeginBlock(futCtx, req)
	})
	require.Equal(t, 1, called)

	VerifyCleared(t, futCtx)
}

func VerifyCleared(t *testing.T, newCtx sdk.Context) {
	t.Log("Verify that the upgrade plan has been cleared")
	bz, err := s.querier(newCtx, []string{types.QueryCurrent}, abci.RequestQuery{})
	require.NoError(t, err)
	require.Nil(t, bz, string(bz))
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
	require.Equal(t, `Upgrade Plan
  Name: test
  height: 100
  Info: .`, types.Plan{Name: "test", Height: 100, Info: ""}.String())

	require.Equal(t, fmt.Sprintf(`Upgrade Plan
  Name: test
  height: 100
  Info: .`), types.Plan{Name: "test", Height: 100, Info: ""}.String())
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
	req := abci.RequestBeginBlock{Header: newCtx.BlockHeader()}
	err := s.keeper.ScheduleUpgrade(s.ctx, types.Plan{Name: "test", Height: s.ctx.BlockHeight() + 1})
	s.keeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan types.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return vm, nil
	})
	s.keeper.SetUpgradeInitializer("test", func() error { return nil })
	require.NoError(t, err)
	t.Log("Verify if upgrade happens without skip upgrade")
	require.NotPanics(t, func() {
		s.module.BeginBlock(newCtx, req)
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
		preRun      func() (sdk.Context, abci.RequestBeginBlock)
		expectPanic bool
	}{
		{
			"test not panic: no scheduled upgrade or applied upgrade is present",
			func() (sdk.Context, abci.RequestBeginBlock) {
				req := abci.RequestBeginBlock{Header: s.ctx.BlockHeader()}
				return s.ctx, req
			},
			false,
		},
		{
			"test not panic: upgrade handler is present for last applied upgrade",
			func() (sdk.Context, abci.RequestBeginBlock) {
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

				req := abci.RequestBeginBlock{Header: newCtx.BlockHeader()}
				return newCtx, req
			},
			false,
		},
	}

	for _, tc := range testCases {
		ctx, req := tc.preRun()

		require.NotPanics(t, func() {
			s.module.BeginBlock(ctx, req)
		})
	}
}
