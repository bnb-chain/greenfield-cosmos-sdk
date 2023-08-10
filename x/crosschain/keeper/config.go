package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

type crossChainConfig struct {
	srcChainID     sdk.ChainID
	destBscChainId sdk.ChainID
	destOpChainId  sdk.ChainID

	nameToChannelID map[string]sdk.ChannelID
	channelIDToName map[sdk.ChannelID]string
	channelIDToApp  map[sdk.ChannelID]sdk.CrossChainApplication
}

func newCrossChainCfg() *crossChainConfig {
	config := &crossChainConfig{
		srcChainID:      0,
		destBscChainId:  0,
		destOpChainId:   0,
		nameToChannelID: make(map[string]sdk.ChannelID),
		channelIDToName: make(map[sdk.ChannelID]string),
		channelIDToApp:  make(map[sdk.ChannelID]sdk.CrossChainApplication),
	}
	return config
}
