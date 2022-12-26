package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// slashing message types
const (
	TypeMsgUnjail  = "unjail"
	TypeMsgImpeach = "impeach"
)

// verify interface at compile time
var _ sdk.Msg = &MsgUnjail{}

// NewMsgUnjail creates a new MsgUnjail instance
func NewMsgUnjail(validatorAddr sdk.AccAddress) *MsgUnjail {
	return &MsgUnjail{
		ValidatorAddr: validatorAddr.String(),
	}
}

func (msg MsgUnjail) Route() string { return RouterKey }
func (msg MsgUnjail) Type() string  { return TypeMsgUnjail }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	valAddr, _ := sdk.AccAddressFromHexUnsafe(msg.ValidatorAddr)
	return []sdk.AccAddress{valAddr}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does a sanity check on the provided message
func (msg MsgUnjail) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(msg.ValidatorAddr); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("validator input address: %s", err)
	}
	return nil
}

// NewMsgImpeach creates a new MsgImpeach instance
func NewMsgImpeach(valAddr sdk.AccAddress, from sdk.AccAddress) *MsgImpeach {
	return &MsgImpeach{
		ValidatorAddress: valAddr.String(),
		From:             from.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgImpeach) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgImpeach) Type() string { return TypeMsgImpeach }

// GetSigners implements the sdk.Msg interface.
func (msg MsgImpeach) GetSigners() []sdk.AccAddress {
	fromAddr, _ := sdk.AccAddressFromHexUnsafe(msg.From)
	return []sdk.AccAddress{fromAddr}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgImpeach) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgImpeach) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(msg.From); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid account address: %s", err)
	}

	if _, err := sdk.AccAddressFromHexUnsafe(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	return nil
}
