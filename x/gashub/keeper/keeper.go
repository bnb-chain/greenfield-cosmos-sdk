package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper encodes/decodes accounts using the go-amino (binary)
// encoding/decoding library.
type Keeper struct {
	key           storetypes.StoreKey
	cdc           codec.BinaryCodec
	paramSubspace paramtypes.Subspace
}

// NewKeeper returns a new gashub keeper
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramstore paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		key:           key,
		cdc:           cdc,
		paramSubspace: paramstore,
	}
}

// Logger returns a module-specific logger.
func (ghk Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetCodec return codec.Codec object used by the keeper
func (ghk Keeper) GetCodec() codec.BinaryCodec { return ghk.cdc }

func (ghk Keeper) RegisterGasCalculators(ctx sdk.Context) {
	params := ghk.GetParams(ctx)
	registerGasCalculators(params)
}
