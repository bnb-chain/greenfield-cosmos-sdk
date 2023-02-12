package simulation

// DONTCOVER

import (
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"encoding/json"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMsgGasParamsSet),
			func(r *rand.Rand) string {
				paramsBytes, err := json.Marshal(GenMsgGasParams(r))
				if err != nil {
					panic(err)
				}
				return string(paramsBytes)
			},
		),
	}
}
