package ethsecp256k1

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/stretchr/testify/require"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func TestPrivKey(t *testing.T) {
	// validate type and equality
	privKey, err := GenPrivKey()
	require.NoError(t, err)
	require.True(t, privKey.Equals(privKey)) //nolint:gocritic
	require.Implements(t, (*cryptotypes.PrivKey)(nil), privKey)

	// validate inequality
	privKey2, err := GenPrivKey()
	require.NoError(t, err)
	require.False(t, privKey.Equals(privKey2))

	// validate Ethereum address equality
	addr := privKey.PubKey().Address()
	key, err := privKey.ToECDSA()
	require.NoError(t, err)
	expectedAddr := crypto.PubkeyToAddress(key.PublicKey)
	require.Equal(t, expectedAddr.Bytes(), addr.Bytes())

	// validate we can sign some bytes
	msg := []byte("hello world")
	sigHash := crypto.Keccak256Hash(msg)
	expectedSig, err := secp256k1.Sign(sigHash.Bytes(), privKey.Bytes())
	require.NoError(t, err)

	sig, err := privKey.Sign(sigHash.Bytes())
	require.NoError(t, err)
	require.Equal(t, expectedSig, sig)
}

func TestPrivKey_PubKey(t *testing.T) {
	privKey, err := GenPrivKey()
	require.NoError(t, err)

	// validate type and equality
	pubKey := &PubKey{
		Key: privKey.PubKey().Bytes(),
	}
	require.Implements(t, (*cryptotypes.PubKey)(nil), pubKey)

	// validate inequality
	privKey2, err := GenPrivKey()
	require.NoError(t, err)
	require.False(t, pubKey.Equals(privKey2.PubKey()))

	// validate signature
	msg := []byte("hello world")
	sigHash := crypto.Keccak256Hash(msg)
	sig, err := privKey.Sign(sigHash.Bytes())
	require.NoError(t, err)

	res := pubKey.VerifySignature(msg, sig)
	require.True(t, res)
}
