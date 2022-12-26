package bls

import (
	"crypto/subtle"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/prysmaticlabs/prysm/crypto/bls"
)

var (
	_ cryptotypes.PrivKey  = &PrivKey{}
	_ codec.AminoMarshaler = &PrivKey{}
)

// GenerateKey generates a new random private key. It returns an error upon
// failure.
func GenerateKey() (*PrivKey, error) {
	secretKey, err := bls.RandKey()
	if err != nil {
		return nil, err
	}

	return &PrivKey{
		Key: secretKey.Marshal(),
	}, nil
}

// Bytes returns the byte representation of the BLS Private Key.
func (privKey *PrivKey) Bytes() []byte {
	if privKey == nil {
		return nil
	}
	bz := make([]byte, len(privKey.Key))
	copy(bz, privKey.Key)

	return bz
}

// PubKey returns the BLS private key's public key. If the privkey is not valid
// it returns a nil value.
func (privKey *PrivKey) PubKey() cryptotypes.PubKey {
	secretKey, err := bls.SecretKeyFromBytes(privKey.Bytes())
	if err != nil {
		return nil
	}

	return &PubKey{
		Key: secretKey.PublicKey().Marshal(),
	}
}

// Equals returns true if two BLS private keys are equal and false otherwise.
func (privKey *PrivKey) Equals(other cryptotypes.LedgerPrivKey) bool {
	return privKey.Type() == other.Type() && subtle.ConstantTimeCompare(privKey.Bytes(), other.Bytes()) == 1
}

// Type returns eth_bls
func (privKey *PrivKey) Type() string {
	return KeyType
}

// MarshalAmino overrides Amino binary marshaling.
func (privKey *PrivKey) MarshalAmino() ([]byte, error) {
	return privKey.Key, nil
}

// UnmarshalAmino overrides Amino binary marshaling.
func (privKey *PrivKey) UnmarshalAmino(bz []byte) error {
	privKey.Key = bz

	return nil
}

// MarshalAminoJSON overrides Amino JSON marshaling.
func (privKey *PrivKey) MarshalAminoJSON() ([]byte, error) {
	// When we marshal to Amino JSON, we don't marshal the "key" field itself,
	// just its contents (i.e. the key bytes).
	return privKey.MarshalAmino()
}

// UnmarshalAminoJSON overrides Amino JSON marshaling.
func (privKey *PrivKey) UnmarshalAminoJSON(bz []byte) error {
	return privKey.UnmarshalAmino(bz)
}

// Sign a message using a secret key - in a beacon/validator client.
//
// In IETF draft BLS specification:
// Sign(SK, message) -> signature: a signing algorithm that generates
//
//	a deterministic signature given a secret key SK and a message.
func (privKey *PrivKey) Sign(digestBz []byte) ([]byte, error) {
	secretKey, err := bls.SecretKeyFromBytes(privKey.Bytes())
	if err != nil {
		return nil, err
	}

	return secretKey.Sign(digestBz).Marshal(), nil
}
