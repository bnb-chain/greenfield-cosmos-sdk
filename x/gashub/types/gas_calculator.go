package types

import (
	"fmt"

	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
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
		switch authorization := authorization.(type) {
		case *staking.StakeAuthorization:
			allowList := authorization.GetAllowList().GetAddress()
			denyList := authorization.GetDenyList().GetAddress()
			num = len(allowList) + len(denyList)
		case *bank.SendAuthorization:
			coins := authorization.SpendLimit
			num = len(coins)
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
		switch feeAllowance := feeAllowance.(type) {
		case *feegrant.AllowedMsgAllowance:
			allowedMsg := feeAllowance.AllowedMessages
			num = len(allowedMsg)
		}

		totalGas := fixedGas + uint64(num)*gasPerItem
		return totalGas, nil
	}
}

var fixedMsgGasCalculatorGen = func(params Params) GasCalculator {
	fixedGas := params.GetFixedMsgGas()
	return FixedGasCalculator(fixedGas)
}

var msgGrantGasCalculatorGen = func(params Params) GasCalculator {
	fixedGas := params.GetFixedMsgGas()
	gasPerItem := params.GetMsgGrantPerItemGas()
	return GrantCalculator(fixedGas, gasPerItem)
}

var msgMultiSendGasCalculatorGen = func(params Params) GasCalculator {
	fixedGas := params.GetFixedMsgGas()
	gasPerItem := params.GetMsgMultiSendPerItemGas()
	return MultiSendCalculator(fixedGas, gasPerItem)
}

var msgGrantAllowanceGasCalculatorGen = func(params Params) GasCalculator {
	fixedGas := params.GetFixedMsgGas()
	gasPerItem := params.GetMsgGrantAllowancePerItemGas()
	return GrantAllowanceCalculator(fixedGas, gasPerItem)
}

func init() {
	RegisterCalculatorGen(types.MsgTypeURL(&authz.MsgGrant{}), msgGrantGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&authz.MsgRevoke{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&authz.MsgExec{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&bank.MsgSend{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&bank.MsgMultiSend{}), msgMultiSendGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgWithdrawDelegatorReward{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgWithdrawValidatorCommission{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgSetWithdrawAddress{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgFundCommunityPool{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&feegrant.MsgGrantAllowance{}), msgGrantAllowanceGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&feegrant.MsgRevokeAllowance{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgSubmitProposal{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgVote{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgVoteWeighted{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgDeposit{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&slashing.MsgUnjail{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&slashing.MsgImpeach{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgEditValidator{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgDelegate{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgUndelegate{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgBeginRedelegate{}), fixedMsgGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgCancelUnbondingDelegation{}), fixedMsgGasCalculatorGen)
}
