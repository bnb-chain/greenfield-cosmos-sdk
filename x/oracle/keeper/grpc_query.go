package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// InturnRelayer returns current in-turn relayer and its relaying start and end time
func (k Keeper) InturnRelayer(c context.Context, req *types.QueryInturnRelayerRequest) (*types.QueryInturnRelayerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	_, relayerInterval := k.GetRelayerParams(ctx)
	return k.GetInturnRelayer(ctx, relayerInterval, req.ClaimSrcChain)
}
