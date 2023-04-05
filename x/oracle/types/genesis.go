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

// DefaultGenesisState - default GenesisState
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

	if data.Params.RelayerRewardShare <= 0 {
		return fmt.Errorf("the relayer reward share must be positive, is %d", data.Params.RelayerRewardShare)
	}

	if data.Params.RelayerRewardShare > 100 {
		return fmt.Errorf("the relayer reward share should not be larger than 100, is %d", data.Params.RelayerRewardShare)
	}

	if data.Params.RelayerInterval <= 0 {
		return fmt.Errorf("the relayer interval should be positive, is %d", data.Params.RelayerInterval)
	}
	return nil
}
