package keeper

import (
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// GetAuthority returns the x/gashub module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetParams returns the total set of x/gashub parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the x/gashub parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetMsgGasParams get the MsgGasParams associated with a msg type url
func (k Keeper) GetMsgGasParams(ctx sdk.Context, msgTypeUrl string) types.MsgGasParams {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetMsgGasParamsKey(msgTypeUrl))
	var mgp types.MsgGasParams
	k.cdc.MustUnmarshal(b, &mgp)
	return mgp
}

// SetMsgGasParams set the provided MsgGasParams in the gashub store
func (k Keeper) SetMsgGasParams(ctx sdk.Context, mgp types.MsgGasParams) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&mgp)
	store.Set(types.GetMsgGasParamsKey(mgp.MsgTypeUrl), b)
}

// SetAllMsgGasParams set all the provided MsgGasParams in the gashub store
func (k Keeper) SetAllMsgGasParams(ctx sdk.Context, mgps []*types.MsgGasParams) {
	for _, mgp := range mgps {
		if mgp == nil {
			continue
		}
		k.SetMsgGasParams(ctx, *mgp)
	}
}

// HasMsgGasParams check existence of the MsgGasParams associated with a msg type url
func (k Keeper) HasMsgGasParams(ctx sdk.Context, msgTypeUrl string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetMsgGasParamsKey(msgTypeUrl))
}

// DeleteMsgGasParams delete the MsgGasParams associated with provided msg type urls
func (k Keeper) DeleteMsgGasParams(ctx sdk.Context, msgTypeUrls ...string) {
	store := ctx.KVStore(k.storeKey)
	for _, msgTypeUrl := range msgTypeUrls {
		store.Delete(types.GetMsgGasParamsKey(msgTypeUrl))
	}
}

// IterateMsgGasParams iterate over msg types
func (k Keeper) IterateMsgGasParams(ctx sdk.Context, handler func(msgTypeUrl string, mgp *types.MsgGasParams) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.MsgGasParamsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var mgp types.MsgGasParams
		k.cdc.MustUnmarshal(iter.Value(), &mgp)
		msgTypeUrl := types.GetMsgTypeUrl(iter.Key())
		if handler(msgTypeUrl, &mgp) {
			break
		}
	}
}

// GetAllMsgGasParams get all MsgGasParams
func (k Keeper) GetAllMsgGasParams(ctx sdk.Context) (mgps []*types.MsgGasParams) {
	k.IterateMsgGasParams(ctx, func(msgTypeUrl string, mgp *types.MsgGasParams) bool {
		mgps = append(mgps, mgp)
		return false
	})

	return mgps
}
