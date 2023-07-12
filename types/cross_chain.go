package types

import (
	"encoding/binary"
	"fmt"
	"math/big"
)

type CrossChainPackageType uint8

type (
	ChannelID uint8
	ChainID   uint16
)

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

type CrossChainApplication interface {
	ExecuteSynPackage(ctx Context, appCtx *CrossChainAppContext, payload []byte) ExecuteResult
	ExecuteAckPackage(ctx Context, header *CrossChainAppContext, payload []byte) ExecuteResult
	// When the ack application crash, payload is the payload of the origin package.
	ExecuteFailAckPackage(ctx Context, header *CrossChainAppContext, payload []byte) ExecuteResult
}

type CrossChainAppContext struct {
	SrcChainId ChainID
	Sequence   uint64
	Header     *PackageHeader
}

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

	SynPackageHeaderLength = 2*CrossChainFeeLength + TimestampLength + PackageTypeLength
	AckPackageHeaderLength = CrossChainFeeLength + TimestampLength + PackageTypeLength
)

func GetPackageHeaderLength(packageType CrossChainPackageType) int {
	if packageType == SynCrossChainPackageType {
		return SynPackageHeaderLength
	}
	return AckPackageHeaderLength
}

type PackageHeader struct {
	PackageType CrossChainPackageType
	Timestamp   uint64
	RelayerFee  *big.Int // relayer fee is the relayer fee paid to relayer src source chain to dest chain
	// ack relayer fee is the relayer fee paid to relayer for the ack or fail ack package if there is any
	// Ack and FailAck packages don't have ack relayer fee, since there is no corresponding ack or fail ack packages
	AckRelayerFee *big.Int
}

var NilAckRelayerFee = big.NewInt(0) // For ack packages, the ack relayer fee should be nil, and it would not be encoded into package header

func EncodePackageHeader(header PackageHeader) []byte {
	packageHeader := make([]byte, GetPackageHeaderLength(header.PackageType))
	packageHeader[0] = uint8(header.PackageType)

	timestampBytes := make([]byte, TimestampLength)
	binary.BigEndian.PutUint64(timestampBytes, header.Timestamp)
	copy(packageHeader[PackageTypeLength:PackageTypeLength+TimestampLength], timestampBytes)

	relayerFeeLength := len(header.RelayerFee.Bytes())
	copy(packageHeader[AckPackageHeaderLength-relayerFeeLength:AckPackageHeaderLength], header.RelayerFee.Bytes())

	// add ack relayer fee to header for syn package
	if header.PackageType == SynCrossChainPackageType {
		ackRelayerFeeLength := len(header.AckRelayerFee.Bytes())
		copy(packageHeader[SynPackageHeaderLength-ackRelayerFeeLength:SynPackageHeaderLength], header.AckRelayerFee.Bytes())
	}

	return packageHeader
}

func DecodePackageHeader(packageHeader []byte) (PackageHeader, error) {
	if len(packageHeader) == 0 {
		return PackageHeader{}, fmt.Errorf("empty package header")
	}

	packageType := CrossChainPackageType(packageHeader[0])
	if !IsValidCrossChainPackageType(packageType) {
		return PackageHeader{}, fmt.Errorf("package type %d is invalid", packageType)
	}

	headerLength := GetPackageHeaderLength(packageType)
	if len(packageHeader) < headerLength {
		err := fmt.Errorf("length of packageHeader is less than %d", headerLength)
		return PackageHeader{}, err
	}

	timestamp := binary.BigEndian.Uint64(packageHeader[PackageTypeLength : PackageTypeLength+TimestampLength])

	relayerFee := big.NewInt(0).SetBytes(packageHeader[PackageTypeLength+TimestampLength : AckPackageHeaderLength])

	header := PackageHeader{
		PackageType:   packageType,
		Timestamp:     timestamp,
		RelayerFee:    relayerFee,
		AckRelayerFee: big.NewInt(0),
	}

	if packageType == SynCrossChainPackageType {
		header.AckRelayerFee = big.NewInt(0).SetBytes(packageHeader[AckPackageHeaderLength:SynPackageHeaderLength])
	}

	return header, nil
}
