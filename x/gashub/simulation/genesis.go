package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// Simulation parameter constants
const (
	MaxTxSize     = "max_tx_size"
	MinGasPerByte = "min_gas_per_byte"
	FixedMsgGas   = "fixed_msg_gas"
)

// GenMaxTxSize randomized MaxTxSize
func GenMaxTxSize(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 2500, 5000))
}

// GenMinGasPerByte randomized MinGasPerByte
func GenMinGasPerByte(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 2500, 5000))
}

// GenMsgGas randomized msg gas consumption
func GenMsgGas(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 1e5, 1e7))
}

// RandomizedGenState generates a random GenesisState for auth
func RandomizedGenState(simState *module.SimulationState) {
	var maxTxSize uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxTxSize, &maxTxSize, simState.Rand,
		func(r *rand.Rand) { maxTxSize = GenMaxTxSize(r) },
	)

	var minGasPerByte uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinGasPerByte, &minGasPerByte, simState.Rand,
		func(r *rand.Rand) { minGasPerByte = GenMinGasPerByte(r) },
	)

	var fixedMsgGas uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, FixedMsgGas, &fixedMsgGas, simState.Rand,
		func(r *rand.Rand) { fixedMsgGas = GenMsgGas(r) },
	)

	params := types.NewParams(maxTxSize, minGasPerByte, fixedMsgGas, fixedMsgGas, fixedMsgGas, fixedMsgGas)

	gashubGenesis := types.NewGenesisState(params)

	bz, err := json.MarshalIndent(&gashubGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated gashub parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gashubGenesis)
}
