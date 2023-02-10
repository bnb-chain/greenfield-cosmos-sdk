package types

import (
	"fmt"

	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type (
	GasCalculator          func(msg types.Msg) (uint64, error)
	GasCalculatorGenerator func(params Params) GasCalculator
)

var calculatorsGen = make(map[string]GasCalculatorGenerator)

var ErrInvalidMsgGas = fmt.Errorf("msg gas param is invalid")

func RegisterCalculatorGen(msgType string, feeCalcGen GasCalculatorGenerator) {
	calculatorsGen[msgType] = feeCalcGen
}

func GetGasCalculatorGen(msgType string) GasCalculatorGenerator {
	return calculatorsGen[msgType]
}

func FixedGasCalculator(amount uint64) GasCalculator {
	return func(msg types.Msg) (uint64, error) {
		if amount == 0 {
			return 0, errors.Wrapf(ErrInvalidMsgGas, "msg type: %s", types.MsgTypeURL(msg))
		}
		return amount, nil
	}
}

func GrantCalculator(fixedGas, gasPerItem uint64) GasCalculator {
	return func(msg types.Msg) (uint64, error) {
		if fixedGas == 0 || gasPerItem == 0 {
			return 0, errors.Wrapf(ErrInvalidMsgGas, "msg type: %s", types.MsgTypeURL(msg))
		}

		msgGrant := msg.(*authz.MsgGrant)
		var num int
		authorization, err := msgGrant.GetAuthorization()
		if err != nil {
			return 0, err
		}
		if authorization, ok := authorization.(*staking.StakeAuthorization); ok {
			allowList := authorization.GetAllowList().GetAddress()
			denyList := authorization.GetDenyList().GetAddress()
			num = len(allowList) + len(denyList)
		}

		totalGas := fixedGas + uint64(num)*gasPerItem
		return totalGas, nil
	}
}

func MultiSendCalculator(fixedGas, gasPerItem uint64) GasCalculator {
	return func(msg types.Msg) (uint64, error) {
		if fixedGas == 0 || gasPerItem == 0 {
			return 0, errors.Wrapf(ErrInvalidMsgGas, "msg type: %s", types.MsgTypeURL(msg))
		}

		msgMultiSend := msg.(*bank.MsgMultiSend)
		var num int
		if len(msgMultiSend.Inputs) > len(msgMultiSend.Outputs) {
			num = len(msgMultiSend.Inputs)
		} else {
			num = len(msgMultiSend.Outputs)
		}
		totalGas := fixedGas + uint64(num)*gasPerItem
		return totalGas, nil
	}
}

func GrantAllowanceCalculator(fixedGas, gasPerItem uint64) GasCalculator {
	return func(msg types.Msg) (uint64, error) {
		if fixedGas == 0 || gasPerItem == 0 {
			return 0, errors.Wrapf(ErrInvalidMsgGas, "msg type: %s", types.MsgTypeURL(msg))
		}

		msgGrantAllowance := msg.(*feegrant.MsgGrantAllowance)
		var num int
		feeAllowance, err := msgGrantAllowance.GetFeeAllowanceI()
		if err != nil {
			return 0, err
		}
		if feeAllowance, ok := feeAllowance.(*feegrant.AllowedMsgAllowance); ok {
			num = len(feeAllowance.AllowedMessages)
		}

		totalGas := fixedGas + uint64(num)*gasPerItem
		return totalGas, nil
	}
}

var msgGrantGasCalculatorGen = func(params Params) GasCalculator {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	for _, gasParams := range msgGasParamsSet {
		if gasParams.GetMsg_type_url() == "/cosmos.authz.v1beta1.MsgGrant" {
			p := gasParams.GetParams()
			if len(p) != 2 {
				panic("wrong params for /cosmos.authz.v1beta1.MsgGrant")
			}
			fixedGas := p[0]
			gasPerItem := p[1]
			return GrantCalculator(fixedGas, gasPerItem)
		}
	}
	panic("no params for /cosmos.authz.v1beta1.MsgGrant")
}

var msgMultiSendGasCalculatorGen = func(params Params) GasCalculator {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	for _, gasParams := range msgGasParamsSet {
		if gasParams.GetMsg_type_url() == "/cosmos.bank.v1beta1.MsgMultiSend" {
			p := gasParams.GetParams()
			if len(p) != 2 {
				panic("wrong params for /cosmos.bank.v1beta1.MsgMultiSend")
			}
			fixedGas := p[0]
			gasPerItem := p[1]
			return MultiSendCalculator(fixedGas, gasPerItem)
		}
	}
	panic("no params for /cosmos.bank.v1beta1.MsgMultiSend")
}

var msgGrantAllowanceGasCalculatorGen = func(params Params) GasCalculator {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	for _, gasParams := range msgGasParamsSet {
		if gasParams.GetMsg_type_url() == "/cosmos.feegrant.v1beta1.MsgGrantAllowance" {
			p := gasParams.GetParams()
			if len(p) != 2 {
				panic("wrong params for /cosmos.feegrant.v1beta1.MsgGrantAllowance")
			}
			fixedGas := p[0]
			gasPerItem := p[1]
			return GrantAllowanceCalculator(fixedGas, gasPerItem)
		}
	}
	panic("no params for /cosmos.feegrant.v1beta1.MsgGrantAllowance")
}

func init() {
	// for fixed gas msgs
	for _, gasParams := range DefaultMsgGasParamsSet {
		if len(gasParams.GetParams()) != 1 {
			continue
		}
		msgType := gasParams.GetMsg_type_url()
		RegisterCalculatorGen(msgType, func(params Params) GasCalculator {
			msgGasParamsSet := params.GetMsgGasParamsSet()
			for _, gasParams := range msgGasParamsSet {
				if gasParams.GetMsg_type_url() == msgType {
					p := gasParams.GetParams()
					if len(p) != 1 {
						panic(fmt.Sprintf("wrong params for %s", msgType))
					}
					fixedGas := p[0]
					return FixedGasCalculator(fixedGas)
				}
			}
			panic(fmt.Sprintf("no params for %s", msgType))
		})
	}

	// for dynamic gas msgs
	RegisterCalculatorGen("/cosmos.authz.v1beta1.MsgGrant", msgGrantGasCalculatorGen)
	RegisterCalculatorGen("/cosmos.feegrant.v1beta1.MsgGrantAllowance", msgGrantAllowanceGasCalculatorGen)
	RegisterCalculatorGen("/cosmos.bank.v1beta1.MsgMultiSend", msgMultiSendGasCalculatorGen)
}
