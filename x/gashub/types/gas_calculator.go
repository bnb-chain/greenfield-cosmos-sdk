package types

import (
	"fmt"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	oracletypes "github.com/cosmos/cosmos-sdk/x/oracle/types"
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
	res, ok := calculatorsGen[msgType]
	// todo: this is a temporary default fee, remove this after all msg types are registered
	if !ok {
		res = func(params Params) GasCalculator {
			return FixedGasCalculator(1e5)
		}
	}
	return res
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
	RegisterCalculatorGen(types.MsgTypeURL(&staking.MsgCreateValidator{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgCreateValidatorGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen(types.MsgTypeURL(&oracletypes.MsgClaim{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgClaimGas()
		return FixedGasCalculator(fixedGas)
	})

	RegisterCalculatorGen(types.MsgTypeURL(&govv1beta1.MsgSubmitProposal{}), func(params Params) GasCalculator {
		fixedGas := params.GetMsgSubmitProposalGas()
		return FixedGasCalculator(fixedGas)
	})

	// these msgs are from greenfield, so the msg types need to be hard coded.
	RegisterCalculatorGen("/bnbchain.greenfield.bridge.MsgTransferOut", func(params Params) GasCalculator {
		fixedGas := params.GetMsgTransferOutGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen("/bnbchain.greenfield.sp.MsgCreateStorageProvider", func(params Params) GasCalculator {
		fixedGas := params.GetMsgCreateStorageProviderGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen("/bnbchain.greenfield.sp.MsgEditStorageProvider", func(params Params) GasCalculator {
		fixedGas := params.GetMsgEditStorageProviderGas()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen("/bnbchain.greenfield.sp.MsgDeposit", func(params Params) GasCalculator {
		fixedGas := params.GetMsgSpDepositGas()
		return FixedGasCalculator(fixedGas)
	})

	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgCreateBucket", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageCreateBucket()
		return FixedGasCalculator(fixedGas)
	})

	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgDeleteBucket", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageDeleteBucket()
		return FixedGasCalculator(fixedGas)
	})

	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgCreateObject", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageCreateObject()
		return FixedGasCalculator(fixedGas)
	})

	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgDeleteObject", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageDeleteObject()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgSealObject", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageSealObject()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgCopyObject", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageCopyObject()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgRejectSealObject", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageRejectSealObject()
		return FixedGasCalculator(fixedGas)
	})

	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgCreateGroup", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageCreateGroup()
		return FixedGasCalculator(fixedGas)
	})

	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgDeleteGroup", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageDeleteGroup()
		return FixedGasCalculator(fixedGas)
	})

	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgLeaveGroup", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageLeaveGroup()
		return FixedGasCalculator(fixedGas)
	})
	RegisterCalculatorGen("/bnbchain.greenfield.storage.MsgUpdateGroupMember", func(params Params) GasCalculator {
		fixedGas := params.GetMsgStorageUpdateGroupMember()
		return FixedGasCalculator(fixedGas)
	})
}
