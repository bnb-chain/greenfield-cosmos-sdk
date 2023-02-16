package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gashub MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpdateMsgGasParams(goCtx context.Context, msg *types.MsgUpdateMsgGasParams) (*types.MsgUpdateMsgGasParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signers := msg.GetSigners()
	if len(signers) != 1 || !signers[0].Equals(authtypes.NewModuleAddress(gov.ModuleName)) {
		return nil, slashingtypes.ErrSignerNotGovModule
	}

	params := k.GetParams(ctx)
	msgGasParamsSet := params.MsgGasParamsSet
	newMsgGasParamsSet := msg.NewParamsSet
	for _, newParams := range newMsgGasParamsSet {
		if err := types.ValidateMsgGasParams(newParams); err != nil {
			return nil, err
		}

		typeUrl := newParams.MsgTypeUrl

		fromValue := ""
		for idx, msgGasParams := range msgGasParamsSet {
			if msgGasParams.MsgTypeUrl == typeUrl {
				fromValue = msgGasParams.String()
				msgGasParamsSet[idx] = newParams
				break
			}
		}
		if fromValue == "" {
			params.MsgGasParamsSet = append(params.MsgGasParamsSet, newParams)
		}

		// register gas calculator
		if err := registerSingleGasCalculator(newParams); err != nil {
			return nil, err
		}

		ctx.EventManager().EmitTypedEvent(
			&types.EventUpdateMsgGasParams{
				MsgTypeUrl: newParams.MsgTypeUrl,
				FromValue:  fromValue,
				ToValue:    newParams.String(),
			},
		)
	}

	k.SetParams(ctx, params)

	return &types.MsgUpdateMsgGasParamsResponse{}, nil
}
