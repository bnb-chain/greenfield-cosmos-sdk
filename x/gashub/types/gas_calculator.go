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
		if gasParams.GetMsgTypeUrl() == "/cosmos.authz.v1beta1.MsgGrant" {
			p, ok := gasParams.GasParams.(*MsgGasParams_DynamicType)
			if !ok {
				panic("type conversion failed for /cosmos.authz.v1beta1.MsgGrant")
			}
			return GrantCalculator(p.DynamicType.FixedGas, p.DynamicType.GasPerItem)
		}
	}
	panic("no params for /cosmos.authz.v1beta1.MsgGrant")
}

var msgMultiSendGasCalculatorGen = func(params Params) GasCalculator {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	for _, gasParams := range msgGasParamsSet {
		if gasParams.GetMsgTypeUrl() == "/cosmos.bank.v1beta1.MsgMultiSend" {
			p, ok := gasParams.GasParams.(*MsgGasParams_DynamicType)
			if !ok {
				panic("type conversion failed for /cosmos.bank.v1beta1.MsgMultiSend")
			}
			return MultiSendCalculator(p.DynamicType.FixedGas, p.DynamicType.GasPerItem)
		}
	}
	panic("no params for /cosmos.bank.v1beta1.MsgMultiSend")
}

var msgGrantAllowanceGasCalculatorGen = func(params Params) GasCalculator {
	msgGasParamsSet := params.GetMsgGasParamsSet()
	for _, gasParams := range msgGasParamsSet {
		if gasParams.GetMsgTypeUrl() == "/cosmos.feegrant.v1beta1.MsgGrantAllowance" {
			p, ok := gasParams.GasParams.(*MsgGasParams_DynamicType)
			if !ok {
				panic("type conversion failed for /cosmos.feegrant.v1beta1.MsgGrantAllowance")
			}
			return MultiSendCalculator(p.DynamicType.FixedGas, p.DynamicType.GasPerItem)
		}
	}
	panic("no params for /cosmos.feegrant.v1beta1.MsgGrantAllowance")
}

func init() {
	defaultMsgGasParamsSet := []*MsgGasParams{
		NewMsgGasParamsWithFixedGas("/cosmos.authz.v1beta1.MsgExec", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.authz.v1beta1.MsgRevoke", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.bank.v1beta1.MsgSend", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgFundCommunityPool", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.feegrant.v1beta1.MsgRevokeAllowance", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgDeposit", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgSubmitProposal", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgVote", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgVoteWeighted", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.oracle.v1.MsgClaim", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.slashing.v1beta1.MsgImpeach", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.slashing.v1beta1.MsgUnjail", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgBeginRedelegate", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCreateValidator", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgDelegate", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgEditValidator", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgUndelegate", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.bridge.MsgTransferOut", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.sp.MsgCreateStorageProvider", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.sp.MsgDeposit", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.sp.MsgEditStorageProvider", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCopyObject", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreateBucket", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreateGroup", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreateObject", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDeleteBucket", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDeleteGroup", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgLeaveGroup", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgRejectSealObject", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgSealObject", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgUpdateGroupMember", 1e5),
		NewMsgGasParamsWithDynamicGas("/cosmos.authz.v1beta1.MsgGrant", 1e5, 1e5),
		NewMsgGasParamsWithDynamicGas("/cosmos.bank.v1beta1.MsgMultiSend", 1e5, 1e5),
		NewMsgGasParamsWithDynamicGas("/cosmos.feegrant.v1beta1.MsgGrantAllowance", 1e5, 1e5),
	}
	// for fixed gas msgs
	for _, gasParams := range defaultMsgGasParamsSet {
		if gasParams.GetGasType() != GasType_GAS_TYPE_FIXED {
			continue
		}
		msgType := gasParams.GetMsgTypeUrl()
		RegisterCalculatorGen(msgType, func(params Params) GasCalculator {
			msgGasParamsSet := params.GetMsgGasParamsSet()
			for _, gasParams := range msgGasParamsSet {
				if gasParams.GetMsgTypeUrl() == msgType {
					p, ok := gasParams.GasParams.(*MsgGasParams_FixedType)
					if !ok {
						panic(fmt.Errorf("unpack failed for %s", msgType))
					}
					return FixedGasCalculator(p.FixedType.FixedGas)
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
