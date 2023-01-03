package types

import (
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

var msgGrantGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgGrantGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgGrantGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgRevokeGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgRevokeGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgRevokeGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgExecGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgExecGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgExecGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgSendGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgSendGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgSendGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgMultiSendGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgMultiSendGas()
	if msgGas == 0 {
		return MultiSendCalculator(DefaultMsgSendGas)
	}
	return MultiSendCalculator(msgGas)
}

var msgWithdrawDelegatorRewardGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgWithdrawDelegatorRewardGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgWithdrawDelegatorRewardGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgWithdrawValidatorCommissionGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgWithdrawValidatorCommissionGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgWithdrawValidatorCommissionGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgSetWithdrawAddressGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgSetWithdrawAddressGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgSetWithdrawAddressGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgFundCommunityPoolGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgFundCommunityPoolGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgFundCommunityPoolGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgGrantAllowanceGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgGrantAllowanceGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgGrantAllowanceGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgRevokeAllowanceGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgRevokeAllowanceGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgRevokeAllowanceGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgSubmitProposalGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgSubmitProposalGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgSubmitProposalGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgVoteGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgVoteGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgVoteGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgVoteWeightedGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgVoteWeightedGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgVoteWeightedGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgDepositGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgDepositGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgDepositGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgUnjailGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgUnjailGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgUnjailGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgImpeachGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgImpeachGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgImpeachGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgEditValidatorGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgEditValidatorGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgEditValidatorGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgDelegateGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgDelegateGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgDelegateGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgUndelegateGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgUndelegateGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgUndelegateGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgBeginRedelegateGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgBeginRedelegateGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgBeginRedelegateGas)
	}
	return FixedGasCalculator(msgGas)
}

var msgCancelUnbondingDelegationGasCalculatorGen = func(params Params) GasCalculator {
	msgGas := params.GetMsgCancelUnbondingDelegationGas()
	if msgGas == 0 {
		return FixedGasCalculator(DefaultMsgCancelUnbondingDelegationGas)
	}
	return FixedGasCalculator(msgGas)
}

func init() {
	RegisterCalculatorGen(types.MsgTypeURL(&authz.MsgGrant{}), msgGrantGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&authz.MsgRevoke{}), msgRevokeGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&authz.MsgExec{}), msgExecGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&bank.MsgSend{}), msgSendGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&bank.MsgMultiSend{}), msgMultiSendGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgWithdrawDelegatorReward{}), msgWithdrawDelegatorRewardGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgWithdrawValidatorCommission{}), msgWithdrawValidatorCommissionGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgSetWithdrawAddress{}), msgSetWithdrawAddressGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&distribution.MsgFundCommunityPool{}), msgFundCommunityPoolGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&feegrant.MsgGrantAllowance{}), msgGrantAllowanceGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&feegrant.MsgRevokeAllowance{}), msgRevokeAllowanceGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgSubmitProposal{}), msgSubmitProposalGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgVote{}), msgVoteGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgVoteWeighted{}), msgVoteWeightedGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&gov.MsgDeposit{}), msgDepositGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&slashing.MsgUnjail{}), msgUnjailGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&slashing.MsgImpeach{}), msgImpeachGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgEditValidator{}), msgEditValidatorGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgDelegate{}), msgDelegateGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgUndelegate{}), msgUndelegateGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgBeginRedelegate{}), msgBeginRedelegateGasCalculatorGen)
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgCancelUnbondingDelegation{}), msgCancelUnbondingDelegationGasCalculatorGen)
}
