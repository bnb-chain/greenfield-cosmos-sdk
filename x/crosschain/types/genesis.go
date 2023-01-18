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

// ValidateGenesis validates the cross chain genesis parameters
func ValidateGenesis(data GenesisState) error {
	balance, valid := big.NewInt(0).SetString(data.Params.InitModuleBalance, 10)
	if !valid {
		return fmt.Errorf("invalid module balance, is %s", data.Params.InitModuleBalance)
	}

	if balance.Cmp(big.NewInt(0)) < 0 {
		return fmt.Errorf("init module balance should be positive, is %s", balance.String())
	}

	return nil
}
