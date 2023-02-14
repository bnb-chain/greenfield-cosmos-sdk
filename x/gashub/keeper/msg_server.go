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
	newMsgGasParams := msg.NewParams

	msgGasParamsSet := params.MsgGasParamsSet
	typeUrl := msg.NewParams.MsgTypeUrl

	fromValue := ""
	for idx, msgGasParams := range msgGasParamsSet {
		if msgGasParams.MsgTypeUrl == typeUrl {
			fromValue = msgGasParams.String()
			msgGasParamsSet[idx] = newMsgGasParams
			break
		}
	}

	// for new msg gas type, add it to params and register gas calculator
	if fromValue == "" {
		params.MsgGasParamsSet = append(params.MsgGasParamsSet, newMsgGasParams)
		err := registerGasCalculator(newMsgGasParams)
		if err != nil {
			return nil, err
		}
	}

	k.SetParams(ctx, params)

	ctx.EventManager().EmitTypedEvent(
		&types.EventUpdateMsgGasParams{
			MsgTypeUrl: msg.NewParams.MsgTypeUrl,
			FromValue:  fromValue,
			ToValue:    newMsgGasParams.String(),
		},
	)

	return &types.MsgUpdateMsgGasParamsResponse{}, nil
}
