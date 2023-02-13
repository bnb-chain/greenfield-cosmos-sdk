package keeper

import (
	"context"
	"fmt"

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

	isNew := msg.IsNew
	params := k.GetParams(ctx)
	newMsgGasParams := msg.NewParams

	if isNew {
		params.MsgGasParamsSet = append(params.MsgGasParamsSet, newMsgGasParams)
	} else {
		msgGasParamsSet := params.MsgGasParamsSet
		typeUrl := msg.NewParams.MsgTypeUrl

		var find bool
		for idx, msgGasParams := range msgGasParamsSet {
			if msgGasParams.MsgTypeUrl == typeUrl {
				msgGasParamsSet[idx] = newMsgGasParams
				find = true
				break
			}
		}
		if !find {
			return nil, fmt.Errorf("msg type not find: %s", typeUrl)
		}
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventUpdateMsgGasParams{
			MsgTypeUrl: msg.NewParams.MsgTypeUrl,
			GasType:    uint64(msg.NewParams.GasType),
		},
	)

	return &types.MsgUpdateMsgGasParamsResponse{}, nil
}
