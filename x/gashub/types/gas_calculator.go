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
			num = len(authorization.SpendLimit)
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
			num = len(feeAllowance.AllowedMessages)
		case *feegrant.PeriodicAllowance:
			spendLimit := len(feeAllowance.PeriodSpendLimit)
			canSpend := len(feeAllowance.PeriodCanSpend)
			if spendLimit > canSpend {
				num = spendLimit
			} else {
				num = canSpend
			}
		case *feegrant.BasicAllowance:
			num = len(feeAllowance.SpendLimit)
		}

		totalGas := fixedGas + uint64(num)*gasPerItem
		return totalGas, nil
	}
}

var msgGrantGasCalculatorGen = func(params Params) GasCalculator {
	fixedGas := params.GetMsgGrantFixedGas()
	gasPerItem := params.GetMsgGrantPerItemGas()
	return GrantCalculator(fixedGas, gasPerItem)
}

var msgMultiSendGasCalculatorGen = func(params Params) GasCalculator {
	fixedGas := params.GetMsgMultiSendFixedGas()
	gasPerItem := params.GetMsgMultiSendPerItemGas()
	return MultiSendCalculator(fixedGas, gasPerItem)
}

var msgGrantAllowanceGasCalculatorGen = func(params Params) GasCalculator {
	fixedGas := params.GetMsgGrantAllowanceFixedGas()
	gasPerItem := params.GetMsgGrantAllowancePerItemGas()
	return GrantAllowanceCalculator(fixedGas, gasPerItem)
}

func init() {
	RegisterCalculatorGen(types.MsgTypeURL(&authz.MsgGrant{}), msgGrantGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&authz.MsgRevoke{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgRevokeGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&authz.MsgExec{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgExecGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&bank.MsgSend{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgSendGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&bank.MsgMultiSend{}), msgMultiSendGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgWithdrawDelegatorReward{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgWithdrawDelegatorRewardGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgWithdrawValidatorCommission{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgWithdrawValidatorCommissionGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgSetWithdrawAddress{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgSetWithdrawAddressGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgFundCommunityPool{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgFundCommunityPoolGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&feegrant.MsgGrantAllowance{}), msgGrantAllowanceGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&feegrant.MsgRevokeAllowance{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgRevokeAllowanceGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgSubmitProposal{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgSubmitProposalGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgVote{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgVoteGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgVoteWeighted{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgVoteWeightedGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgDeposit{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgDepositGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&slashing.MsgUnjail{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgUnjailGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&slashing.MsgImpeach{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgImpeachGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgEditValidator{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgEditValidatorGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgDelegate{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgDelegateGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgUndelegate{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgUndelegateGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgBeginRedelegate{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgBeginRedelegateGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgCancelUnbondingDelegation{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgCancelUnbondingDelegationGas()
		return FixedGasCalculator(fixedGas)
	})
}
