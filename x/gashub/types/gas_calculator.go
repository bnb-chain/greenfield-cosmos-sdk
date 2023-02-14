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

var MsgGrantGasCalculatorGen = func(params Params) GasCalculator {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	for _, gasParams := range msgGasParamsSet {
		if p := gasParams.GetGrantType(); p != nil {
			return GrantCalculator(p.FixedGas, p.GasPerItem)
		}
	}
	panic("no params for /cosmos.authz.v1beta1.MsgGrant")
}

var MsgMultiSendGasCalculatorGen = func(params Params) GasCalculator {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	for _, gasParams := range msgGasParamsSet {
		if p := gasParams.GetMultiSendType(); p != nil {
			return MultiSendCalculator(p.FixedGas, p.GasPerItem)
		}
	}
	panic("no params for /cosmos.bank.v1beta1.MsgMultiSend")
}

var MsgGrantAllowanceGasCalculatorGen = func(params Params) GasCalculator {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	for _, gasParams := range msgGasParamsSet {
		if p := gasParams.GetGrantAllowanceType(); p != nil {
			return GrantAllowanceCalculator(p.FixedGas, p.GasPerItem)
		}
	}
	panic("no params for /cosmos.feegrant.v1beta1.MsgGrantAllowance")
}
