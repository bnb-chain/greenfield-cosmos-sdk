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
	registerGasCalculators(data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (ghk Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := ghk.GetParams(ctx)

	return types.NewGenesisState(params)
}

func registerGasCalculators(params types.Params) {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	for _, gasParams := range msgGasParamsSet {
		err := registerGasCalculator(gasParams)
		if err != nil {
			panic(err)
		}
	}
}

func registerGasCalculator(gasParams *types.MsgGasParams) error {
	msgType := gasParams.GetMsgTypeUrl()

	switch gasParams.GasParams.(type) {
	case *types.MsgGasParams_FixedType:
		types.RegisterCalculatorGen(msgType, func(params types.Params) types.GasCalculator {
			msgGasParamsSet := params.GetMsgGasParamsSet()
			for _, gasParams := range msgGasParamsSet {
				if gasParams.GetMsgTypeUrl() == msgType {
					p := gasParams.GetFixedType()
					if p == nil {
						panic(fmt.Errorf("get msg gas params failed for %s", msgType))
					}
					return types.FixedGasCalculator(p.FixedGas)
				}
			}
			panic(fmt.Sprintf("no params for %s", msgType))
		})
	case *types.MsgGasParams_GrantType:
		types.RegisterCalculatorGen(msgType, types.MsgGrantGasCalculatorGen)
	case *types.MsgGasParams_MultiSendType:
		types.RegisterCalculatorGen(msgType, types.MsgMultiSendGasCalculatorGen)
	case *types.MsgGasParams_GrantAllowanceType:
		types.RegisterCalculatorGen(msgType, types.MsgGrantAllowanceGasCalculatorGen)
	default:
		return fmt.Errorf("unknown gas params type")
	}

	return nil
}
