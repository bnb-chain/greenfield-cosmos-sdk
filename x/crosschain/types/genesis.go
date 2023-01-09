package types

import (
	"fmt"
	"math/big"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {
	relayerFee := big.NewInt(0)
	relayerFee, valid := relayerFee.SetString(data.Params.RelayerFee, 10)

	if !valid {
		return fmt.Errorf("invalid relayer fee, is %s", data.Params.RelayerFee)
	}

	if relayerFee.Cmp(big.NewInt(0)) < 0 {
		return fmt.Errorf("relayer fee should not be negative, is %s", data.Params.RelayerFee)
	}

	return nil
}
