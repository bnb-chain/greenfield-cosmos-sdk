package keeper_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/golang/mock/gomock"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/prysmaticlabs/prysm/crypto/bls/blst"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/willf/bitset"

	otestutil "github.com/cosmos/cosmos-sdk/x/oracle/testutil"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/mint"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/cosmos/cosmos-sdk/x/oracle/keeper"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type TestSuite struct {
	suite.Suite

	ctx sdk.Context

	oracleKeeper keeper.Keeper

	bankKeeper       *types.MockBankKeeper
	crossChainKeeper *types.MockCrossChainKeeper
	stakingKeeper    *types.MockStakingKeeper

	msgServer   types.MsgServer
	queryClient types.QueryClient
}

func (s *TestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx

	ctrl := gomock.NewController(s.T())

	crossChainKeeper := types.NewMockCrossChainKeeper(ctrl)
	bankKeeper := types.NewMockBankKeeper(ctrl)
	stakingKeeper := types.NewMockStakingKeeper(ctrl)

	s.bankKeeper = bankKeeper
	s.crossChainKeeper = crossChainKeeper
	s.stakingKeeper = stakingKeeper

	s.oracleKeeper = keeper.NewKeeper(encCfg.Codec, key, "fee", types.ModuleName, crossChainKeeper, bankKeeper, stakingKeeper)

	s.oracleKeeper.SetParams(s.ctx, types.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, s.oracleKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	s.msgServer = keeper.NewMsgServerImpl(s.oracleKeeper)
	s.queryClient = queryClient
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestProcessClaim() {
	s.oracleKeeper.SetParams(s.ctx, types.Params{
		RelayerTimeout:     5,
		RelayerRewardShare: 50,
		RelayerInterval:    600,
	})

	newValidators, blsKeys := createValidators(s.T())

	s.stakingKeeper.EXPECT().GetHistoricalInfo(gomock.Any(), gomock.Any()).Return(stakingtypes.HistoricalInfo{
		Header: s.ctx.BlockHeader(),
		Valset: newValidators,
	}, true).AnyTimes()

	validatorMap := make(map[string]int, 0)
	for idx, validator := range newValidators {
		validatorMap[validator.RelayerAddress] = idx
	}

	msgClaim := types.MsgClaim{
		FromAddress:    newValidators[0].RelayerAddress,
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

	blsSig := otestutil.GenerateBlsSig(blsKeys, blsSignBytes[:])
	msgClaim.VoteAddressSet = valBitSet.Bytes()
	msgClaim.AggSignature = blsSig

	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, _, err := s.oracleKeeper.CheckClaim(s.ctx, &msgClaim)
	s.Require().Nil(err, "error should be nil")

	// wrong validator set
	wrongValBitSet := bitset.New(256)
	for i := 0; i < 10; i++ {
		wrongValBitSet.Set(uint(i))
	}
	msgClaim.VoteAddressSet = wrongValBitSet.Bytes()
	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, _, err = s.oracleKeeper.CheckClaim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "error should not be nil")
	s.Require().Contains(err.Error(), "number of validator set is larger than validators")

	// wrong validator set
	wrongValBitSet = bitset.New(256)
	wrongValBitSet.Set(uint(validatorMap[newValidators[0].RelayerAddress]))
	wrongValBitSet.Set(uint(validatorMap[newValidators[1].RelayerAddress]))
	msgClaim.VoteAddressSet = wrongValBitSet.Bytes()
	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, _, err = s.oracleKeeper.CheckClaim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "error should not be nil")
	s.Require().Contains(err.Error(), "not enough validators voted")

	// wrong sig
	msgClaim.VoteAddressSet = valBitSet.Bytes()
	msgClaim.AggSignature = bytes.Repeat([]byte{2}, 96)

	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, _, err = s.oracleKeeper.CheckClaim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "error should not be nil")
	s.Require().Contains(err.Error(), "BLS signature converts failed")
}

