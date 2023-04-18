package tx

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkmath "github.com/cosmos/cosmos-sdk/math"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/group"
)

func TestEIP712Handler(t *testing.T) {
	privKey, pubkey, addr := testdata.KeyTestPubAddrEthSecp256k1()
	_, feePayerPubKey, feePayerAddr := testdata.KeyTestPubAddrEthSecp256k1()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &banktypes.MsgSend{})
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	txConfig := NewTxConfig(marshaler, []signingtypes.SignMode{signingtypes.SignMode_SIGN_MODE_EIP_712})
	txBuilder := txConfig.NewTxBuilder()

	chainID := "greenfield_9000"
	testMemo := "some test memo"
	testMsg := banktypes.NewMsgSend(addr, addr, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(1))))
	accNum, accSeq := uint64(1), uint64(2)

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

	err := txBuilder.SetMsgs(testMsg)
	require.NoError(t, err)
	txBuilder.SetMemo(testMemo)
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
	signingData.ChainID = "greenfield_9000-1"
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

func TestMoreMsgs(t *testing.T) {
	_, pubkey, addr := testdata.KeyTestPubAddrEthSecp256k1()
	_, feePayerPubKey, feePayerAddr := testdata.KeyTestPubAddrEthSecp256k1()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &vesting.MsgCreateVestingAccount{})
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &govtypes.MsgSubmitProposal{})
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &banktypes.MsgSend{})
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &feegrant.MsgGrantAllowance{})
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &group.MsgCreateGroup{})
	interfaceRegistry.RegisterInterface(
		"cosmos.feegrant.v1beta1.FeeAllowanceI",
		(*feegrant.FeeAllowanceI)(nil),
		&feegrant.BasicAllowance{},
		&feegrant.PeriodicAllowance{},
		&feegrant.AllowedMsgAllowance{},
	)
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	txConfig := NewTxConfig(marshaler, []signingtypes.SignMode{signingtypes.SignMode_SIGN_MODE_EIP_712})
	txBuilder := txConfig.NewTxBuilder()

	chainID := "greenfield_9000-1"
	testMemo := "some test memo"
	accNum, accSeq := uint64(1), uint64(2)

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

	txBuilder.SetMemo(testMemo)
	txBuilder.SetFeeAmount(fee.Amount)
	txBuilder.SetFeePayer(feePayerAddr)
	txBuilder.SetGasLimit(fee.GasLimit)

	err := txBuilder.SetSignatures(sig, feePayerSig)
	require.NoError(t, err)

	signingData := signing.SignerData{
		Address:       addr.String(),
		ChainID:       chainID,
		AccountNumber: accNum,
		Sequence:      accSeq,
		PubKey:        pubkey,
	}
	modeHandler := signModeEip712Handler{}

	msgSend := banktypes.NewMsgSend(addr, addr, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(1))))
	msgProposal, _ := govtypes.NewMsgSubmitProposal(
		[]sdk.Msg{msgSend},
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(1))),
		"test",
		"test",
		"test",
		"test",
		true,
	)
	expiration := time.Now()
	basic := feegrant.BasicAllowance{
		SpendLimit: nil,
		Expiration: &expiration,
	}
	period := feegrant.PeriodicAllowance{
		Basic:  basic,
		Period: time.Duration(3600),
	}
	msgGrantAllowance1, _ := feegrant.NewMsgGrantAllowance(&basic, addr, addr)
	msgGrantAllowance2, _ := feegrant.NewMsgGrantAllowance(&period, addr, addr)
	msgCreateGroup := &group.MsgCreateGroup{
		Admin: addr.String(),
		Members: []group.MemberRequest{
			{
				Address:  addr.String(),
				Weight:   "1",
				Metadata: "metaData",
			},
		},
		Metadata: "metaData",
	}
	testCases := []sdk.Msg{
		msgProposal,
		msgSend,
		msgGrantAllowance1,
		msgGrantAllowance2,
		msgCreateGroup,
	}

	for _, tc := range testCases {
		err = txBuilder.SetMsgs(tc)
		require.NoError(t, err)
		signBytes, err := modeHandler.GetSignBytes(signingtypes.SignMode_SIGN_MODE_EIP_712, signingData, txBuilder.GetTx())
		require.NoError(t, err)
		require.NotNil(t, signBytes)
	}
}
