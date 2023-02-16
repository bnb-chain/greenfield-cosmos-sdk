package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// SetParams sets the gashub module's parameters.
func (ghk Keeper) SetParams(ctx sdk.Context, params types.Params) {
	ghk.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the gashub module's parameters.
func (ghk Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	ghk.paramSubspace.GetParamSet(ctx, &params)
	return
}
