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

// ValidateGenesis validates the cross chain genesis parameters
func ValidateGenesis(data GenesisState) error {
	if data.Params.InitModuleBalance.IsNil() || !data.Params.InitModuleBalance.IsPositive() {
		return fmt.Errorf("init module balance should be positive, is %s", data.Params.InitModuleBalance.String())
	}

	return nil
}
