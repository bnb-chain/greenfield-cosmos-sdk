package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	GetLastValidators(ctx sdk.Context) (validators []types.Validator)
	GetHistoricalInfo(ctx sdk.Context, height int64) (types.HistoricalInfo, bool)
	BondDenom(ctx sdk.Context) (res string)
}

type CrossChainKeeper interface {
	CreateRawIBCPackage(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID,
		packageType sdk.CrossChainPackageType, packageLoad []byte) (uint64, error)
	GetCrossChainApp(channelID sdk.ChannelID) sdk.CrossChainApplication
	GetSrcChainID() sdk.ChainID
	IsDestChainSupported(chainID sdk.ChainID) bool
	GetReceiveSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) uint64
	IncrReceiveSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID)
}

type BankKeeper interface {
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}
