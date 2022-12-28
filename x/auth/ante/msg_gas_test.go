package ante_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *AnteTestSuite) TestMsgGas() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// keys and addresses
	_, _, addr1 := testdata.KeyEthSecp256k1TestPubAddr()
	_, _, addr2 := testdata.KeyEthSecp256k1TestPubAddr()
	_, _, addr3 := testdata.KeyEthSecp256k1TestPubAddr()
	_, _, addr4 := testdata.KeyEthSecp256k1TestPubAddr()
	addrs := []sdk.AccAddress{addr1, addr2, addr3, addr4}

	// set accounts
	for i, addr := range addrs {
		acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr)
		suite.Require().NoError(acc.SetAccountNumber(uint64(i)))
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
	}

	msgSend := bank.NewMsgSend(addr1, addr2, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100))))
	in := make([]bank.Input, 3)
	for i := 0; i < 3; i++ {
		in[i] = bank.NewInput(addrs[i], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100))))
	}
	msgMultiSend := bank.NewMsgMultiSend(in, []bank.Output{bank.NewOutput(addr4, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(300))))})

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()

	mgd := ante.NewConsumeMsgGasDecorator(suite.app.AccountKeeper, suite.app.GashubKeeper)
	antehandler := sdk.ChainAnteDecorators(mgd)

	type testCase struct {
		name        string
		msg         sdk.Msg
		expectedGas uint64
	}
	testCases := []testCase{
		{"MsgSend", msgSend, 1000000},
		{"MsgMultiSend", msgMultiSend, 2400000},
	}
	for _, tc := range testCases {
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder() // Create new txBuilder for each test

		suite.Require().NoError(suite.txBuilder.SetMsgs(tc.msg))
		suite.txBuilder.SetFeeAmount(feeAmount)
		suite.txBuilder.SetGasLimit(gasLimit)

		tx, err := suite.CreateTestTx(nil, nil, nil, suite.ctx.ChainID())
		suite.Require().NoError(err)

		gasConsumedBefore := suite.ctx.GasMeter().GasConsumed()
		_, err = antehandler(suite.ctx, tx, false)
		gasConsumedAfter := suite.ctx.GasMeter().GasConsumed()
		suite.Require().Equal(tc.expectedGas, gasConsumedAfter-gasConsumedBefore)
	}
}
