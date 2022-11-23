package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultRelayerFeeParam uint64 = 1e6 // decimal is 8
)

var (
	KeyParamRelayerFee = []byte("RelayerFee")
)

func DefaultParams() Params {
	return Params{
		RelayerFee: DefaultRelayerFeeParam,
	}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyParamRelayerFee, p.RelayerFee, validateRelayerFee),
	}
}

func validateRelayerFee(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("the syn_package_fee must be positive: %d", v)
	}

	return nil
}
