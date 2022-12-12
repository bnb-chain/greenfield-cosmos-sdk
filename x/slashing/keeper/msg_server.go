package keeper

import (
	"context"
	"time"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the slashing MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// Unjail implements MsgServer.Unjail method.
// Validators must submit a transaction to unjail itself after
// having been jailed (and thus unbonded) for downtime
func (k msgServer) Unjail(goCtx context.Context, msg *types.MsgUnjail) (*types.MsgUnjailResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, valErr := sdk.ValAddressFromHex(msg.ValidatorAddr)
	if valErr != nil {
		return nil, valErr
	}
	err := k.Keeper.Unjail(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr),
		),
	)

	return &types.MsgUnjailResponse{}, nil
}

// KickOut defines a method for removing an existing validator after gov proposal passes.
func (k msgServer) KickOut(goCtx context.Context, msg *types.MsgKickOut) (*types.MsgKickOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signers := msg.GetSigners()
	if len(signers) != 1 || !signers[0].Equals(k.ak.GetModuleAddress(gov.ModuleName)) {
		return nil, types.ErrSignerNotGovModule
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
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

	// Jail to a big enough time (Dec 31, 9999 - 23:59:59 GMT)
	k.JailUntil(ctx, consAddr, time.Unix(253402300799, 0))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeKickOut,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.ValidatorAddress),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From),
		),
	})

	return &types.MsgKickOutResponse{}, nil
}
