package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var DefaultInitModuleBalance sdkmath.Int

func init() {
	initModuleBalance, ok := sdkmath.NewIntFromString("2000000000000000000000000") // 2M
	if !ok {
		panic("invalid init module balance")
	}
	DefaultInitModuleBalance = initModuleBalance
}

var KeyParamInitModuleBalance = []byte("InitModuleBalance")

func DefaultParams() Params {
	return Params{
		InitModuleBalance: DefaultInitModuleBalance,
	}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyParamInitModuleBalance, &p.InitModuleBalance, validateInitModuleBalance),
	}
}

func validateInitModuleBalance(i interface{}) error {
	v, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("init module balance should not be nil")
	}

	if !v.IsPositive() {
		return fmt.Errorf("init module balance should be positive, is %s", v.String())
	}

	return nil
}
