package hd

import (
	"strings"

	"github.com/cosmos/go-bip39"
	util "github.com/wealdtech/go-eth2-util"

	"github.com/cosmos/cosmos-sdk/crypto/keys/eth/bls"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types"
)

// PubKeyType defines an algorithm to derive key-pairs which can be used for cryptographic signing.
type PubKeyType string

const (
	// MultiType implies that a pubkey is a multisignature
	MultiType = PubKeyType("multi")
	// Secp256k1Type uses the Bitcoin secp256k1 ECDSA parameters.
	Secp256k1Type = PubKeyType("secp256k1")
	// Ed25519Type represents the Ed25519Type signature system.
	// It is currently not supported for end-user keys (wallets/ledgers).
	Ed25519Type = PubKeyType("ed25519")
	// Sr25519Type represents the Sr25519Type signature system.
	Sr25519Type = PubKeyType("sr25519")
)

// Secp256k1 uses the Bitcoin secp256k1 ECDSA parameters.
var Secp256k1 = secp256k1Algo{}

type (
	DeriveFn   func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error)
	GenerateFn func(bz []byte) types.PrivKey
)

type WalletGenerator interface {
	Derive(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error)
	Generate(bz []byte) types.PrivKey
}

type secp256k1Algo struct{}

func (s secp256k1Algo) Name() PubKeyType {
	return Secp256k1Type
}

// Derive derives and returns the secp256k1 private key for the given seed and HD path.
func (s secp256k1Algo) Derive() DeriveFn {
	return func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error) {
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		masterPriv, ch := ComputeMastersFromSeed(seed)
		if len(hdPath) == 0 {
			return masterPriv[:], nil
		}
		derivedKey, err := DerivePrivateKeyForPath(masterPriv, ch, hdPath)

		return derivedKey, err
	}
}

// Generate generates a secp256k1 private key from the given bytes.
func (s secp256k1Algo) Generate() GenerateFn {
	return func(bz []byte) types.PrivKey {
		bzArr := make([]byte, secp256k1.PrivKeySize)
		copy(bzArr, bz)

		return &secp256k1.PrivKey{Key: bzArr}
	}
}

const (
	// BLSType uses the ethereum BLS parameters.
	BLSType = PubKeyType(bls.KeyType)
)

// EthBLS uses the Bitcoin eth_bls parameters.
var EthBLS = ethBLSAlgo{}

type ethBLSAlgo struct{}

// Name returns eth_bls
func (s ethBLSAlgo) Name() PubKeyType {
	return BLSType
}

// Derive derives and returns the eth_bls private key for the given seed and HD path.
func (s ethBLSAlgo) Derive() DeriveFn {
	// Derive derives and returns the eth_bls private key for the given mnemonic and HD path.
	return func(mnemonic, bip39Passphrase, path string) ([]byte, error) {
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		privKey, err := util.PrivateKeyFromSeedAndPath(
			seed, strings.ReplaceAll(path, "'", ""),
		)
		if err != nil {
			return nil, err
		}

		return privKey.Marshal(), nil
	}
}

// Generate generates a eth_bls private key from the given bytes.
func (s ethBLSAlgo) Generate() GenerateFn {
	return func(bz []byte) types.PrivKey {
		return &bls.PrivKey{
			Key: bz,
		}
	}
}
