package keeper_test

import (
	"math/big"
	"time"

	"github.com/willf/bitset"

	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/testutil"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type DummyCrossChainApp struct{}

func (ta *DummyCrossChainApp) ExecuteSynPackage(ctx sdk.Context, header *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	return sdk.ExecuteResult{}
}

func (ta *DummyCrossChainApp) ExecuteAckPackage(ctx sdk.Context, header *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	return sdk.ExecuteResult{}
}

func (ta *DummyCrossChainApp) ExecuteFailAckPackage(ctx sdk.Context, header *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	return sdk.ExecuteResult{}
}

func (s *TestSuite) TestClaim() {
	s.app.CrossChainKeeper.RegisterChannel("test", sdk.ChannelID(1), &DummyCrossChainApp{})

	s.app.OracleKeeper.SetParams(s.ctx, types.Params{
		RelayerTimeout:     5,
		RelayerRewardShare: 50,
		RelayerInterval:    600,
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

	payloadHeader := sdk.EncodePackageHeader(sdk.PackageHeader{
		PackageType:   sdk.SynCrossChainPackageType,
		Timestamp:     1992,
		RelayerFee:    big.NewInt(1),
		AckRelayerFee: big.NewInt(1),
	})

	testPackage := types.Package{
		ChannelId: 1,
		Sequence:  0,
		Payload:   append(payloadHeader, []byte("test payload")...),
	}

	packageBytes, err := rlp.EncodeToBytes([]types.Package{testPackage})
	s.Require().Nil(err, "encode package error")

	msgClaim := types.MsgClaim{
		FromAddress:    validators[0].RelayerAddress,
		SrcChainId:     56,
		DestChainId:    1,
		Sequence:       0,
		Timestamp:      1992,
		Payload:        packageBytes,
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
	_, err = s.msgServer.Claim(s.ctx, &msgClaim)
	s.Require().Nil(err, "process claim msg error")
}

func (s *TestSuite) TestInvalidClaim() {
	s.app.CrossChainKeeper.RegisterChannel("test", sdk.ChannelID(1), &DummyCrossChainApp{})

	s.app.OracleKeeper.SetParams(s.ctx, types.Params{
		RelayerTimeout:     5,
		RelayerRewardShare: 50,
		RelayerInterval:    600,
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
		SrcChainId:     65,
		DestChainId:    1,
		Sequence:       0,
		Timestamp:      1992,
		Payload:        []byte("invalid payload"),
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

	// invalid src chain id
	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, err := s.msgServer.Claim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "process claim should return error")
	s.Require().Contains(err.Error(), "src chain id is invalid")

	s.app.CrossChainKeeper.SetDestChainID(sdk.ChainID(65))

	// invalid payload
	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, err = s.msgServer.Claim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "process claim should return error")
	s.Require().Contains(err.Error(), "decode payload error")

	// invalid timestamp
	payloadHeader := sdk.EncodePackageHeader(sdk.PackageHeader{
		PackageType:   sdk.SynCrossChainPackageType,
		Timestamp:     1993,
		RelayerFee:    big.NewInt(1),
		AckRelayerFee: big.NewInt(1),
	})
	testPackage := types.Package{
		ChannelId: 1,
		Sequence:  0,
		Payload:   append(payloadHeader, []byte("test payload")...),
	}

	packageBytes, err := rlp.EncodeToBytes([]types.Package{testPackage})
	s.Require().Nil(err, "encode package error")

	msgClaim.Payload = packageBytes
	blsSignBytes = msgClaim.GetBlsSignBytes()
	blsSig = testutil.GenerateBlsSig(blsKeys, blsSignBytes[:])
	msgClaim.AggSignature = blsSig

	s.ctx = s.ctx.WithBlockTime(time.Unix(int64(msgClaim.Timestamp), 0))
	_, err = s.msgServer.Claim(s.ctx, &msgClaim)
	s.Require().NotNil(err, "process claim should return error")
	s.Require().Contains(err.Error(), "is not the same in payload header")
}
