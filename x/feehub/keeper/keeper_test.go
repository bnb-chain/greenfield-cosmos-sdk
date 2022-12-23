package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/feehub/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *IntegrationTestSuite) SetupTest() {
	app := simapp.Setup(suite.T(), false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.FeehubKeeper.SetParams(ctx, types.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.FeehubKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = queryClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
