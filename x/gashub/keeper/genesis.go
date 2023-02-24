package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// InitGenesis - Init store state from genesis data
func (ghk Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	ghk.SetParams(ctx, data.Params)

	// init gas calculators from genesis data
	registerAllGasCalculators(data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (ghk Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := ghk.GetParams(ctx)

	return types.NewGenesisState(params)
}

func registerAllGasCalculators(params types.Params) {
	msgGasParamsSet := params.GetMsgGasParamsSet()
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
		types.RegisterCalculatorGen(msgTypeUrl, types.FixedGasCalculatorGen(msgTypeUrl))
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
