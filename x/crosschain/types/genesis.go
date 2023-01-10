package types

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
	return nil
}
