package testutil

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MockCrossChainKeeper struct{}

func (m MockCrossChainKeeper) GetDestBscChainID() sdk.ChainID {
	return sdk.ChainID(0)
}

func (m MockCrossChainKeeper) CreateRawIBCPackageWithFee(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID, packageType sdk.CrossChainPackageType,
	packageLoad []byte, relayerFee, ackRelayerFee *big.Int,
) (uint64, error) {
	return 0, nil
}

func (m MockCrossChainKeeper) RegisterChannel(name string, id sdk.ChannelID, app sdk.CrossChainApplication) error {
	return nil
}

func (m MockCrossChainKeeper) GetSendSequence(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID) uint64 {
	return 0
}

func (m MockCrossChainKeeper) GetReceiveSequence(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID) uint64 {
	return 0
}

func (m MockCrossChainKeeper) IsDestChainSupported(chainID sdk.ChainID) bool {
	return true
}
