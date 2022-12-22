package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/feehub/types"
)

// Simulation parameter constants
const (
	MaxTxSize     = "max_tx_size"
	MinGasPerByte = "min_gas_per_byte"
	MsgSendGas    = "msg_send_gas"
)

// GenMaxTxSize randomized MaxTxSize
func GenMaxTxSize(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 2500, 5000))
}

// GenMinGasPerByte randomized MinGasPerByte
func GenMinGasPerByte(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 2500, 5000))
}

// GenMsgSendGas randomized MsgSendGas
func GenMsgSendGas(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 1e6, 1e7))
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

	var msgSendGas uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MsgSendGas, &msgSendGas, simState.Rand,
		func(r *rand.Rand) { msgSendGas = GenMsgSendGas(r) },
	)

	params := types.NewParams(maxTxSize, minGasPerByte, msgSendGas)

	feehubGenesis := types.NewGenesisState(params)

	bz, err := json.MarshalIndent(&feehubGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated feehub parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(feehubGenesis)
}
