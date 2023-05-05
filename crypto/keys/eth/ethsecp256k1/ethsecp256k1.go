package ethsecp256k1

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

const (
	// PrivKeySize defines the size of the PrivKey bytes
	PrivKeySize = 32
	// PubKeySize defines the size of the PubKey bytes
	PubKeySize = 33
	// KeyType is the string constant for the Secp256k1 algorithm
	KeyType = "eth_secp256k1"
)

// Amino encoding names
const (
	// PrivKeyName defines the amino encoding name for the EthSecp256k1 private key
	PrivKeyName = "ethereum/PrivKeyEthSecp256k1"
	// PubKeyName defines the amino encoding name for the EthSecp256k1 public key
	PubKeyName = "ethereum/PubKeyEthSecp256k1"
)

// RegisterInterfaces adds eth_secp256k1 PubKey to pubkey registry
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &PubKey{})
}
