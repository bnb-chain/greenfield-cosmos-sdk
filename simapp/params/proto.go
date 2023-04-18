//go:build !test_amino
// +build !test_amino

package params

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
)

// MakeTestEncodingConfig creates an EncodingConfig for a non-amino based test configuration.
// This function should be used only internally (in the SDK).
// App user shouldn't create new codecs - use the app.AppCodec instead.
// [DEPRECATED]
func MakeTestEncodingConfig() EncodingConfig {
	cdc := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(codec, []signingtypes.SignMode{
		signingtypes.SignMode_SIGN_MODE_EIP_712,
	})

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          txCfg,
		Amino:             cdc,
	}
}
