package types

import (
	"fmt"

	"sigs.k8s.io/yaml"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultMaxTxSize     uint64 = 1024
	DefaultMinGasPerByte uint64 = 5
)

// Parameter keys
var (
	KeyMaxTxSize       = []byte("MaxTxSize")
	KeyMinGasPerByte   = []byte("MinGasPerByte")
	KeyMsgGasParamsSet = []byte("MsgGasParamsSet")
)

var _ paramtypes.ParamSet = &Params{}

// NewMsgGasParamsWithFixedGas creates a new MsgGasParams object with fixed gas
func NewMsgGasParamsWithFixedGas(msgTypeUrl string, gas uint64) *MsgGasParams {
	return &MsgGasParams{
		MsgTypeUrl: msgTypeUrl,
		GasParams:  &MsgGasParams_FixedType{FixedType: &MsgGasParams_FixedGasParams{FixedGas: gas}},
	}
}

// NewMsgGasParamsWithDynamicGas creates a new MsgGasParams object with dynamic gas
func NewMsgGasParamsWithDynamicGas(msgTypeUrl string, msgGasParams isMsgGasParams_GasParams) *MsgGasParams {
	return &MsgGasParams{
		MsgTypeUrl: msgTypeUrl,
		GasParams:  msgGasParams,
	}
}

// NewParams creates a new Params object
func NewParams(
	maxTxSize, minGasPerByte uint64, msgGasParamsSet []*MsgGasParams,
) Params {
	return Params{
		MaxTxSize:       maxTxSize,
		MinGasPerByte:   minGasPerByte,
		MsgGasParamsSet: msgGasParamsSet,
	}
}

// ParamKeyTable for gashub module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value
// pairs of gashub's parameters.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxTxSize, &p.MaxTxSize, validateMaxTxSize),
		paramtypes.NewParamSetPair(KeyMinGasPerByte, &p.MinGasPerByte, validateMinGasPerByte),
		paramtypes.NewParamSetPair(KeyMsgGasParamsSet, &p.MsgGasParamsSet, validateMsgGasParams),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
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
		NewMsgGasParamsWithFixedGas("/cosmos.slashing.v1beta1.MsgUnjail", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgBeginRedelegate", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCreateValidator", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgDelegate", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgEditValidator", 1e5),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgUndelegate", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.bridge.MsgTransferOut", 1e5),
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
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreatePaymentAccount", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDeposit", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgWithdraw", 1e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDisableRefund", 1e5),
		NewMsgGasParamsWithDynamicGas(
			"/cosmos.authz.v1beta1.MsgGrant",
			&MsgGasParams_GrantType{
				GrantType: &MsgGasParams_DynamicGasParams{
					FixedGas:   1e5,
					GasPerItem: 1e5,
				},
			},
		),
		NewMsgGasParamsWithDynamicGas(
			"/cosmos.bank.v1beta1.MsgMultiSend",
			&MsgGasParams_MultiSendType{
				MultiSendType: &MsgGasParams_DynamicGasParams{
					FixedGas:   1e5,
					GasPerItem: 1e5,
				},
			},
		),
		NewMsgGasParamsWithDynamicGas(
			"/cosmos.feegrant.v1beta1.MsgGrantAllowance",
			&MsgGasParams_GrantAllowanceType{
				GrantAllowanceType: &MsgGasParams_DynamicGasParams{
					FixedGas:   1e5,
					GasPerItem: 1e5,
				},
			},
		),
	}
	return Params{
		MaxTxSize:       DefaultMaxTxSize,
		MinGasPerByte:   DefaultMinGasPerByte,
		MsgGasParamsSet: defaultMsgGasParamsSet,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateMaxTxSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid max tx size: %d", v)
	}

	return nil
}

func validateMinGasPerByte(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid min gas per byte: %d", v)
	}

	return nil
}

func validateMsgGasParams(i interface{}) error {
	v, ok := i.([]*MsgGasParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	anyUnpacker := codectypes.NewInterfaceRegistry()
	RegisterInterfaces(anyUnpacker)
	for _, msgGasParams := range v {
		switch p := msgGasParams.GasParams.(type) {
		case *MsgGasParams_FixedType:
			if p.FixedType.FixedGas == 0 {
				return fmt.Errorf("invalid gas. cannot be zero")
			}
		case *MsgGasParams_GrantType:
			if p.GrantType.FixedGas == 0 || p.GrantType.GasPerItem == 0 {
				return fmt.Errorf("invalid gas. cannot be zero")
			}
		case *MsgGasParams_MultiSendType:
			if p.MultiSendType.FixedGas == 0 || p.MultiSendType.GasPerItem == 0 {
				return fmt.Errorf("invalid gas. cannot be zero")
			}
		case *MsgGasParams_GrantAllowanceType:
			if p.GrantAllowanceType.FixedGas == 0 || p.GrantAllowanceType.GasPerItem == 0 {
				return fmt.Errorf("invalid gas. cannot be zero")
			}
		default:
			return fmt.Errorf("unknown or unspecified gas type")
		}
	}

	return nil
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateMaxTxSize(p.MaxTxSize); err != nil {
		return err
	}
	if err := validateMinGasPerByte(p.MinGasPerByte); err != nil {
		return err
	}
	if err := validateMsgGasParams(p.MsgGasParamsSet); err != nil {
		return err
	}

	return nil
}
