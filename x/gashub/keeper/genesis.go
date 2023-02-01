package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// InitGenesis - Init store state from genesis data
//
// CONTRACT: old coins from the FeeCollectionKeeper need to be transferred through
// a genesis port script to the new fee collector account
func (ghk Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	ghk.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (ghk Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := ghk.GetParams(ctx)

	return types.NewGenesisState(params)
}
