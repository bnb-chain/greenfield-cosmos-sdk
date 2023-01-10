package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func DefaultParams() Params {
	return Params{}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}
