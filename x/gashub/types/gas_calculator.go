package types

import (
	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/gashub/errors"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type (
	GasCalculator          func(msg types.Msg) (uint64, error)
	GasCalculatorGenerator func(mgp MsgGasParams) GasCalculator
)

var (
	calculatorsGen = make(map[string]GasCalculatorGenerator)

	FixedGasCalculatorGen = func(mgh MsgGasParams) GasCalculator {
		if fixedTyp := mgh.GetFixedType(); fixedTyp != nil {
			return FixedGasCalculator(fixedTyp.FixedGas)
		}
		return nil
	}

	MsgGrantGasCalculatorGen = func(mgh MsgGasParams) GasCalculator {
		if grantTyp := mgh.GetGrantType(); grantTyp != nil {
			return GrantCalculator(grantTyp.FixedGas, grantTyp.GasPerItem)
		}
		return nil
	}

	MsgMultiSendGasCalculatorGen = func(mgh MsgGasParams) GasCalculator {
		if multiSendTyp := mgh.GetMultiSendType(); multiSendTyp != nil {
			return MultiSendCalculator(multiSendTyp.FixedGas, multiSendTyp.GasPerItem)
		}
		return nil
	}

	MsgGrantAllowanceGasCalculatorGen = func(mgh MsgGasParams) GasCalculator {
		if grantAllowanceTyp := mgh.GetGrantAllowanceType(); grantAllowanceTyp != nil {
			return GrantAllowanceCalculator(grantAllowanceTyp.FixedGas, grantAllowanceTyp.GasPerItem)
		}
		return nil
	}
)

func RegisterCalculatorGen(msgType string, feeCalcGen GasCalculatorGenerator) {
	calculatorsGen[msgType] = feeCalcGen
}

func GetGasCalculatorGen(msgType string) GasCalculatorGenerator {
	return calculatorsGen[msgType]
}

func FixedGasCalculator(amount uint64) GasCalculator {
	return func(msg types.Msg) (uint64, error) {
		if amount == 0 {
			return 0, errorsmod.Wrapf(errors.ErrInvalidMsgGasParams, "msg type: %s", types.MsgTypeURL(msg))
		}
		return amount, nil
	}
}

func GrantCalculator(fixedGas, gasPerItem uint64) GasCalculator {
	return func(msg types.Msg) (uint64, error) {
		if fixedGas == 0 || gasPerItem == 0 {
			return 0, errorsmod.Wrapf(errors.ErrInvalidMsgGasParams, "msg type: %s", types.MsgTypeURL(msg))
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
			return 0, errorsmod.Wrapf(errors.ErrInvalidMsgGasParams, "msg type: %s", types.MsgTypeURL(msg))
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
			return 0, errorsmod.Wrapf(errors.ErrInvalidMsgGasParams, "msg type: %s", types.MsgTypeURL(msg))
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
