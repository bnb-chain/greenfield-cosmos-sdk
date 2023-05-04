package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// InitGenesis - Init store state from genesis data
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, mgh := range genState.GetMsgGasParams() {
		k.SetMsgGasParams(ctx, mgh)
	}

	// register gas calculators
	k.RegisterGasCalculators(ctx)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)

	var mghs []types.MsgGasParams
	for _, mgh := range k.GetAllMsgGasParams(ctx) {
		mghs = append(mghs, *mgh)
	}
	return types.NewGenesisState(params, mghs)
}

func registerAllGasCalculators(msgGasParamsSet []*types.MsgGasParams) {
	for _, gasParams := range msgGasParamsSet {
		err := registerSingleGasCalculator(gasParams)
		if err != nil {
			panic(err)
		}
	}
}

func registerSingleGasCalculator(gasParams *types.MsgGasParams) error {
	msgTypeUrl := gasParams.GetMsgTypeUrl()

	switch gasParams.GasParams.(type) {
	case *types.MsgGasParams_FixedType:
		types.RegisterCalculatorGen(msgTypeUrl, types.FixedGasCalculatorGen)
	case *types.MsgGasParams_GrantType:
		types.RegisterCalculatorGen(msgTypeUrl, types.MsgGrantGasCalculatorGen)
	case *types.MsgGasParams_MultiSendType:
		types.RegisterCalculatorGen(msgTypeUrl, types.MsgMultiSendGasCalculatorGen)
	case *types.MsgGasParams_GrantAllowanceType:
		types.RegisterCalculatorGen(msgTypeUrl, types.MsgGrantAllowanceGasCalculatorGen)
	default:
		return fmt.Errorf("unknown gas params type")
	}

	return nil
}
