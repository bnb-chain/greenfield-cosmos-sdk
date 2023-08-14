package keeper_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crosschain/keeper"
	testutil2 "github.com/cosmos/cosmos-sdk/x/crosschain/testutil"
	"github.com/cosmos/cosmos-sdk/x/crosschain/types"
	govtestutil "github.com/cosmos/cosmos-sdk/x/gov/testutil"
	"github.com/cosmos/cosmos-sdk/x/mint"
)

type TestSuite struct {
	suite.Suite

	crossChainKeeper keeper.Keeper

	ctx         sdk.Context
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (s *TestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx
	// gomock initializations
	ctrl := gomock.NewController(s.T())
	bankKeeper := govtestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := govtestutil.NewMockStakingKeeper(ctrl)
	s.crossChainKeeper = keeper.NewKeeper(
		encCfg.Codec,
		key,
		authtypes.NewModuleAddress(types.ModuleName).String(),
		stakingKeeper,
		bankKeeper,
	)

	err := s.crossChainKeeper.SetParams(s.ctx, types.DefaultParams())
	s.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(testCtx.Ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.crossChainKeeper)

	s.queryClient = types.NewQueryClient(queryHelper)
	s.msgServer = keeper.NewMsgServerImpl(s.crossChainKeeper)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestIncrSendSequence() {
	beforeSequence := s.crossChainKeeper.GetSendSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	s.crossChainKeeper.IncrSendSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	afterSequence := s.crossChainKeeper.GetSendSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	s.Require().EqualValues(afterSequence, beforeSequence+1)
}

func (s *TestSuite) TestIncrReceiveSequence() {
	beforeSequence := s.crossChainKeeper.GetReceiveSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	s.crossChainKeeper.IncrReceiveSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	afterSequence := s.crossChainKeeper.GetReceiveSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	s.Require().EqualValues(afterSequence, beforeSequence+1)
}

func (s *TestSuite) TestRegisterChannel() {
	testChannelName := "test channel"
	testChannelId := sdk.ChannelID(100)

	err := s.crossChainKeeper.RegisterChannel(testChannelName, testChannelId, &testutil2.MockCrossChainApplication{})

	s.Require().NoError(err)

	app := s.crossChainKeeper.GetCrossChainApp(testChannelId)
	s.Require().NotNil(app)

	// check duplicate name
	err = s.crossChainKeeper.RegisterChannel(testChannelName, testChannelId, app)
	s.Require().ErrorContains(err, "duplicated channel name")

	// check duplicate channel id
	err = s.crossChainKeeper.RegisterChannel("another channel", testChannelId, app)
	s.Require().ErrorContains(err, "duplicated channel id")

	// check nil app
	err = s.crossChainKeeper.RegisterChannel("another channel", sdk.ChannelID(101), nil)
	s.Require().ErrorContains(err, "nil cross chain app")
}

func (s *TestSuite) TestSetChannelSendPermission() {
	s.crossChainKeeper.SetChannelSendPermission(s.ctx, sdk.ChainID(1), sdk.ChannelID(1), sdk.ChannelAllow)

	permission := s.crossChainKeeper.GetChannelSendPermission(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))
	s.Require().EqualValues(sdk.ChannelAllow, permission)
}

func (s *TestSuite) TestUpdateChannelPermission() {
	s.crossChainKeeper.RegisterChannel("test", 1, &testutil2.MockCrossChainApplication{})
	s.crossChainKeeper.SetDestBscChainID(1)

	s.crossChainKeeper.SetChannelSendPermission(s.ctx, sdk.ChainID(1), sdk.ChannelID(1), sdk.ChannelAllow)

	permissions := []*types.ChannelPermission{
		&types.ChannelPermission{
			DestChainId: 1,
			ChannelId:   1,
			Permission:  0,
		},
	}
	err := s.crossChainKeeper.UpdatePermissions(s.ctx, permissions)
	s.Require().NoError(err)

	permission := s.crossChainKeeper.GetChannelSendPermission(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))
	s.Require().EqualValues(sdk.ChannelForbidden, permission)
}
