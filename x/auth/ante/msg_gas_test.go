package ante_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	gashubtypes "github.com/cosmos/cosmos-sdk/x/gashub/types"
	"github.com/golang/mock/gomock"
)

func TestMsgGas(t *testing.T) {
	type testCase struct {
		name        string
		malleate    func(*AnteTestSuite) sdk.Msg
		expectedGas uint64
	}
	testCases := []testCase{
		{
			"Fixed gas type",
			func(suite *AnteTestSuite) sdk.Msg {
				accs := suite.CreateTestAccounts(2)

				msg := bank.NewMsgSend(accs[0].acc.GetAddress(), accs[1].acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100))))

				typeUrl := sdk.MsgTypeURL(msg)
				msgSendGasParams := gashubtypes.MsgGasParams{
					MsgTypeUrl: typeUrl,
					GasParams:  &gashubtypes.MsgGasParams_FixedType{FixedType: &gashubtypes.MsgGasParams_FixedGasParams{FixedGas: 1200}},
				}
				suite.gashubKeeper.EXPECT().GetMsgGasParams(gomock.Any(), typeUrl).Return(msgSendGasParams)

				return msg
			},
			1200,
		},
		{
			"Dynamic gas type",
			func(suite *AnteTestSuite) sdk.Msg {
				accs := suite.CreateTestAccounts(4)

				msg := bank.NewMsgMultiSend(
					bank.NewInput(accs[0].acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(300)))),
					[]bank.Output{
						bank.NewOutput(accs[1].acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100)))),
						bank.NewOutput(accs[2].acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100)))),
						bank.NewOutput(accs[3].acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100)))),
					},
				)

				typeUrl := sdk.MsgTypeURL(msg)
				msgSendGasParams := gashubtypes.MsgGasParams{
					MsgTypeUrl: typeUrl,
					GasParams:  &gashubtypes.MsgGasParams_MultiSendType{MultiSendType: &gashubtypes.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
				}
				suite.gashubKeeper.EXPECT().GetMsgGasParams(gomock.Any(), typeUrl).Return(msgSendGasParams)
				return msg
			},
			3200,
		},
	}
	for _, tc := range testCases {
		suite := SetupTestSuite(t, true)
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
		suite.ctx = suite.ctx.WithBlockHeight(1)

		suite.gashubKeeper.EXPECT().GetParams(gomock.Any())

		require.NoError(t, suite.txBuilder.SetMsgs(tc.malleate(suite)))

		tx, err := suite.CreateTestTx(suite.ctx, nil, nil, nil, suite.ctx.ChainID(), signing.SignMode_SIGN_MODE_DIRECT)
		require.NoError(t, err)

		mgd := ante.NewConsumeMsgGasDecorator(suite.accountKeeper, suite.gashubKeeper)
		anteHandler := sdk.ChainAnteDecorators(mgd)

		gasConsumedBefore := suite.ctx.GasMeter().GasConsumed()
		_, err = anteHandler(suite.ctx, tx, true)
		require.NoError(t, err)

		gasConsumedAfter := suite.ctx.GasMeter().GasConsumed()
		require.Equal(t, tc.expectedGas, gasConsumedAfter-gasConsumedBefore)
	}
}
