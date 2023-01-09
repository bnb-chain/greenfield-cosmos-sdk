package types

import (
	"fmt"
	"math/big"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultRelayerFeeParam string = "1"
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
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	relayerFee := big.NewInt(0)
	relayerFee, valid := relayerFee.SetString(v, 10)

	if !valid {
		return fmt.Errorf("invalid relayer fee, %s", v)
	}

	if relayerFee.Cmp(big.NewInt(0)) < 0 {
		return fmt.Errorf("invalid relayer fee, %s", v)
	}

	return nil
}
