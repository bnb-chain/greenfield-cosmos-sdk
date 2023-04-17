package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

func (suite *KeeperTestSuite) TestExportGenesis() {
	ctx := suite.ctx

	suite.Require().NoError(suite.gashubKeeper.SetParams(ctx, types.DefaultParams()))
	exportGenesis := suite.gashubKeeper.ExportGenesis(ctx)

	suite.Require().Equal(types.DefaultParams().MaxTxSize, exportGenesis.Params.MaxTxSize)
	suite.Require().Equal(types.DefaultParams().MinGasPerByte, exportGenesis.Params.MinGasPerByte)
}

func (suite *KeeperTestSuite) TestInitGenesis() {
	g := types.DefaultGenesisState()
	k := suite.gashubKeeper
	k.InitGenesis(suite.ctx, g)

	// Check that the genesis state was set correctly.
	msgSendParams := k.GetMsgGasParams(suite.ctx, sdk.MsgTypeURL(&banktypes.MsgSend{}))
	suite.Require().Equal(uint64(1200), msgSendParams.GetFixedType().FixedGas)

	msgMultiSendParams := k.GetMsgGasParams(suite.ctx, sdk.MsgTypeURL(&banktypes.MsgMultiSend{}))
	suite.Require().Equal(uint64(800), msgMultiSendParams.GetMultiSendType().FixedGas)
	suite.Require().Equal(uint64(800), msgMultiSendParams.GetMultiSendType().GasPerItem)
}
