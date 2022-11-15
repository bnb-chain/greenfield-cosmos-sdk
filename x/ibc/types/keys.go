package types

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "ibc"
	StoreKey     = ModuleName
	QuerierRoute = StoreKey

	prefixLength      = 1
	srcChainIdLength  = 2
	destChainIDLength = 2
	channelIDLength   = 1
	sequenceLength    = 8

	totalPackageKeyLength = prefixLength + srcChainIdLength + destChainIDLength + channelIDLength + sequenceLength

	MaxSideChainIdLength = 20
	SequenceLength       = 8

	GovChannelId = sdk.ChannelID(9)

	QueryParameters = "parameters"
)

var (
	PrefixForIbcPackageKey = []byte{0x00}

	PrefixForSendSequenceKey    = []byte{0xf0}
	PrefixForReceiveSequenceKey = []byte{0xf1}

	PrefixForChannelPermissionKey = []byte{0xc0}
)

func BuildIBCPackageKey(srcChainID, destChainID sdk.ChainID, channelID sdk.ChannelID, sequence uint64) []byte {
	key := make([]byte, totalPackageKeyLength)

	copy(key[:prefixLength], PrefixForIbcPackageKey)
	binary.BigEndian.PutUint16(key[prefixLength:srcChainIdLength+prefixLength], uint16(srcChainID))
	binary.BigEndian.PutUint16(key[prefixLength+srcChainIdLength:prefixLength+srcChainIdLength+destChainIDLength], uint16(destChainID))
	copy(key[prefixLength+srcChainIdLength+destChainIDLength:], []byte{byte(channelID)})
	binary.BigEndian.PutUint64(key[prefixLength+srcChainIdLength+destChainIDLength+channelIDLength:], sequence)

	return key
}

type ChannelPermissionSetting struct {
	SideChainId string                `json:"side_chain_id"`
	ChannelId   sdk.ChannelID         `json:"channel_id"`
	Permission  sdk.ChannelPermission `json:"permission"`
}

func (c *ChannelPermissionSetting) Check() error {
	if len(c.SideChainId) == 0 || len(c.SideChainId) > MaxSideChainIdLength {
		return fmt.Errorf("invalid side chain id")
	}
	if c.ChannelId == GovChannelId {
		return fmt.Errorf("gov channel id is forbidden to set")
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

func BuildChannelPermissionsPrefixKey(destChainID sdk.ChainID) []byte {
	key := make([]byte, prefixLength+destChainIDLength)

	copy(key[:prefixLength], PrefixForChannelPermissionKey)
	binary.BigEndian.PutUint16(key[prefixLength:prefixLength+destChainIDLength], uint16(destChainID))
	return key
}