func (s *TestSuite) TestKeeper_IsRelayerValid() {
	s.oracleKeeper.SetParams(s.ctx, types.Params{
		RelayerTimeout:     5,
		RelayerRewardShare: 50,
		RelayerInterval:    600,
	})

	vals := make([]stakingtypes.Validator, 5)
	for i := range vals {
		pk := ed25519.GenPrivKey().PubKey()

		val := newValidator(s.T(), sdk.AccAddress(pk.Address()), pk)
		privKey, _ := blst.RandKey()
		val.BlsKey = privKey.PublicKey().Marshal()
		vals[i] = val
	}

	s.stakingKeeper.EXPECT().GetHistoricalInfo(gomock.Any(), gomock.Any()).Return(stakingtypes.HistoricalInfo{
		Header: s.ctx.BlockHeader(),
		Valset: vals,
	}, true).AnyTimes()

	val0Addr := vals[0].RelayerAddress
	val1Addr := vals[1].RelayerAddress
	val3Addr := vals[3].RelayerAddress

	// in-turn relayer is relayer 3 and interval is [1800, 2400)
	tests := []struct {
		claimMsg     types.MsgClaim
		blockTime    int64
		expectedPass bool
		errorMsg     string
	}{
		{ // out-turn relayer within the in-turn relayer interval, and not exceeding the timeout, so is not eligible to relay
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
			false,
			"",
		},
		{ // out-turn relayer within the in-turn relayer interval, but exceeding the timeout, so is eligible to relay
			types.MsgClaim{
				FromAddress:    val1Addr,
				SrcChainId:     1,
				DestChainId:    2,
				Sequence:       1,
				Timestamp:      1800,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			1992,
			true,
			"",
		},
		// in-turn relayer, even exceeding timeout
		{
			types.MsgClaim{
				FromAddress:    val3Addr,
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
		// in-turn relayer, within the timeout
		{
			types.MsgClaim{
				FromAddress:    val3Addr,
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
	}

	for idx, test := range tests {
		s.ctx = s.ctx.WithBlockTime(time.Unix(test.blockTime, 0))
		relayer := sdk.MustAccAddressFromHex(test.claimMsg.FromAddress)
		isValid, err := s.oracleKeeper.IsRelayerValid(s.ctx, relayer, vals, test.claimMsg.Timestamp)

		if test.expectedPass {
			s.Require().Nil(err)
			s.Require().True(isValid, fmt.Sprintf("test case %d should be right", idx))
		} else {
			s.Require().False(isValid, fmt.Sprintf("test case %d should be false", idx))
		}
	}
}

// Creates a new validators and asserts the error check.
func newValidator(t *testing.T, operator sdk.AccAddress, pubKey cryptotypes.PubKey) stakingtypes.Validator {
	v, err := stakingtypes.NewSimpleValidator(sdk.ValAddress(operator), pubKey, stakingtypes.Description{})
	require.NoError(t, err)
	return v
}

func createValidators(t *testing.T) ([]stakingtypes.Validator, []bls.SecretKey) {
	addrs := simtestutil.CreateIncrementalAccounts(5)
	valAddrs := simtestutil.ConvertAddrsToValAddrs(addrs)
	pks := simtestutil.CreateTestPubKeys(5)

	privKey1, _ := blst.RandKey()
	privKey2, _ := blst.RandKey()
	privKey3, _ := blst.RandKey()

	blsKeys := []bls.SecretKey{privKey1, privKey2, privKey3}

	val1 := newValidator(t, sdk.AccAddress(valAddrs[0]), pks[0])
	val1.BlsKey = privKey1.PublicKey().Marshal()

	val2 := newValidator(t, sdk.AccAddress(valAddrs[1]), pks[1])
	val2.BlsKey = privKey2.PublicKey().Marshal()

	val3 := newValidator(t, sdk.AccAddress(valAddrs[2]), pks[2])
	val3.BlsKey = privKey3.PublicKey().Marshal()

	vals := []stakingtypes.Validator{val1, val2, val3}

	return vals, blsKeys
}
