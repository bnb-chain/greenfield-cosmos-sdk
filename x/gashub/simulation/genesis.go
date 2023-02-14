package simulation

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
	"github.com/gogo/protobuf/jsonpb"
)

// Simulation parameter constants
const (
	MaxTxSize     = "max_tx_size"
	MinGasPerByte = "min_gas_per_byte"
	MinGasPrice   = "min_gas_price"
	MsgGas        = "msg_gas"
)

// GenMaxTxSize randomized MaxTxSize
func GenMaxTxSize(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 2500, 5000))
}

// GenMinGasPerByte randomized MinGasPerByte
func GenMinGasPerByte(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 5, 50))
}

// GenMinGasPrice randomized MinGasPrice
func GenMinGasPrice(r *rand.Rand) string {
	amount := simulation.RandIntBetween(r, 1, 10)
	return strconv.FormatInt(int64(amount), 10) + "gweibnb"
}

// GenMsgGasParams randomized msg gas consumption
func GenMsgGasParams(r *rand.Rand) *types.MsgGasParams {
	msgTypeUrl := simulation.RandStringOfLength(r, 12)
	gas := uint64(simulation.RandIntBetween(r, 1e5, 1e7))
	return types.NewMsgGasParamsWithFixedGas(msgTypeUrl, gas)
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

	var minGasPrice string
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinGasPrice, &minGasPrice, simState.Rand,
		func(r *rand.Rand) { minGasPrice = GenMinGasPrice(r) },
	)

	var msgGasParams *types.MsgGasParams
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MsgGas, &msgGasParams, simState.Rand,
		func(r *rand.Rand) { msgGasParams = GenMsgGasParams(r) },
	)

	params := types.NewParams(maxTxSize, minGasPerByte, minGasPrice, []*types.MsgGasParams{msgGasParams})

	gashubGenesis := types.NewGenesisState(params)

	cdc := jsonpb.Marshaler{Indent: " "}
	bz, err := cdc.MarshalToString(&gashubGenesis.Params)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated gashub parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gashubGenesis)
}
