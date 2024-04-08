package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/golang/mock/gomock"
	"github.com/willf/bitset"

	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/keeper"
	"github.com/cosmos/cosmos-sdk/x/oracle/testutil"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type DummyCrossChainApp struct{}

type packUnpackTest struct {
	def      string
	unpacked interface{}
	packed   string
}

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
	newValidators, blsKeys := createValidators(s.T())

	s.stakingKeeper.EXPECT().GetHistoricalInfo(gomock.Any(), gomock.Any()).Return(stakingtypes.HistoricalInfo{
		Header: s.ctx.BlockHeader(),
		Valset: newValidators,
	}, true).AnyTimes()

	s.crossChainKeeper.EXPECT().GetSrcChainID().Return(sdk.ChainID(1)).AnyTimes()
	s.crossChainKeeper.EXPECT().IsDestChainSupported(sdk.ChainID(56)).Return(true).AnyTimes()
	s.crossChainKeeper.EXPECT().GetReceiveSequence(gomock.Any(), gomock.Any(), types.RelayPackagesChannelId).Return(uint64(0)).AnyTimes()
	s.crossChainKeeper.EXPECT().GetReceiveSequence(gomock.Any(), gomock.Any(), sdk.ChannelID(1)).Return(uint64(0)).AnyTimes()
	s.crossChainKeeper.EXPECT().CreateRawIBCPackageWithFee(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(0), nil).AnyTimes()
	s.crossChainKeeper.EXPECT().GetCrossChainApp(sdk.ChannelID(1)).Return(&DummyCrossChainApp{}).AnyTimes()
	s.crossChainKeeper.EXPECT().IncrReceiveSequence(gomock.Any(), gomock.Any(), gomock.Any()).Return().AnyTimes()
	s.stakingKeeper.EXPECT().BondDenom(gomock.Any()).Return("BNB").AnyTimes()
	s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	validatorMap := make(map[string]int, 0)
	for idx, validator := range newValidators {
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
		FromAddress:    newValidators[0].RelayerAddress,
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
	newValidators, blsKeys := createValidators(s.T())

	s.stakingKeeper.EXPECT().GetHistoricalInfo(gomock.Any(), gomock.Any()).Return(stakingtypes.HistoricalInfo{
		Header: s.ctx.BlockHeader(),
		Valset: newValidators,
	}, true).AnyTimes()

	s.crossChainKeeper.EXPECT().GetSrcChainID().Return(sdk.ChainID(1)).AnyTimes()
	s.crossChainKeeper.EXPECT().IsDestChainSupported(sdk.ChainID(65)).Return(false)
	s.crossChainKeeper.EXPECT().GetReceiveSequence(gomock.Any(), gomock.Any(), types.RelayPackagesChannelId).Return(uint64(0)).AnyTimes()
	s.crossChainKeeper.EXPECT().GetReceiveSequence(gomock.Any(), gomock.Any(), sdk.ChannelID(1)).Return(uint64(0)).AnyTimes()
	s.crossChainKeeper.EXPECT().CreateRawIBCPackageWithFee(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(uint64(0), nil).AnyTimes()
	s.crossChainKeeper.EXPECT().GetCrossChainApp(sdk.ChannelID(1)).Return(&DummyCrossChainApp{}).AnyTimes()
	s.crossChainKeeper.EXPECT().IncrReceiveSequence(gomock.Any(), gomock.Any(), gomock.Any()).Return().AnyTimes()
	s.stakingKeeper.EXPECT().BondDenom(gomock.Any()).Return("BNB").AnyTimes()
	s.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	validatorMap := make(map[string]int, 0)
	for idx, validator := range newValidators {
		validatorMap[validator.RelayerAddress] = idx
	}

	msgClaim := types.MsgClaim{
		FromAddress:    newValidators[0].RelayerAddress,
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

	s.crossChainKeeper.EXPECT().IsDestChainSupported(sdk.ChainID(65)).Return(true).AnyTimes()

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

func (s *TestSuite) TestMultiMessageDecode() {
	msg1, _ := hexutil.Decode("000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000e35fa931a00000000000000000000000000000000000000000000000000001626218b45860000000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e149600000000000000000000000000000000000000000000000000000000000000e10200000000000000000000000000000000000000000000000000000000000000200000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e1496000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000057465737431000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	msg2, _ := hexutil.Decode("000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000e35fa931a00000000000000000000000000000000000000000000000000001626218b45860000000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e149600000000000000000000000000000000000000000000000000000000000000e10200000000000000000000000000000000000000000000000000000000000000200000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e1496000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000057465737432000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	tests := []packUnpackTest{
		{
			def:      `[{"type": "bytes[]"}]`,
			packed:   "000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000022000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000e35fa931a00000000000000000000000000000000000000000000000000001626218b45860000000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e149600000000000000000000000000000000000000000000000000000000000000e10200000000000000000000000000000000000000000000000000000000000000200000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e1496000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000005746573743100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000e35fa931a00000000000000000000000000000000000000000000000000001626218b45860000000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e149600000000000000000000000000000000000000000000000000000000000000e10200000000000000000000000000000000000000000000000000000000000000200000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e1496000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000057465737432000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			unpacked: [][]byte{msg1, msg2},
		},
	}

	for i, test := range tests {
		encb, err := hex.DecodeString(test.packed)
		s.Require().Nilf(err, "invalid hex %s: %v", test.packed, err)

		messages, err := keeper.DecodeMultiMessage(encb)
		s.Require().Nilf(err, "test %d (%v) failed: %v", i, test.def, err)

		for _, message := range messages {
			fmt.Println("message", hex.EncodeToString(message))

			channelId, msgBytes, ackRelayFee, err := keeper.DecodeMessage(message)
			s.Require().Nil(err, "unpack error")

			fmt.Println(channelId, msgBytes, ackRelayFee)
		}
	}
}

func (s *TestSuite) TestMultiAckMessageEncode() {
	msg1, _ := hexutil.Decode("000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000e35fa931a00000000000000000000000000000000000000000000000000001626218b45860000000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e149600000000000000000000000000000000000000000000000000000000000000e10200000000000000000000000000000000000000000000000000000000000000200000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e1496000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000057465737431000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	msg2, _ := hexutil.Decode("000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000e35fa931a00000000000000000000000000000000000000000000000000001626218b45860000000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e149600000000000000000000000000000000000000000000000000000000000000e10200000000000000000000000000000000000000000000000000000000000000200000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e1496000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000057465737432000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	data := [][]byte{msg1, msg2}

	_, err := keeper.EncodeMultiAckMessage(data)
	s.Require().Nil(err, "EncodeMultiAckMessage error")
}
