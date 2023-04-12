package bls

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

const (
	// KeyType is the string constant for the BLS algorithm
	KeyType = "eth_bls"
)

// Amino encoding names
const (
	// PrivKeyName defines the amino encoding name for the EthBLS private key
	PrivKeyName = "ethereum/PrivKeyEthBLS"
	// PubKeyName defines the amino encoding name for the EthBLS public key
	PubKeyName = "ethereum/PubKeyEthBLS"
)

// RegisterInterfaces adds BLS PubKey to pubkey registry
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &PubKey{})
}
