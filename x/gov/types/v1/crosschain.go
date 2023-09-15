package v1

import (
	"encoding/hex"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewCrossChainParamsChange(key string, values, targets []string) *CrossChainParamsChange {
	return &CrossChainParamsChange{Key: key, Values: values, Targets: targets}
}

func (m *CrossChainParamsChange) ValidateBasic() error {
	if m.Key == "" || len(m.Values) == 0 || len(m.Targets) == 0 {
		return types.ErrEmptyChange
	}
	if len(m.Values) != len(m.Targets) {
		return types.ErrAddressSizeNotMatch
	}
	if m.Key != types.KeyUpgrade && len(m.Values) > 1 {
		return types.ErrExceedParamsChangeLimit
	}

	for i := 0; i < len(m.Values); i++ {
		value := m.Values[i]
		target := m.Targets[i]
		if len(strings.TrimSpace(value)) == 0 {
			return types.ErrEmptyValue
		}
		if len(strings.TrimSpace(target)) == 0 {
			return types.ErrEmptyTarget
		}
		if m.Key == types.KeyUpgrade {
			_, err := sdk.AccAddressFromHexUnsafe(value)
			if err != nil {
				return types.ErrAddressNotValid
			}
		} else {
			_, err := hex.DecodeString(value)
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidValue, "value is not valid %s", value)
			}
		}
		_, err := sdk.AccAddressFromHexUnsafe(target)
		if err != nil {
			return types.ErrAddressNotValid
		}
	}
	return nil
}
