package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// slashing message types
const (
	TypeMsgUnjail  = "unjail"
	TypeMsgKickOut = "kick_out"
)

// verify interface at compile time
var _ sdk.Msg = &MsgUnjail{}

// NewMsgUnjail creates a new MsgUnjail instance
//
//nolint:interfacer
func NewMsgUnjail(validatorAddr sdk.ValAddress) *MsgUnjail {
	return &MsgUnjail{
		ValidatorAddr: validatorAddr.String(),
	}
}

func (msg MsgUnjail) Route() string { return RouterKey }
func (msg MsgUnjail) Type() string  { return TypeMsgUnjail }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	valAddr, _ := sdk.ValAddressFromHex(msg.ValidatorAddr)
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does a sanity check on the provided message
func (msg MsgUnjail) ValidateBasic() error {
	if _, err := sdk.ValAddressFromHex(msg.ValidatorAddr); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("validator input address: %s", err)
	}
	return nil
}

// NewMsgKickOut creates a new MsgKickOut instance
//
//nolint:interfacer
func NewMsgKickOut(valAddr sdk.ValAddress, from sdk.AccAddress) *MsgKickOut {
	return &MsgKickOut{
		ValidatorAddress: valAddr.String(),
		From:             from.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgKickOut) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgKickOut) Type() string { return TypeMsgKickOut }

// GetSigners implements the sdk.Msg interface.
func (msg MsgKickOut) GetSigners() []sdk.AccAddress {
	fromAddr, _ := sdk.AccAddressFromBech32(msg.From)
	return []sdk.AccAddress{fromAddr}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgKickOut) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgKickOut) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid account address: %s", err)
	}

	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	return nil
}
