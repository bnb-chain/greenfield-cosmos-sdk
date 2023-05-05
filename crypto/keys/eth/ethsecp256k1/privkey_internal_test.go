package ethsecp256k1

import (
	"testing"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

var _ cryptotypes.PrivKey = &PrivKey{}

func TestSKSuite(t *testing.T) {
	suite.Run(t, new(SKSuite))
}

type SKSuite struct{ CommonSuite }

func (suite *SKSuite) TestEquals() {
	require := suite.Require()

	skOther, err := GenPrivKey()
	require.NoError(err)
	require.False(suite.sk.Equals(skOther))

	skOther2 := &PrivKey{skOther.Bytes()}
	require.True(skOther.Equals(skOther2))
	require.True(skOther2.Equals(skOther), "Equals must be reflexive")
}

func (suite *SKSuite) TestBytes() {
	bz := suite.sk.Bytes()
	suite.Len(bz, PrivKeySize)
	var sk *PrivKey
	suite.Nil(sk.Bytes())
}

func (suite *SKSuite) TestMarshalProto() {
	require := suite.Require()

	/**** test structure marshalling ****/

	var sk PrivKey
	bz, err := proto.Marshal(suite.sk)
	require.NoError(err)
	require.NoError(proto.Unmarshal(bz, &sk))
	require.True(sk.Equals(suite.sk))

	/**** test structure marshalling with codec ****/

	sk = PrivKey{}
	registry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	bz, err = cdc.Marshal(suite.sk.(*PrivKey))
	require.NoError(err)
	require.NoError(cdc.Unmarshal(bz, &sk))
	require.True(sk.Equals(suite.sk))

	const bufSize = 100
	bz2 := make([]byte, bufSize)
	skCpy := suite.sk.(*PrivKey)
	_, err = skCpy.MarshalTo(bz2)
	require.NoError(err)
	require.Len(bz2, bufSize)
	require.Equal(bz, bz2[:sk.Size()])

	bz2 = make([]byte, bufSize)
	_, err = skCpy.MarshalToSizedBuffer(bz2)
	require.NoError(err)
	require.Len(bz2, bufSize)
	require.Equal(bz, bz2[(bufSize-sk.Size()):])
}

func (suite *SKSuite) TestSign() {
	require := suite.Require()

	msg := crypto.CRandBytes(1000)
	sig, err := suite.sk.Sign(msg)
	require.NoError(err)
	sigCpy := make([]byte, len(sig))
	copy(sigCpy, sig)
	require.True(suite.pk.VerifySignature(msg, sigCpy))

	// Mutate the signature
	for i := range sig {
		sigCpy[i] ^= byte(i + 1)
		require.False(suite.pk.VerifySignature(msg, sigCpy))
	}

	// Mutate the message
	msg[1] ^= byte(2)
	require.False(suite.pk.VerifySignature(msg, sig))
}
