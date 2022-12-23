package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feehub/types"
)

// SetParams sets the auth module's parameters.
func (fhk Keeper) SetParams(ctx sdk.Context, params types.Params) {
	fhk.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (fhk Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	fhk.paramSubspace.GetParamSet(ctx, &params)
	return
}
