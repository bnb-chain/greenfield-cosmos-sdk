package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	GetLastValidators(ctx sdk.Context) (validators []types.Validator)
	GetHistoricalInfo(ctx sdk.Context, height int64) (types.HistoricalInfo, bool)
	BondDenom(ctx sdk.Context) (res string)
}

type CrossChainKeeper interface {
	CreateRawIBCPackageWithFee(ctx sdk.Context, channelID sdk.ChannelID,
		packageType sdk.CrossChainPackageType, packageLoad []byte, relayerFee *big.Int, ackRelayerFee *big.Int,
	) (uint64, error)
	GetCrossChainApp(channelID sdk.ChannelID) sdk.CrossChainApplication
	GetSrcChainID() sdk.ChainID
	IsDestChainSupported(chainID sdk.ChainID) bool
	GetReceiveSequence(ctx sdk.Context, channelID sdk.ChannelID) uint64
	IncrReceiveSequence(ctx sdk.Context, channelID sdk.ChannelID)
}

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
