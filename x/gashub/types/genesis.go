package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Validate performs basic validation of supply genesis data returning an
// error for any failed validation criteria.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	seenMsgGasParams := make(map[string]bool)

	for _, mgp := range gs.GetMsgGasParams() {
		if _, exists := seenMsgGasParams[mgp.MsgTypeUrl]; exists {
			return fmt.Errorf("duplicate msg gas params found: '%s'", mgp.MsgTypeUrl)
		}
		if err := mgp.Validate(); err != nil {
			return err
		}
		seenMsgGasParams[mgp.MsgTypeUrl] = true
	}

	return nil
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, msgGasParamsSet []MsgGasParams) *GenesisState {
	return &GenesisState{
		Params:       params,
		MsgGasParams: msgGasParamsSet,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() *GenesisState {
	defaultMsgGasParamsSet := []MsgGasParams{
		*NewMsgGasParamsWithFixedGas("/cosmos.auth.v1beta1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.bank.v1beta1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.consensus.v1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.crosschain.v1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.gashub.v1beta1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.mint.v1beta1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.oracle.v1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.slashing.v1beta1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/greenfield.bridge.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/greenfield.sp.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/greenfield.payment.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/greenfield.challenge.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/greenfield.permission.MsgUpdateParams", 0),
		*NewMsgGasParamsWithFixedGas("/cosmos.authz.v1beta1.MsgExec", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.authz.v1beta1.MsgRevoke", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.bank.v1beta1.MsgSend", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.feegrant.v1beta1.MsgRevokeAllowance", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgDeposit", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgSubmitProposal", 2e8),
		*NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgVote", 2e7),
		*NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgVoteWeighted", 2e7),
		*NewMsgGasParamsWithFixedGas("/cosmos.oracle.v1.MsgClaim", 1e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.slashing.v1beta1.MsgUnjail", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgBeginRedelegate", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCreateValidator", 2e8),
		*NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgDelegate", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgEditValidator", 2e7),
		*NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgUndelegate", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.bridge.MsgTransferOut", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.sp.MsgCreateStorageProvider", 2e8),
		*NewMsgGasParamsWithFixedGas("/greenfield.sp.MsgDeposit", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.sp.MsgEditStorageProvider", 2e7),
		*NewMsgGasParamsWithFixedGas("/greenfield.sp.MsgUpdateSpStoragePrice", 2e7),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgCreateBucket", 2.4e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgDeleteBucket", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgMirrorBucket", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgUpdateBucketInfo", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgCreateObject", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgSealObject", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgMirrorObject", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgRejectSealObject", 1.2e4),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgDeleteObject", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgCopyObject", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgCancelCreateObject", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgUpdateObjectInfo", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgDiscontinueObject", 2.4e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgDiscontinueBucket", 2.4e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgCreateGroup", 2.4e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgDeleteGroup", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgLeaveGroup", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgUpdateGroupMember", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgMirrorGroup", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgPutPolicy", 2.4e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.storage.MsgDeletePolicy", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.payment.MsgCreatePaymentAccount", 2e5),
		*NewMsgGasParamsWithFixedGas("/greenfield.payment.MsgDeposit", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.payment.MsgWithdraw", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.payment.MsgDisableRefund", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.challenge.MsgSubmit", 1.2e3),
		*NewMsgGasParamsWithFixedGas("/greenfield.challenge.MsgAttest", 1e2),
		*NewMsgGasParamsWithDynamicGas(
			"/cosmos.authz.v1beta1.MsgGrant",
			&MsgGasParams_GrantType{
				GrantType: &MsgGasParams_DynamicGasParams{
					FixedGas:   8e2,
					GasPerItem: 8e2,
				},
			},
		),
		*NewMsgGasParamsWithDynamicGas(
			"/cosmos.bank.v1beta1.MsgMultiSend",
			&MsgGasParams_MultiSendType{
				MultiSendType: &MsgGasParams_DynamicGasParams{
					FixedGas:   8e2,
					GasPerItem: 8e2,
				},
			},
		),
		*NewMsgGasParamsWithDynamicGas(
			"/cosmos.feegrant.v1beta1.MsgGrantAllowance",
			&MsgGasParams_GrantAllowanceType{
				GrantAllowanceType: &MsgGasParams_DynamicGasParams{
					FixedGas:   8e2,
					GasPerItem: 8e2,
				},
			},
		),
	}
	return NewGenesisState(DefaultParams(), defaultMsgGasParamsSet)
}

// GetGenesisStateFromAppState returns x/gashub GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}
