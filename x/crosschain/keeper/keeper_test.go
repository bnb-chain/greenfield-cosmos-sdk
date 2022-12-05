package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crosschain/testutil"
	"github.com/cosmos/cosmos-sdk/x/crosschain/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

type TestSuite struct {
	suite.Suite

	app *simapp.SimApp
	ctx sdk.Context
}

func (s *TestSuite) SetupTest() {
	app := simapp.Setup(s.T(), false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})

	app.CrossChainKeeper.SetParams(ctx, types.DefaultParams())

	s.app = app
	s.ctx = ctx
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestIncrSendSequence() {
	beforeSequence := s.app.CrossChainKeeper.GetSendSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	s.app.CrossChainKeeper.IncrSendSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	afterSequence := s.app.CrossChainKeeper.GetSendSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	s.Require().EqualValues(afterSequence, beforeSequence+1)
}

func (s *TestSuite) TestIncrReceiveSequence() {
	beforeSequence := s.app.CrossChainKeeper.GetReceiveSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	s.app.CrossChainKeeper.IncrReceiveSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	afterSequence := s.app.CrossChainKeeper.GetReceiveSequence(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))

	s.Require().EqualValues(afterSequence, beforeSequence+1)
}

func (s *TestSuite) TestRegisterDestChain() {
	testChainId := sdk.ChainID(100)

	err := s.app.CrossChainKeeper.RegisterDestChain(testChainId)

	s.Require().NoError(err)

	// check duplicate channel id
	err = s.app.CrossChainKeeper.RegisterDestChain(testChainId)
	s.Require().ErrorContains(err, "duplicated destination chain chainID")
}

func (s *TestSuite) TestRegisterChannel() {
	testChannelName := "test channel"
	testChannelId := sdk.ChannelID(100)

	err := s.app.CrossChainKeeper.RegisterChannel(testChannelName, testChannelId, &testutil.MockCrossChainApplication{})

	s.Require().NoError(err)

	app := s.app.CrossChainKeeper.GetCrossChainApp(testChannelId)
	s.Require().NotNil(app)

	// check duplicate name
	err = s.app.CrossChainKeeper.RegisterChannel(testChannelName, testChannelId, app)
	s.Require().ErrorContains(err, "duplicated channel name")

	// check duplicate channel id
	err = s.app.CrossChainKeeper.RegisterChannel("another channel", testChannelId, app)
	s.Require().ErrorContains(err, "duplicated channel id")

	// check nil app
	err = s.app.CrossChainKeeper.RegisterChannel("another channel", sdk.ChannelID(101), nil)
	s.Require().ErrorContains(err, "nil cross chain app")
}

func (s *TestSuite) TestCreateIBCPackage() {
	sequence, err := s.app.CrossChainKeeper.CreateRawIBCPackage(s.ctx, sdk.ChainID(1), sdk.ChannelID(1), sdk.CrossChainPackageType(1), []byte("test payload"))
	s.Require().NoError(err)
	s.Require().EqualValues(0, sequence)

	_, err = s.app.CrossChainKeeper.GetIBCPackage(s.ctx, sdk.ChainID(1), sdk.ChannelID(1), sequence)
	s.Require().NoError(err)
}

func (s *TestSuite) TestSetChannelSendPermission() {
	s.app.CrossChainKeeper.SetChannelSendPermission(s.ctx, sdk.ChainID(1), sdk.ChannelID(1), sdk.ChannelAllow)

	permission := s.app.CrossChainKeeper.GetChannelSendPermission(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))
	s.Require().EqualValues(sdk.ChannelAllow, permission)
}
