package keeper

import (
	"context"

	"cosmossdk.io/errors"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the slashing MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// UpdateParams implements MsgServer.UpdateParams method.
// It defines a method to update the x/slashing module parameters.
func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// Unjail implements MsgServer.Unjail method.
// Validators must submit a transaction to unjail itself after
// having been jailed (and thus unbonded) for downtime
func (k msgServer) Unjail(goCtx context.Context, msg *types.MsgUnjail) (*types.MsgUnjailResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, valErr := sdk.AccAddressFromHexUnsafe(msg.ValidatorAddr)
	if valErr != nil {
		return nil, valErr
	}
	err := k.Keeper.Unjail(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.MsgUnjailResponse{}, nil
}

// Impeach defines a method for removing an existing validator after gov proposal passes.
func (k msgServer) Impeach(goCtx context.Context, msg *types.MsgImpeach) (*types.MsgImpeachResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signers := msg.GetSigners()
	if len(signers) != 1 || !signers[0].Equals(authtypes.NewModuleAddress(govtypes.ModuleName)) {
		return nil, types.ErrSignerNotGovModule
	}

	valAddr, err := sdk.AccAddressFromHexUnsafe(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// validator must already be registered
	validator := k.sk.Validator(ctx, valAddr)
	if validator == nil {
		return nil, types.ErrNoValidatorForAddress
	}

	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return nil, err
	}

	// Jail the validator if not already jailed. This will begin unbonding the
	// validator if not already unbonding (tombstoned).
	if !validator.IsJailed() {
		k.Jail(ctx, consAddr)
	}

	// Jail forever.
	k.JailForever(ctx, consAddr)

	return &types.MsgImpeachResponse{}, nil
}
