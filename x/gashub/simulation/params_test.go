package simulation_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/x/gashub/simulation"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

func TestParamChanges(t *testing.T) {
	s := rand.NewSource(1)
	r := rand.New(s)

	expected := []struct {
		composedKey string
		key         string
		simValue    string
		subspace    string
	}{
		{"gashub/MaxTxSize", "MaxTxSize", "\"3081\"", "gashub"},
	}

	paramChanges := simulation.ParamChanges(r)

	require.Len(t, paramChanges, 1)

	for i, p := range paramChanges {
		require.Equal(t, expected[i].composedKey, p.ComposedKey())
		require.Equal(t, expected[i].key, p.Key())
		require.Equal(t, expected[i].simValue, p.SimValue()(r))
		require.Equal(t, expected[i].subspace, p.Subspace())
	}
}

func TestMsgUrl(t *testing.T) {
	defaultMsgGasParamsSet := []*types.MsgGasParams{
		types.NewMsgGasParamsWithFixedGas("/cosmos.authz.v1beta1.MsgExec", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.authz.v1beta1.MsgRevoke", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.bank.v1beta1.MsgSend", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgFundCommunityPool", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.feegrant.v1beta1.MsgRevokeAllowance", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgDeposit", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgSubmitProposal", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgVote", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgVoteWeighted", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.oracle.v1.MsgClaim", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.slashing.v1beta1.MsgUnjail", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgBeginRedelegate", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCreateValidator", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgDelegate", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgEditValidator", 12e3),
		types.NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgUndelegate", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.bridge.MsgTransferOut", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.sp.MsgDeposit", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.sp.MsgEditStorageProvider", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCopyObject", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreateBucket", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreateGroup", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreateObject", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDeleteBucket", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDeleteGroup", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgLeaveGroup", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgRejectSealObject", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgSealObject", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgUpdateGroupMember", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.payment.MsgCreatePaymentAccount", 2e6),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.payment.MsgDeposit", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.payment.MsgWithdraw", 12e3),
		types.NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.payment.MsgDisableRefund", 12e3),
		types.NewMsgGasParamsWithDynamicGas(
			"/cosmos.authz.v1beta1.MsgGrant",
			&types.MsgGasParams_GrantType{
				GrantType: &types.MsgGasParams_DynamicGasParams{
					FixedGas:   8e3,
					GasPerItem: 8e3,
				},
			},
		),
		types.NewMsgGasParamsWithDynamicGas(
			"/cosmos.bank.v1beta1.MsgMultiSend",
			&types.MsgGasParams_MultiSendType{
				MultiSendType: &types.MsgGasParams_DynamicGasParams{
					FixedGas:   8e3,
					GasPerItem: 8e3,
				},
			},
		),
		types.NewMsgGasParamsWithDynamicGas(
			"/cosmos.feegrant.v1beta1.MsgGrantAllowance",
			&types.MsgGasParams_GrantAllowanceType{
				GrantAllowanceType: &types.MsgGasParams_DynamicGasParams{
					FixedGas:   8e3,
					GasPerItem: 8e3,
				},
			},
		),
	}

	for _, param := range defaultMsgGasParamsSet {
		fmt.Println(param.MsgTypeUrl)
	}
}
