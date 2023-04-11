package keeper

import (
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type hexCodec struct{}

var _ address.Codec = &hexCodec{}

func NewHexCodec() address.Codec {
	return hexCodec{}
}

// StringToBytes encodes text to bytes
func (bc hexCodec) StringToBytes(text string) ([]byte, error) {
	addr, err := sdk.AccAddressFromHexUnsafe(text)
	if err != nil {
		return nil, err
	}
	return addr.Bytes(), nil
}

// BytesToString decodes bytes to text
func (bc hexCodec) BytesToString(bz []byte) (string, error) {
	return sdk.AccAddress(bz).String(), nil
}
