package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CrossChainKeeper interface {
	CreateRawIBCPackageWithFee(ctx sdk.Context, channelID sdk.ChannelID, packageType sdk.CrossChainPackageType,
		packageLoad []byte, relayerFee *big.Int, ackRelayerFee *big.Int, callbackGasPrice *big.Int,
	) (uint64, error)

	RegisterChannel(name string, id sdk.ChannelID, app sdk.CrossChainApplication) error
}
