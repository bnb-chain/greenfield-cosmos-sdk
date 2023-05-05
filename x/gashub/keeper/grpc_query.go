package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

var _ types.QueryServer = Keeper{}

// Params returns parameters of auth module
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) MsgGasParams(goCtx context.Context, req *types.QueryMsgGasParamsRequest) (*types.QueryMsgGasParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	resp := &types.QueryMsgGasParamsResponse{}
	if len(req.MsgTypeUrls) > 0 {
		for _, url := range req.MsgTypeUrls {
			if mgp, ok := k.getMsgGasParams(ctx, url); ok {
				resp.MsgGasParams = append(resp.MsgGasParams, mgp)
			}
		}
	} else {
		store := ctx.KVStore(k.storeKey)
		mgpStore := prefix.NewStore(store, types.MsgGasParamsPrefix)
		mgps, pageRes, err := query.GenericFilteredPaginate(k.cdc, mgpStore, req.Pagination, func(key []byte, result *types.MsgGasParams) (*types.MsgGasParams, error) {
			return result, nil
		}, func() *types.MsgGasParams {
			return &types.MsgGasParams{}
		})
		if err != nil {
			return nil, err
		}

		resp.MsgGasParams = mgps
		resp.Pagination = pageRes
	}

	return resp, nil
}

func (k Keeper) getMsgGasParams(ctx sdk.Context, url string) (*types.MsgGasParams, bool) {
	if has := k.HasMsgGasParams(ctx, url); !has {
		return nil, false
	}

	mgp := k.GetMsgGasParams(ctx, url)
	return &mgp, true
}
