package types

import (
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// bank message types
const (
	TypeMsgUpdateMsgGasParams = "update_msg_gas_params"
)

var _ sdk.Msg = &MsgUpdateMsgGasParams{}

// MsgUpdateMsgGasParams - construct a msg to update msg gas params.
//
//nolint:interfacer
func NewMsgUpdateMsgGasParams(from sdk.AccAddress, msgGasParams *MsgGasParams) *MsgUpdateMsgGasParams {
	return &MsgUpdateMsgGasParams{
		From:      from.String(),
		NewParams: msgGasParams,
	}
}

// Route Implements Msg.
func (msg MsgUpdateMsgGasParams) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgUpdateMsgGasParams) Type() string { return TypeMsgUpdateMsgGasParams }

// ValidateBasic Implements Msg.
func (msg MsgUpdateMsgGasParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(msg.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	if err := validateMsgGasParams(msg.NewParams); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid msg gas params: %s", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateMsgGasParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgUpdateMsgGasParams) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromHexUnsafe(msg.From)
	return []sdk.AccAddress{fromAddress}
}
