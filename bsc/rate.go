package bsc

import (
	"math"
	"math/big"
)

const (
	BNBDecimalOnBFS = 8
	BNBDecimalOnBSC = 18
)

func GetBigIntForDecimal(decimal int) *big.Int {
	floatDecimal := big.NewFloat(math.Pow10(decimal))
	bigIntDecimal := new(big.Int)
	floatDecimal.Int(bigIntDecimal)

	return bigIntDecimal
}

// ConvertBCAmountToBSCAmount can only be used to convert BNB decimal
func ConvertBCAmountToBSCAmount(bcAmount int64) *big.Int {
	decimals := GetBigIntForDecimal(BNBDecimalOnBSC - BNBDecimalOnBFS)
	return big.NewInt(0).Mul(big.NewInt(bcAmount), decimals)
}

// ConvertBSCAmountToBCAmount can only be used to convert BNB decimal
func ConvertBSCAmountToBCAmount(bscAmount *big.Int) int64 {
	decimals := GetBigIntForDecimal(BNBDecimalOnBSC - BNBDecimalOnBFS)
	return big.NewInt(0).Div(bscAmount, decimals).Int64()
}
