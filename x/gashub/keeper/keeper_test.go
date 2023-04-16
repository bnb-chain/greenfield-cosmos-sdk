package keeper_test

import (
	"fmt"
	"strings"
	"testing"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/stretchr/testify/suite"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/keeper"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx          sdk.Context
	gashubKeeper keeper.Keeper

	queryClient types.QueryClient
	msgServer   types.MsgServer

	encCfg moduletestutil.TestEncodingConfig
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	suite.ctx = ctx
	suite.gashubKeeper = keeper.NewKeeper(
		encCfg.Codec,
		key,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	types.RegisterInterfaces(encCfg.InterfaceRegistry)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, suite.gashubKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.msgServer = keeper.NewMsgServerImpl(suite.gashubKeeper)
	suite.encCfg = encCfg
}

func (suite *KeeperTestSuite) TestGetAuthority() {
	NewKeeperWithAuthority := func(authority string) keeper.Keeper {
		return keeper.NewKeeper(
			moduletestutil.MakeTestEncodingConfig().Codec,
			storetypes.NewKVStoreKey(types.StoreKey),
			authority,
		)
	}

	tests := map[string]string{
		"some random account":    "0x319D057ce294319bA1fa5487134608727e1F3e29",
		"gov module account":     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		"another module account": authtypes.NewModuleAddress(minttypes.ModuleName).String(),
	}

	for name, expected := range tests {
		suite.T().Run(name, func(t *testing.T) {
			kpr := NewKeeperWithAuthority(expected)
			actual := kpr.GetAuthority()
			suite.Require().Equal(expected, actual)
		})
	}
}

func (suite *KeeperTestSuite) TestDefaultParams() {
	ctx := suite.ctx
	require := suite.Require()
	params := types.DefaultParams()
	require.Equal(uint64(64*1024), params.MaxTxSize)
	require.Equal(uint64(5), params.MinGasPerByte)

	require.NoError(suite.gashubKeeper.SetParams(ctx, params))
}

func (suite *KeeperTestSuite) TestSetMsgGasParams() {
	ctx, gashubKeeper := suite.ctx, suite.gashubKeeper
	require := suite.Require()

	tests := []struct {
		name       string
		msgTypeUrl string
		params     types.MsgGasParams
	}{
		{
			name:       "fixed type",
			msgTypeUrl: "fixed",
			params: types.MsgGasParams{
				MsgTypeUrl: "fixed",
				GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
			},
		},
		{
			name:       "dynamic type",
			msgTypeUrl: "dynamic",
			params: types.MsgGasParams{
				MsgTypeUrl: "dynamic",
				GasParams:  &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
			},
		},
	}
	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			gashubKeeper.SetMsgGasParams(ctx, tc.params)
			actual := gashubKeeper.GetMsgGasParams(ctx, tc.msgTypeUrl)
			require.Equal(tc.params, actual)
		})
	}
}

