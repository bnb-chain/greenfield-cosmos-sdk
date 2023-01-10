package keeper_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/prysmaticlabs/prysm/crypto/bls/blst"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"github.com/willf/bitset"

	crosschaintypes "github.com/cosmos/cosmos-sdk/x/crosschain/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/cosmos/cosmos-sdk/x/oracle/keeper"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/testutil"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type TestSuite struct {
	suite.Suite

	app *simapp.SimApp
	ctx sdk.Context

	msgServer types.MsgServer
}

func (s *TestSuite) SetupTest() {
	app := simapp.Setup(s.T(), false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})

	app.OracleKeeper.SetParams(ctx, types.DefaultParams())

	s.app = app
	s.ctx = ctx

	s.app.CrossChainKeeper.SetSrcChainID(sdk.ChainID(1))

	coins := sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100000)))
	err := s.app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	s.NoError(err)
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, crosschaintypes.ModuleName, coins)
	s.NoError(err)

	s.msgServer = keeper.NewMsgServerImpl(s.app.OracleKeeper)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestProcessClaim() {
	s.app.OracleKeeper.SetParams(s.ctx, types.Params{
		RelayerTimeout:     5,
		RelayerBackoffTime: 3,
		RelayerRewardShare: 50,
	})

	_, _, newValidators, blsKeys := createValidators(s.T(), s.ctx, s.app, []int64{9, 8, 7})

	validators := s.app.StakingKeeper.GetLastValidators(s.ctx)

	s.app.StakingKeeper.SetHistoricalInfo(s.ctx, s.ctx.BlockHeight(), &stakingtypes.HistoricalInfo{
		Header: s.ctx.BlockHeader(),
		Valset: validators,
	})

	validatorMap := make(map[string]int, 0)
	for idx, validator := range validators {
		validatorMap[validator.RelayerAddress] = idx
	}

	msgClaim := types.MsgClaim{
		FromAddress:    validators[0].RelayerAddress,
		SrcChainId:     1,
		DestChainId:    2,
		Sequence:       1,
		Timestamp:      1992,
		Payload:        []byte("test payload"),
		VoteAddressSet: []uint64{0, 1},
		AggSignature:   []byte("test sig"),
	}

	blsSignBytes := msgClaim.GetBlsSignBytes()

	valBitSet := bitset.New(256)
	for _, newValidator := range newValidators {
		valBitSet.Set(uint(validatorMap[newValidator.RelayerAddress]))
	}

	blsSig := testutil.GenerateBlsSig(blsKeys, blsSignBytes[:])
	msgClaim.VoteAddressSet = valBitSet.Bytes()
	msgClaim.AggSignature = blsSig

	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, err := s.app.OracleKeeper.CheckClaim(s.ctx, &msgClaim)
	s.Require().Nil(err, "error should be nil")

	// not in turn
	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp)+6, 0))
	_, err = s.app.OracleKeeper.CheckClaim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "error should not be nil")
	s.Require().Contains(err.Error(), fmt.Sprintf("relayer(%s) is not in turn", validators[0].RelayerAddress))

	// wrong validator set
	wrongValBitSet := bitset.New(256)
	for i := 0; i < 10; i++ {
		wrongValBitSet.Set(uint(i))
	}
	msgClaim.VoteAddressSet = wrongValBitSet.Bytes()
	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, err = s.app.OracleKeeper.CheckClaim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "error should not be nil")
	s.Require().Contains(err.Error(), "number of validator set is larger than validators")

	// wrong validator set
	wrongValBitSet = bitset.New(256)
	wrongValBitSet.Set(uint(validatorMap[newValidators[0].RelayerAddress]))
	wrongValBitSet.Set(uint(validatorMap[newValidators[1].RelayerAddress]))
	msgClaim.VoteAddressSet = wrongValBitSet.Bytes()
	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, err = s.app.OracleKeeper.CheckClaim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "error should not be nil")
	s.Require().Contains(err.Error(), "not enough validators voted")

	// wrong sig
	msgClaim.VoteAddressSet = valBitSet.Bytes()
	msgClaim.AggSignature = bytes.Repeat([]byte{2}, 96)

	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, err = s.app.OracleKeeper.CheckClaim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "error should not be nil")
	s.Require().Contains(err.Error(), "BLS signature converts failed")
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

	val0Addr := vals[0].RelayerAddress
	val1Addr := vals[1].RelayerAddress

	tests := []struct {
		claimMsg     types.MsgClaim
		blockTime    int64
		expectedPass bool
		errorMsg     string
	}{
		{
			types.MsgClaim{
				FromAddress:    val0Addr,
				SrcChainId:     1,
				DestChainId:    2,
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
				FromAddress:    val1Addr,
				SrcChainId:     1,
				DestChainId:    2,
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
				FromAddress:    val1Addr,
				SrcChainId:     1,
				DestChainId:    2,
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
				FromAddress:    val0Addr,
				SrcChainId:     1,
				DestChainId:    2,
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
				FromAddress:    val0Addr,
				SrcChainId:     1,
				DestChainId:    2,
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
		isInturn, err := s.app.OracleKeeper.IsRelayerInturn(s.ctx, vals, &test.claimMsg)

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
	v, err := stakingtypes.NewSimpleValidator(operator, pubKey, stakingtypes.Description{})
	require.NoError(t, err)
	return v
}

func createValidators(t *testing.T, ctx sdk.Context, app *simapp.SimApp, powers []int64) ([]sdk.AccAddress, []sdk.ValAddress, []stakingtypes.Validator, []bls.SecretKey) {
	addrs := simapp.AddTestAddrsIncremental(app, ctx, 5, app.StakingKeeper.TokensFromConsensusPower(ctx, 300))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	pks := simapp.CreateTestPubKeys(5)

	privKey1, _ := blst.RandKey()
	privKey2, _ := blst.RandKey()
	privKey3, _ := blst.RandKey()

	blsKeys := []bls.SecretKey{privKey1, privKey2, privKey3}

	val1 := teststaking.NewValidator(t, valAddrs[0], pks[0])
	val1.RelayerBlsKey = privKey1.PublicKey().Marshal()

	val2 := teststaking.NewValidator(t, valAddrs[1], pks[1])
	val2.RelayerBlsKey = privKey2.PublicKey().Marshal()

	val3 := teststaking.NewValidator(t, valAddrs[2], pks[2])
	val3.RelayerBlsKey = privKey3.PublicKey().Marshal()

	vals := []stakingtypes.Validator{val1, val2, val3}

	app.StakingKeeper.SetValidator(ctx, val1)
	app.StakingKeeper.SetValidator(ctx, val2)
	app.StakingKeeper.SetValidator(ctx, val3)
	app.StakingKeeper.SetLastValidatorPower(ctx, val1.GetOperator(), 1)
	app.StakingKeeper.SetLastValidatorPower(ctx, val2.GetOperator(), 2)
	app.StakingKeeper.SetLastValidatorPower(ctx, val3.GetOperator(), 3)

	return addrs, valAddrs, vals, blsKeys
}
