package bls

import (
	"bytes"

	cmtcrypto "github.com/cometbft/cometbft/crypto"
	"github.com/prysmaticlabs/prysm/crypto/bls"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

// Address returns the address of the BLS public key.
// The function will return an empty address if the public key is invalid.
func (pubKey *PubKey) Address() cmtcrypto.Address {
	if pubKey == nil {
		return cmtcrypto.Address(nil)
	}

	return cmtcrypto.Address(pubKey.Bytes())
}

// Bytes returns the raw bytes of the BLS public key.
func (pubKey *PubKey) Bytes() []byte {
	if pubKey == nil {
		return nil
	}
	bz := make([]byte, len(pubKey.Key))
	copy(bz, pubKey.Key)

	return bz
}

// String implements the fmt.Stringer interface.
func (pubKey *PubKey) String() string {
	return pubKey.Address().String()
}

// Type returns eth_bls
func (pubKey *PubKey) Type() string {
	return KeyType
}

// Equals returns true if the pubkey type is the same and their bytes are deeply equal.
func (pubKey *PubKey) Equals(other cryptotypes.PubKey) bool {
	return pubKey.Type() == other.Type() && bytes.Equal(pubKey.Bytes(), other.Bytes())
}

// MarshalAmino overrides Amino binary marshalling.
func (pubKey PubKey) MarshalAmino() ([]byte, error) {
	return pubKey.Key, nil
}

// UnmarshalAmino overrides Amino binary marshalling.
func (pubKey *PubKey) UnmarshalAmino(bz []byte) error {
	pubKey.Key = bz

	return nil
}

// MarshalAminoJSON overrides Amino JSON marshalling.
func (pubKey PubKey) MarshalAminoJSON() ([]byte, error) {
	// When we marshal to Amino JSON, we don't marshal the "key" field itself,
	// just its contents (i.e. the key bytes).
	return pubKey.MarshalAmino()
}

// UnmarshalAminoJSON overrides Amino JSON marshalling.
func (pubKey *PubKey) UnmarshalAminoJSON(bz []byte) error {
	return pubKey.UnmarshalAmino(bz)
}

// VerifySignature verifies that the BLS public key created a given signature over
// the provided message.
func (pubKey *PubKey) VerifySignature(msg, sig []byte) bool {
	key, err := bls.PublicKeyFromBytes(pubKey.Bytes())
	if err != nil {
		return false
	}
	signature, err := bls.SignatureFromBytes(sig)
	if err != nil {
		return false
	}

	return signature.Verify(key, msg)
}
