package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// RandomGenesisMaxTxSizeParam creates randomized MaxTxSize
func RandomGenesisMaxTxSizeParam(r *rand.Rand) uint64 {
	// 32kb - 64kb
	return uint64(32*1024 + r.Int63n(32*1024))
}

// RandomGenesisMinGasPerByteParam creates randomized MinGasPerByte
func RandomGenesisMinGasPerByteParam(r *rand.Rand) uint64 {
	// 1 - 10
	return uint64(r.Int63n(10) + 1)
}

// RandomGenesisMsgGasParams randomized values for MsgGasParams
func RandomGenesisMsgGasParams(r *rand.Rand, msgTypes []string) []types.MsgGasParams {
	var msgGasParams []types.MsgGasParams
	for _, msgType := range msgTypes {
		msgGasParams = append(msgGasParams, types.MsgGasParams{
			MsgTypeUrl: msgType,
			GasParams: &types.MsgGasParams_FixedType{
				FixedType: &types.MsgGasParams_FixedGasParams{
					FixedGas: uint64(r.Int63n(1e5)),
				},
			},
		})
	}
	return msgGasParams
}

// RandomizedGenState generates a random GenesisState for auth
func RandomizedGenState(simState *module.SimulationState) {
	maxTxSize := RandomGenesisMaxTxSizeParam(simState.Rand)
	minGasPerByte := RandomGenesisMinGasPerByteParam(simState.Rand)
	msgGasParams := RandomGenesisMsgGasParams(simState.Rand, []string{"testMsg"})

	params := types.NewParams(maxTxSize, minGasPerByte)

	gashubGenesis := types.GenesisState{
		Params:       params,
		MsgGasParams: msgGasParams,
	}

	paramsBytes, err := json.MarshalIndent(&gashubGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated gashub parameters:\n%s\n", paramsBytes)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&gashubGenesis)
}
