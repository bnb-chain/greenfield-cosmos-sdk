package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

var DefaultInitModuleBalance sdkmath.Int

func init() {
	initModuleBalance, ok := sdkmath.NewIntFromString("2000000000000000000000000") // 2M
	if !ok {
		panic("invalid init module balance")
	}
	DefaultInitModuleBalance = initModuleBalance
}

func DefaultParams() Params {
	return Params{
		InitModuleBalance: DefaultInitModuleBalance,
	}
}

func (p *Params) Validate() error {
	if p.InitModuleBalance.IsNil() {
		return fmt.Errorf("init module balance should not be nil")
	}

	if !p.InitModuleBalance.IsPositive() {
		return fmt.Errorf("init module balance should be positive, is %s", p.InitModuleBalance.String())
	}

	return nil
}
