package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateParams{}

// GetSigners returns the signer addresses that are expected to sign the result
// of GetSignBytes.
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromHexUnsafe(msg.Authority)
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (msg MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic MsgUpdateParams message validation.
func (msg MsgUpdateParams) ValidateBasic() error {
	return msg.Params.Validate()
}

// NewMsgSetMsgGasParams Construct a message to set one or more MsgGasParams entries.
func NewMsgSetMsgGasParams(authority string, updateSet []*MsgGasParams, deleteSet []string) *MsgSetMsgGasParams {
	return &MsgSetMsgGasParams{
		Authority: authority,
		UpdateSet: updateSet,
		DeleteSet: deleteSet,
	}
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgSetMsgGasParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for MsgSoftwareUpgrade.
func (msg MsgSetMsgGasParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHexUnsafe(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs basic validation on this MsgSetSendEnabled.
func (msg MsgSetMsgGasParams) ValidateBasic() error {
	if len(msg.Authority) > 0 {
		if _, err := sdk.AccAddressFromHexUnsafe(msg.Authority); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
		}
	}

	seen := map[string]bool{}
	for _, mgp := range msg.UpdateSet {
		if _, alreadySeen := seen[mgp.MsgTypeUrl]; alreadySeen {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicate msg gas params entries found for %q", mgp.MsgTypeUrl)
		}

		seen[mgp.MsgTypeUrl] = true

		if err := mgp.Validate(); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrapf("invalid MsgGasParams entries %q: %s", mgp, err)
		}
	}

	return nil
}
