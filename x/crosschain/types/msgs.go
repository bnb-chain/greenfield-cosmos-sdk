package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgUpdateChannelPermissions{}
	_ sdk.Msg = &MsgMintModuleTokens{}
)

// GetSignBytes implements the LegacyMsg interface.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHexUnsafe(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := m.Params.Validate(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgUpdateChannelPermissions) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateChannelPermissions) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHexUnsafe(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateChannelPermissions) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	for _, per := range m.ChannelPermissions {
		if per == nil {
			return fmt.Errorf("channel permission is nil")
		}
	}

	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgMintModuleTokens) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgMintModuleTokens message.
func (m *MsgMintModuleTokens) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHexUnsafe(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgMintModuleTokens) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if m.Amount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("amount must be positive, is %s", m.Amount)
	}

	return nil
}
