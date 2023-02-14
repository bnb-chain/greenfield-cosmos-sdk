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
	for _, gasParams := range msgGasParamsSet {
		msgType := gasParams.GetMsgTypeUrl()

		switch gasParams.GasParams.(type) {
		case *types.MsgGasParams_FixedType:
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
		case *types.MsgGasParams_DynamicType:
			switch msgType {
			case "/cosmos.feegrant.v1beta1.MsgGrantAllowance":
				types.RegisterCalculatorGen(msgType, types.MsgGrantAllowanceGasCalculatorGen)
			case "/cosmos.bank.v1beta1.MsgMultiSend":
				types.RegisterCalculatorGen(msgType, types.MsgMultiSendGasCalculatorGen)
			case "/cosmos.authz.v1beta1.MsgGrant":
				types.RegisterCalculatorGen(msgType, types.MsgGrantGasCalculatorGen)
			default:
				panic("unknown dynamic msg type")
			}
		default:
			panic("unknown gas consumption type")
		}
	}
}
