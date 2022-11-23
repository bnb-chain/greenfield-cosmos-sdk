package types

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/tendermint/tendermint/crypto"
)

var (
	PegAccount = AccAddress(crypto.AddressHash([]byte("BFSPegAccount"))) // TODO: update if needed
)

type CrossChainPackageType uint8

type ChannelID uint8
type ChainID uint16

const (
	SynCrossChainPackageType     CrossChainPackageType = 0x00
	AckCrossChainPackageType     CrossChainPackageType = 0x01
	FailAckCrossChainPackageType CrossChainPackageType = 0x02
)

type ChannelPermission uint8

const (
	ChannelAllow     ChannelPermission = 1
	ChannelForbidden ChannelPermission = 0
)

func IsValidCrossChainPackageType(packageType CrossChainPackageType) bool {
	return packageType == SynCrossChainPackageType || packageType == AckCrossChainPackageType || packageType == FailAckCrossChainPackageType
}

func ParseChannelID(input string) (ChannelID, error) {
	channelID, err := strconv.Atoi(input)
	if err != nil {
		return ChannelID(0), err
	}
	if channelID > math.MaxInt8 || channelID < 0 {
		return ChannelID(0), fmt.Errorf("channelID must be in [0, 255]")
	}
	return ChannelID(channelID), nil
}

func ParseChainID(input string) (ChainID, error) {
	chainID, err := strconv.Atoi(input)
	if err != nil {
		return ChainID(0), err
	}
	if chainID > math.MaxUint16 || chainID < 0 {
		return ChainID(0), fmt.Errorf("cross chainID must be in [0, 65535]")
	}
	return ChainID(chainID), nil
}

type CrossChainApplication interface {
	ExecuteSynPackage(ctx Context, payload []byte, relayerFee int64) ExecuteResult
	ExecuteAckPackage(ctx Context, payload []byte) ExecuteResult
	// When the ack application crash, payload is the payload of the origin package.
	ExecuteFailAckPackage(ctx Context, payload []byte) ExecuteResult
}

// TODO: define the execute result
type ExecuteResult struct {
	Err     error
	Payload []byte
}

func (c ExecuteResult) IsOk() bool {
	return c.Err == nil
}

func (c ExecuteResult) ErrMsg() string {
	if c.Err == nil {
		return ""
	}
	return c.Err.Error()
}

const (
	CrossChainFeeLength = 32
	PackageTypeLength   = 1
	TimestampLength     = 8
	PackageHeaderLength = CrossChainFeeLength + TimestampLength + PackageTypeLength
)

func EncodePackageHeader(packageType CrossChainPackageType, timestamp uint64, relayerFee big.Int) []byte {
	packageHeader := make([]byte, PackageHeaderLength)
	packageHeader[0] = uint8(packageType)

	timestampBytes := make([]byte, TimestampLength)
	binary.BigEndian.PutUint64(timestampBytes, timestamp)
	copy(packageHeader[CrossChainFeeLength:CrossChainFeeLength+TimestampLength], timestampBytes)

	length := len(relayerFee.Bytes())
	copy(packageHeader[PackageHeaderLength-length:PackageHeaderLength], relayerFee.Bytes())
	return packageHeader
}

func DecodePackageHeader(packageHeader []byte) (packageType CrossChainPackageType, timestamp uint64, relayFee big.Int, err error) {
	if len(packageHeader) < PackageHeaderLength {
		err = fmt.Errorf("length of packageHeader is less than %d", PackageHeaderLength)
		return
	}
	packageType = CrossChainPackageType(packageHeader[0])

	timestamp = binary.BigEndian.Uint64(packageHeader[PackageTypeLength : CrossChainFeeLength+TimestampLength])

	relayFee.SetBytes(packageHeader[PackageTypeLength+TimestampLength : PackageHeaderLength])
	return
}
