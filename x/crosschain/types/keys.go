package types

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var ParamsKey = []byte{0x01}

const (
	ModuleName = "crosschain"
	StoreKey   = ModuleName

	prefixLength      = 1
	srcChainIdLength  = 2
	destChainIDLength = 2
	channelIDLength   = 1
	sequenceLength    = 8

	totalPackageKeyLength = prefixLength + srcChainIdLength + destChainIDLength + channelIDLength + sequenceLength

	MaxSideChainIdLength = 20
	SequenceLength       = 8
)

var (
	PrefixForIbcPackageKey = []byte{0x00}

	PrefixForSendSequenceKey    = []byte{0xf0}
	PrefixForReceiveSequenceKey = []byte{0xf1}

	PrefixForChannelPermissionKey = []byte{0xc0}
)

func BuildCrossChainPackageKey(srcChainID, destChainID sdk.ChainID, channelID sdk.ChannelID, sequence uint64) []byte {
	key := make([]byte, totalPackageKeyLength)

	copy(key[:prefixLength], PrefixForIbcPackageKey)
	binary.BigEndian.PutUint16(key[prefixLength:srcChainIdLength+prefixLength], uint16(srcChainID))
	binary.BigEndian.PutUint16(key[prefixLength+srcChainIdLength:prefixLength+srcChainIdLength+destChainIDLength], uint16(destChainID))
	copy(key[prefixLength+srcChainIdLength+destChainIDLength:], []byte{byte(channelID)})
	binary.BigEndian.PutUint64(key[prefixLength+srcChainIdLength+destChainIDLength+channelIDLength:], sequence)

	return key
}

type ChannelPermissionSetting struct {
	DestChainId string                `json:"dest_chain_id"`
	ChannelId   sdk.ChannelID         `json:"channel_id"`
	Permission  sdk.ChannelPermission `json:"permission"`
}

func (c *ChannelPermissionSetting) Check() error {
	if len(c.DestChainId) == 0 || len(c.DestChainId) > MaxSideChainIdLength {
		return fmt.Errorf("invalid dest chain id")
	}
	if c.Permission != sdk.ChannelAllow && c.Permission != sdk.ChannelForbidden {
		return fmt.Errorf("permission %d is invalid", c.Permission)
	}
	return nil
}

func BuildChannelSequenceKey(destChainID sdk.ChainID, channelID sdk.ChannelID, prefix []byte) []byte {
	key := make([]byte, prefixLength+destChainIDLength+channelIDLength)

	copy(key[:prefixLength], prefix)
	binary.BigEndian.PutUint16(key[prefixLength:prefixLength+destChainIDLength], uint16(destChainID))
	copy(key[prefixLength+destChainIDLength:], []byte{byte(channelID)})
	return key
}

func BuildChannelPermissionKey(destChainID sdk.ChainID, channelID sdk.ChannelID) []byte {
	key := make([]byte, prefixLength+destChainIDLength+channelIDLength)

	copy(key[:prefixLength], PrefixForChannelPermissionKey)
	binary.BigEndian.PutUint16(key[prefixLength:prefixLength+destChainIDLength], uint16(destChainID))
	copy(key[prefixLength+destChainIDLength:], []byte{byte(channelID)})
	return key
}
