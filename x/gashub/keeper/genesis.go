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
	initGasCalculators(data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (ghk Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := ghk.GetParams(ctx)

	return types.NewGenesisState(params)
}

func initGasCalculators(params types.Params) {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	// for fixed gas msgs
	for _, gasParams := range msgGasParamsSet {
		if _, ok := gasParams.GasParams.(*types.MsgGasParams_FixedType); !ok {
			continue
		}
		msgType := gasParams.GetMsgTypeUrl()
		types.RegisterCalculatorGen(msgType, func(params types.Params) types.GasCalculator {
			msgGasParamsSet := params.GetMsgGasParamsSet()
			for _, gasParams := range msgGasParamsSet {
				if gasParams.GetMsgTypeUrl() == msgType {
					p, ok := gasParams.GasParams.(*types.MsgGasParams_FixedType)
					if !ok {
						panic(fmt.Errorf("unpack failed for %s", msgType))
					}
					return types.FixedGasCalculator(p.FixedType.FixedGas)
				}
			}
			panic(fmt.Sprintf("no params for %s", msgType))
		})
	}

	// for dynamic gas msgs
	types.RegisterCalculatorGen("/cosmos.authz.v1beta1.MsgGrant", types.MsgGrantGasCalculatorGen)
	types.RegisterCalculatorGen("/cosmos.feegrant.v1beta1.MsgGrantAllowance", types.MsgGrantAllowanceGasCalculatorGen)
	types.RegisterCalculatorGen("/cosmos.bank.v1beta1.MsgMultiSend", types.MsgMultiSendGasCalculatorGen)
}
