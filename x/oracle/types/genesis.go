package types

import (
	"fmt"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {
	if data.Params.RelayerTimeout <= 0 {
		return fmt.Errorf("relayer timeout should be positive, is %d", data.Params.RelayerTimeout)
	}

	if data.Params.RelayerBackoffTime <= 0 {
		return fmt.Errorf("the relayer backoff time must be positive, is %d", data.Params.RelayerBackoffTime)
	}
	return nil
}
