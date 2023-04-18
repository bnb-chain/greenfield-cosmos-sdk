package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// InitGenesis - Init store state from genesis data
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, mgh := range genState.GetMsgGasParams() {
		k.SetMsgGasParams(ctx, mgh)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)

	var mghs []types.MsgGasParams
	for _, mgh := range k.GetAllMsgGasParams(ctx) {
		mghs = append(mghs, *mgh)
	}
	return types.NewGenesisState(params, mghs)
}
