package tx

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestEIP712Handler(t *testing.T) {
	privKey, pubkey, addr := testdata.KeyEthSecp256k1TestPubAddr()
	_, feePayerPubKey, feePayerAddr := testdata.KeyEthSecp256k1TestPubAddr()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &testdata.TestMsg{})
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	txConfig := NewTxConfig(marshaler, []signingtypes.SignMode{signingtypes.SignMode_SIGN_MODE_EIP_712})
	txBuilder := txConfig.NewTxBuilder()

	chainID := "ethermint_9000"
	memo := "some test memo"
	msgs := []sdk.Msg{banktypes.NewMsgSend(addr, addr, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(1))))}
	accNum, accSeq := uint64(1), uint64(2) // Arbitrary account number/sequence

	sigData := &signingtypes.SingleSignatureData{
		SignMode: signingtypes.SignMode_SIGN_MODE_EIP_712,
	}

	sig := signingtypes.SignatureV2{
		PubKey:   pubkey,
		Data:     sigData,
		Sequence: accSeq,
	}
	feePayerSig := signingtypes.SignatureV2{
		PubKey:   feePayerPubKey,
		Data:     sigData,
		Sequence: accSeq,
	}

	fee := txtypes.Fee{Amount: sdk.NewCoins(sdk.NewInt64Coin("atom", 150)), GasLimit: 20000}
	tip := &txtypes.Tip{Amount: sdk.NewCoins(sdk.NewInt64Coin("tip-token", 10))}

	err := txBuilder.SetMsgs(msgs...)
	require.NoError(t, err)
	txBuilder.SetMemo(memo)
	txBuilder.SetFeeAmount(fee.Amount)
	txBuilder.SetFeePayer(feePayerAddr)
	txBuilder.SetGasLimit(fee.GasLimit)
	txBuilder.SetTip(tip)

	err = txBuilder.SetSignatures(sig, feePayerSig)
	require.NoError(t, err)

	signingData := signing.SignerData{
		Address:       addr.String(),
		ChainID:       chainID,
		AccountNumber: accNum,
		Sequence:      accSeq,
		PubKey:        pubkey,
	}

	modeHandler := signModeEip712Handler{}

	t.Log("verify invalid chain ID")
	_, err = modeHandler.GetSignBytes(signingtypes.SignMode_SIGN_MODE_EIP_712, signingData, txBuilder.GetTx())
	require.EqualError(t, err, fmt.Sprintf("failed to parse chainID: %s", signingData.ChainID))

	t.Log("verify GetSignBytes correct")
	signingData.ChainID = "ethermint_9000-1"
	signBytes, err := modeHandler.GetSignBytes(signingtypes.SignMode_SIGN_MODE_EIP_712, signingData, txBuilder.GetTx())
	require.NoError(t, err)
	require.NotNil(t, signBytes)

	t.Log("verify that setting signature doesn't change sign bytes")
	expectedSignBytes := signBytes
	sigData.Signature, err = privKey.Sign(signBytes)
	require.NoError(t, err)
	err = txBuilder.SetSignatures(sig)
	require.NoError(t, err)
	signBytes, err = modeHandler.GetSignBytes(signingtypes.SignMode_SIGN_MODE_EIP_712, signingData, txBuilder.GetTx())
	require.NoError(t, err)
	require.Equal(t, expectedSignBytes, signBytes)
}

func TestEIP712Handler_DefaultMode(t *testing.T) {
	handler := signModeEip712Handler{}
	require.Equal(t, signingtypes.SignMode_SIGN_MODE_EIP_712, handler.DefaultMode())
}

func TestEIP712ModeHandler_nonDIRECT_MODE(t *testing.T) {
	invalidModes := []signingtypes.SignMode{
		signingtypes.SignMode_SIGN_MODE_DIRECT,
		signingtypes.SignMode_SIGN_MODE_TEXTUAL,
		signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		signingtypes.SignMode_SIGN_MODE_UNSPECIFIED,
	}
	for _, invalidMode := range invalidModes {
		t.Run(invalidMode.String(), func(t *testing.T) {
			var dh signModeEip712Handler
			var signingData signing.SignerData
			_, err := dh.GetSignBytes(invalidMode, signingData, nil)
			require.Error(t, err)
			wantErr := fmt.Errorf("expected %s, got %s", signingtypes.SignMode_SIGN_MODE_EIP_712, invalidMode)
			require.Equal(t, err, wantErr)
		})
	}
}

func TestEIP712ModeHandler_nonProtoTx(t *testing.T) {
	var dh signModeEip712Handler
	var signingData signing.SignerData
	tx := new(nonProtoTx)
	_, err := dh.GetSignBytes(signingtypes.SignMode_SIGN_MODE_EIP_712, signingData, tx)
	require.Error(t, err)
	wantErr := fmt.Errorf("can only handle a protobuf Tx, got %T", tx)
	require.Equal(t, err, wantErr)
}