package types

import "github.com/cosmos/cosmos-sdk/types"

type (
	FeeCalculator          func(msg types.Msg) uint64
	FeeCalculatorGenerator func(params Params) FeeCalculator
)

var (
	calculators    = make(map[string]FeeCalculator)
	CalculatorsGen = make(map[string]FeeCalculatorGenerator)
)

func RegisterCalculator(msgType string, feeCalc FeeCalculator) {
	calculators[msgType] = feeCalc
}

func GetCalculatorGen(msgType string) FeeCalculatorGenerator {
	return CalculatorsGen[msgType]
}

func FixedFeeCalculator(amount uint64) FeeCalculator {
	return func(msg types.Msg) uint64 {
		return amount
	}
}

var msgSendFeeCalculatorGen = func(params Params) FeeCalculator {
	msgSendGas := params.GetMsgSendGas()
	if msgSendGas == 0 {
		return FixedFeeCalculator(DefaultMsgSendGas)
	}
	return FixedFeeCalculator(msgSendGas)
}

func init() {
	CalculatorsGen["send"] = msgSendFeeCalculatorGen
}
