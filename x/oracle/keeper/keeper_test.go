package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type TestSuite struct {
	suite.Suite

	app *simapp.SimApp
	ctx sdk.Context
}

func (s *TestSuite) SetupTest() {
	app := simapp.Setup(s.T(), false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})

	app.OracleKeeper.SetParams(ctx, types.DefaultParams())

	s.app = app
	s.ctx = ctx
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestKeeper_IsValidatorInturn() {
	s.app.OracleKeeper.SetParams(s.ctx, types.Params{
		RelayerTimeout:     5,
		RelayerBackoffTime: 3,
	})

	vals := make([]stakingtypes.Validator, 5)
	for i := range vals {
		pk := ed25519.GenPrivKey().PubKey()
		vals[i] = newValidator(s.T(), sdk.ValAddress(pk.Address()), pk)
	}

	val0Addr, err := vals[0].GetConsAddr()
	s.Require().Nil(err)
	val1Addr, err := vals[1].GetConsAddr()
	s.Require().Nil(err)

	tests := []struct {
		claimMsg     types.MsgClaim
		blockTime    int64
		expectedPass bool
		errorMsg     string
	}{
		{
			types.MsgClaim{
				FromAddress:    val0Addr.String(),
				ChainId:        1,
				Sequence:       1,
				Timestamp:      1990,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			1992,
			true,
			"",
		},
		// wrong validator in timeout
		{
			types.MsgClaim{
				FromAddress:    val1Addr.String(),
				ChainId:        1,
				Sequence:       1,
				Timestamp:      1990,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			1992,
			false,
			"",
		},
		// right validator in backoff time
		{
			types.MsgClaim{
				FromAddress:    val1Addr.String(),
				ChainId:        1,
				Sequence:       1,
				Timestamp:      1985,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			1992,
			true,
			"",
		},
		// wrong validator in backoff time
		{
			types.MsgClaim{
				FromAddress:    val0Addr.String(),
				ChainId:        1,
				Sequence:       1,
				Timestamp:      1985,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			1992,
			false,
			"",
		},
		// right validator in backoff time
		{
			types.MsgClaim{
				FromAddress:    val0Addr.String(),
				ChainId:        1,
				Sequence:       1,
				Timestamp:      1970,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			1989,
			true,
			"",
		},
	}

	for _, test := range tests {
		s.ctx = s.ctx.WithBlockTime(time.Unix(test.blockTime, 0))
		isInturn, err := s.app.OracleKeeper.IsValidatorInturn(s.ctx, vals, &test.claimMsg)

		if test.expectedPass {
			s.Require().Nil(err)
			s.Require().True(isInturn)
		} else {
			s.Require().False(isInturn)
		}
	}

}

// Creates a new validators and asserts the error check.
func newValidator(t *testing.T, operator sdk.ValAddress, pubKey cryptotypes.PubKey) stakingtypes.Validator {
	v, err := stakingtypes.NewValidator(operator, pubKey, stakingtypes.Description{})
	require.NoError(t, err)
	return v
}
