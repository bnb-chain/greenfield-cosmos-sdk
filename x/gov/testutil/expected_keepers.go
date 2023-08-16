// This file only used to generate mocks

package testutil

import (
	"math/big"

	math "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

// AccountKeeper extends gov's actual expected AccountKeeper with additional
// methods used in tests.
type AccountKeeper interface {
	types.AccountKeeper

	IterateAccounts(ctx sdk.Context, cb func(account authtypes.AccountI) (stop bool))
}

// BankKeeper extends gov's actual expected BankKeeper with additional
// methods used in tests.
type BankKeeper interface {
	bankkeeper.Keeper
}

// StakingKeeper extends gov's actual expected StakingKeeper with additional
// methods used in tests.
type StakingKeeper interface {
	types.StakingKeeper

	BondDenom(ctx sdk.Context) string
	TokensFromConsensusPower(ctx sdk.Context, power int64) math.Int
}

// CrossChainKeeper defines the expected crossChain keeper
type CrossChainKeeper interface {
	GetDestBscChainID() sdk.ChainID
	CreateRawIBCPackageWithFee(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID, packageType sdk.CrossChainPackageType,
		packageLoad []byte, relayerFee, ackRelayerFee *big.Int,
	) (uint64, error)

	RegisterChannel(name string, id sdk.ChannelID, app sdk.CrossChainApplication) error

	GetSendSequence(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID) uint64

	GetReceiveSequence(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID) uint64

	IsDestChainSupported(chainID sdk.ChainID) bool
}
