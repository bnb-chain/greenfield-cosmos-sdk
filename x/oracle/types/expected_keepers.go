package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	GetLastValidators(ctx sdk.Context) (validators []types.Validator)
}

type CrossChainKeeper interface {
	CreateRawIBCPackage(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID,
		packageType sdk.CrossChainPackageType, packageLoad []byte) (uint64, error)
	GetCrossChainApp(channelID sdk.ChannelID) sdk.CrossChainApplication
	GetReceiveSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) uint64
	IncrReceiveSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID)
}

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
