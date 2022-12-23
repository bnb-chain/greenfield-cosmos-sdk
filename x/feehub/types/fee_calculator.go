package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

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

func MultiSendCalculator(amount uint64) FeeCalculator {
	return func(msg types.Msg) uint64 {
		msgMultiSend := msg.(*bank.MsgMultiSend)
		var num int
		if len(msgMultiSend.Inputs) > len(msgMultiSend.Outputs) {
			num = len(msgMultiSend.Inputs)
		} else {
			num = len(msgMultiSend.Outputs)
		}
		return uint64(num) * amount
	}
}

var msgSendFeeCalculatorGen = func(params Params) FeeCalculator {
	msgSendGas := params.GetMsgSendGas()
	if msgSendGas == 0 {
		return FixedFeeCalculator(DefaultMsgSendGas)
	}
	return FixedFeeCalculator(msgSendGas)
}

var msgMultiSendFeeCalculatorGen = func(params Params) FeeCalculator {
	msgMultiSendGas := params.GetMsgMultiSendGas()
	if msgMultiSendGas == 0 {
		return MultiSendCalculator(DefaultMsgSendGas)
	}
	return MultiSendCalculator(msgMultiSendGas)
}

func init() {
	CalculatorsGen[types.MsgTypeURL(&bank.MsgSend{})] = msgSendFeeCalculatorGen
	CalculatorsGen[types.MsgTypeURL(&bank.MsgMultiSend{})] = msgMultiSendFeeCalculatorGen
}
