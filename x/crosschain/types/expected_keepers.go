package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type CrossChainKeeper interface {
	GetSrcChainID() sdk.ChainID
	IncrSendSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID)
	GetChannelSendPermission(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) sdk.ChannelPermission
	GetSendSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) uint64
}
