package keyring

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
)

var (
	// SupportedAlgorithms defines the list of signing algorithms used on Greenfield:
	//  - eth_secp256k1 (Ethereum)
	//  - eth_bls (Ethereum)
	//  - secp256k1 (For legacy accounts)
	SupportedAlgorithms = SigningAlgoList{hd.EthSecp256k1, hd.EthBLS, hd.Secp256k1}
	// SupportedAlgorithmsLedger defines the list of signing algorithms used on Greenfield for the Ledger device:
	//  - eth_secp256k1 (Ethereum)
	//  - eth_bls (Ethereum)
	//  - secp256k1 (For legacy accounts)
	SupportedAlgorithmsLedger = SigningAlgoList{hd.EthSecp256k1, hd.EthBLS, hd.Secp256k1}
)

// ETHAlgoOption defines a function keys options for the ethereum Secp256k1 curve.
// It supports eth_secp256k1, eth_bls keys for accounts.
func ETHAlgoOption() Option {
	return func(options *Options) {
		options.SupportedAlgos = SupportedAlgorithms
		options.SupportedAlgosLedger = SupportedAlgorithmsLedger
	}
}