func (suite *KeeperTestSuite) TestSetAllMsgGasParams() {
	ctx, gashubKeeper := suite.ctx, suite.gashubKeeper
	require := suite.Require()

	tests := []struct {
		name            string
		msgGasParamsSet []*types.MsgGasParams
	}{
		{
			name:            "nil",
			msgGasParamsSet: nil,
		},
		{
			name:            "empty",
			msgGasParamsSet: []*types.MsgGasParams{{}},
		},
		{
			name: "one case",
			msgGasParamsSet: []*types.MsgGasParams{
				{
					MsgTypeUrl: "fixed",
					GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
				},
			},
		},
		{
			name: "two cases",
			msgGasParamsSet: []*types.MsgGasParams{
				{
					MsgTypeUrl: "fixed",
					GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
				},
				{
					MsgTypeUrl: "dynamic",
					GasParams:  &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
				},
			},
		},
	}

	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			gashubKeeper.SetAllMsgGasParams(ctx, tc.msgGasParamsSet)

			for _, mgp := range tc.msgGasParamsSet {
				actual := gashubKeeper.GetMsgGasParams(ctx, mgp.MsgTypeUrl)
				require.Equal(*mgp, actual)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestDeleteMsgGasParams() {
	ctx, gashubKeeper := suite.ctx, suite.gashubKeeper
	require := suite.Require()

	tests := []struct {
		name       string
		msgTypeUrl string
		params     types.MsgGasParams
	}{
		{
			name:       "fixed type",
			msgTypeUrl: "fixed",
			params: types.MsgGasParams{
				MsgTypeUrl: "fixed",
				GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
			},
		},
		{
			name:       "dynamic type",
			msgTypeUrl: "dynamic",
			params: types.MsgGasParams{
				MsgTypeUrl: "dynamic",
				GasParams:  &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
			},
		},
	}
	for _, tc := range tests {
		suite.T().Run(tc.name, func(t *testing.T) {
			// set
			gashubKeeper.SetMsgGasParams(ctx, tc.params)
			actual := gashubKeeper.GetMsgGasParams(ctx, tc.msgTypeUrl)
			require.Equal(tc.params, actual)

			// delete
			gashubKeeper.DeleteMsgGasParams(ctx, tc.params.MsgTypeUrl)
			actual = gashubKeeper.GetMsgGasParams(ctx, tc.msgTypeUrl)
			require.Equal(types.MsgGasParams{}, actual)
		})
	}
}

func (suite *KeeperTestSuite) TestIterateMsgGasParamsEntries() {
	ctx, gashubKeeper := suite.ctx, suite.gashubKeeper
	require := suite.Require()

	suite.T().Run("no entries to iterate", func(t *testing.T) {
		count := 0
		gashubKeeper.IterateMsgGasParams(ctx, func(_ string, _ *types.MsgGasParams) (stop bool) {
			count++
			return false
		})

		require.Equal(0, count)
	})

	fixed := types.MsgGasParams{
		GasParams: &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
	}
	dynamic := types.MsgGasParams{
		GasParams: &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
	}
	alpha := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	urls := make([]string, len(alpha)*2)
	for i, l := range alpha {
		urls[i*2] = fmt.Sprintf("%siterfixed", l)
		urls[i*2+1] = fmt.Sprintf("%siterdynamic", l)
		fixed.MsgTypeUrl = urls[i*2]
		dynamic.MsgTypeUrl = urls[i*2+1]
		gashubKeeper.SetMsgGasParams(ctx, fixed)
		gashubKeeper.SetMsgGasParams(ctx, dynamic)
	}

	var seen []string
	suite.T().Run("all urls have expected values", func(t *testing.T) {
		gashubKeeper.IterateMsgGasParams(ctx, func(msgTypeUrl string, mgp *types.MsgGasParams) (stop bool) {
			seen = append(seen, msgTypeUrl)
			if strings.HasSuffix(msgTypeUrl, "fixed") {
				fixed.MsgTypeUrl = msgTypeUrl
				require.Equal(fixed, *mgp)
			} else {
				dynamic.MsgTypeUrl = msgTypeUrl
				require.Equal(dynamic, *mgp)
			}
			return false
		})
	})

	suite.T().Run("all urls were seen", func(t *testing.T) {
		require.ElementsMatch(urls, seen)
	})

	gashubKeeper.DeleteMsgGasParams(ctx, urls...)

	suite.T().Run("no entries to iterate again after deleting all of them", func(t *testing.T) {
		count := 0
		gashubKeeper.IterateMsgGasParams(ctx, func(_ string, _ *types.MsgGasParams) (stop bool) {
			count++
			return false
		})

		require.Equal(0, count)
	})
}

func (suite *KeeperTestSuite) TestGetAllMsgGasParamsEntries() {
	ctx, gashubKeeper := suite.ctx, suite.gashubKeeper
	require := suite.Require()

	suite.T().Run("no entries", func(t *testing.T) {
		actual := gashubKeeper.GetAllMsgGasParams(ctx)
		require.Len(actual, 0)
	})

	fixed := types.MsgGasParams{
		GasParams: &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
	}
	dynamic := types.MsgGasParams{
		GasParams: &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
	}
	alpha := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	urls := make([]string, len(alpha)*2)
	for i, l := range alpha {
		urls[i*2] = fmt.Sprintf("%siterfixed", l)
		urls[i*2+1] = fmt.Sprintf("%siterdynamic", l)
		fixed.MsgTypeUrl = urls[i*2]
		dynamic.MsgTypeUrl = urls[i*2+1]
		gashubKeeper.SetMsgGasParams(ctx, fixed)
		gashubKeeper.SetMsgGasParams(ctx, dynamic)
	}

	var seen []string
	suite.T().Run("all urls have expected values", func(t *testing.T) {
		actual := gashubKeeper.GetAllMsgGasParams(ctx)
		for _, mgp := range actual {
			seen = append(seen, mgp.MsgTypeUrl)
			if strings.HasSuffix(mgp.MsgTypeUrl, "fixed") {
				fixed.MsgTypeUrl = mgp.MsgTypeUrl
				require.Equal(fixed, *mgp)
			} else {
				dynamic.MsgTypeUrl = mgp.MsgTypeUrl
				require.Equal(dynamic, *mgp)
			}
		}
	})

	suite.T().Run("all urls were seen", func(t *testing.T) {
		require.ElementsMatch(urls, seen)
	})

	for _, url := range urls {
		gashubKeeper.DeleteMsgGasParams(ctx, url)
	}

	suite.T().Run("no entries again after deleting all of them", func(t *testing.T) {
		actual := gashubKeeper.GetAllMsgGasParams(ctx)
		require.Len(actual, 0)
	})
}

func (suite *KeeperTestSuite) TestSetParams() {
	ctx, gashubKeeper := suite.ctx, suite.gashubKeeper
	require := suite.Require()

	expected := []uint64{1234, 5}
	params := types.NewParams(expected[0], expected[1])
	require.NoError(gashubKeeper.SetParams(ctx, params))

	suite.Run("stored params are as expected", func() {
		actual := gashubKeeper.GetParams(ctx)
		require.Equal(actual.MaxTxSize, expected[0])
		require.Equal(actual.MinGasPerByte, expected[1])
	})
}
