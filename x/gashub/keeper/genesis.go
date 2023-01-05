package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// InitGenesis - Init store state from genesis data
//
// CONTRACT: old coins from the FeeCollectionKeeper need to be transferred through
// a genesis port script to the new fee collector account
func (fhk Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	fhk.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (fhk Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := fhk.GetParams(ctx)

	return types.NewGenesisState(params)
}
