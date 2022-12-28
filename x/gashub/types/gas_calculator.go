package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type (
	GasCalculator          func(msg types.Msg) uint64
	GasCalculatorGenerator func(params Params) GasCalculator
)

var calculatorsGen = make(map[string]GasCalculatorGenerator)

func RegisterCalculatorGen(msgType string, feeCalcGen GasCalculatorGenerator) {
	calculatorsGen[msgType] = feeCalcGen
}

func GetGasCalculatorGen(msgType string) GasCalculatorGenerator {
	return calculatorsGen[msgType]
}

func FixedGasCalculator(amount uint64) GasCalculator {
	return func(msg types.Msg) uint64 {
		return amount
	}
}

func MultiSendCalculator(amount uint64) GasCalculator {
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

var msgSendFeeCalculatorGen = func(params Params) GasCalculator {
	msgSendGas := params.GetMsgSendGas()
	if msgSendGas == 0 {
		return FixedGasCalculator(DefaultMsgSendGas)
	}
	return FixedGasCalculator(msgSendGas)
}

var msgMultiSendFeeCalculatorGen = func(params Params) GasCalculator {
	msgMultiSendGas := params.GetMsgMultiSendGas()
	if msgMultiSendGas == 0 {
		return MultiSendCalculator(DefaultMsgSendGas)
	}
	return MultiSendCalculator(msgMultiSendGas)
}

func init() {
	RegisterCalculatorGen(types.MsgTypeURL(&bank.MsgSend{}), msgSendFeeCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&bank.MsgMultiSend{}), msgMultiSendFeeCalculatorGen)
}
