package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Params returns params of the cross chain module.
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// CrossChainPackage returns the specified cross chain package
func (k Keeper) CrossChainPackage(c context.Context, req *types.QueryCrossChainPackageRequest) (*types.QueryCrossChainPackageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pack, err := k.GetCrossChainPackage(ctx, sdk.ChainID(req.DestChainId), sdk.ChannelID(req.ChannelId), req.Sequence)
	if err != nil {
		return nil, err
	}
	return &types.QueryCrossChainPackageResponse{
		Package: pack,
	}, nil
}

// SendSequence returns the send sequence of the channel
func (k Keeper) SendSequence(c context.Context, req *types.QuerySendSequenceRequest) (*types.QuerySendSequenceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	sequence := k.GetSendSequence(ctx, sdk.ChainID(req.DestChainId), sdk.ChannelID(req.ChannelId))

	return &types.QuerySendSequenceResponse{
		Sequence: sequence,
	}, nil
}

// ReceiveSequence returns the receive sequence of the channel
func (k Keeper) ReceiveSequence(c context.Context, req *types.QueryReceiveSequenceRequest) (*types.QueryReceiveSequenceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	sequence := k.GetReceiveSequence(ctx, sdk.ChainID(req.DestChainId), sdk.ChannelID(req.ChannelId))

	return &types.QueryReceiveSequenceResponse{
		Sequence: sequence,
	}, nil
}
