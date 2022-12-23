package bls

import (
	"testing"

	proto "github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

var _ cryptotypes.PubKey = (*PubKey)(nil)

func TestPKSuite(t *testing.T) {
	suite.Run(t, new(PKSuite))
}

type CommonSuite struct {
	suite.Suite
	pk *PubKey // cryptotypes.PubKey
	sk cryptotypes.PrivKey
}

func (suite *CommonSuite) SetupSuite() {
	sk, err := GenerateKey()
	suite.Require().NoError(err)
	suite.sk = sk
	suite.pk = sk.PubKey().(*PubKey)
}

type PKSuite struct{ CommonSuite }

func (suite *PKSuite) TestType() {
	suite.Require().Equal(KeyType, suite.pk.Type())
}

func (suite *PKSuite) TestBytes() {
	var pk *PubKey
	suite.Nil(pk.Bytes())
}

func (suite *PKSuite) TestEquals() {
	require := suite.Require()

	skOther, err := GenerateKey()
	require.NoError(err)
	pkOther := skOther.PubKey()
	pkOther2 := &PubKey{Key: skOther.PubKey().Bytes()}

	require.False(suite.pk.Equals(pkOther))
	require.True(pkOther.Equals(pkOther2))
	require.True(pkOther2.Equals(pkOther))
	require.True(pkOther.Equals(pkOther), "Equals must be reflexive")
}

func (suite *PKSuite) TestMarshalProto() {
	require := suite.Require()

	/**** test structure marshalling ****/

	var pk PubKey
	bz, err := proto.Marshal(suite.pk)
	require.NoError(err)
	require.NoError(proto.Unmarshal(bz, &pk))
	require.True(pk.Equals(suite.pk))

	/**** test structure marshalling with codec ****/

	pk = PubKey{}
	registry := types.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	bz, err = cdc.Marshal(suite.pk)
	require.NoError(err)
	require.NoError(cdc.Unmarshal(bz, &pk))
	require.True(pk.Equals(suite.pk))

	const bufSize = 100
	bz2 := make([]byte, bufSize)
	pkCpy := suite.pk
	_, err = pkCpy.MarshalTo(bz2)
	require.NoError(err)
	require.Len(bz2, bufSize)
	require.Equal(bz, bz2[:pk.Size()])

	bz2 = make([]byte, bufSize)
	_, err = pkCpy.MarshalToSizedBuffer(bz2)
	require.NoError(err)
	require.Len(bz2, bufSize)
	require.Equal(bz, bz2[(bufSize-pk.Size()):])

	/**** test interface marshalling ****/
	bz, err = cdc.MarshalInterface(suite.pk)
	require.NoError(err)
	var pkI cryptotypes.PubKey
	err = cdc.UnmarshalInterface(bz, &pkI)
	require.EqualError(err, "no registered implementations of type types.PubKey")

	RegisterInterfaces(registry)
	require.NoError(cdc.UnmarshalInterface(bz, &pkI))
	require.True(pkI.Equals(suite.pk))

	require.Error(cdc.UnmarshalInterface(bz, nil), "nil should fail")
}
