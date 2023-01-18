package types

import (
	"fmt"
	"math/big"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const DefaultInitModuleBalance string = "1000000000000000000000000000"

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
		paramtypes.NewParamSetPair(KeyParamInitModuleBalance, p.InitModuleBalance, validateInitModuleBalance),
	}
}

func validateInitModuleBalance(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	balance, valid := big.NewInt(0).SetString(v, 10)
	if !valid {
		return fmt.Errorf("invalid module balance, is %s", v)
	}

	if balance.Cmp(big.NewInt(0)) < 0 {
		return fmt.Errorf("init module balance should be positive, is %s", v)
	}

	return nil
}
