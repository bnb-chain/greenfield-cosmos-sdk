package types

import (
	"fmt"

	"sigs.k8s.io/yaml"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultMaxTxSize     uint64 = 32 * 1024 // 32kb
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
		paramtypes.NewParamSetPair(KeyMsgGasParamsSet, &p.MsgGasParamsSet, ValidateMsgGasParams),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	defaultMsgGasParamsSet := []*MsgGasParams{
		NewMsgGasParamsWithFixedGas("/cosmos.authz.v1beta1.MsgExec", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.authz.v1beta1.MsgRevoke", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.bank.v1beta1.MsgSend", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.feegrant.v1beta1.MsgRevokeAllowance", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgDeposit", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgSubmitProposal", 2e8),
		NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgVote", 2e7),
		NewMsgGasParamsWithFixedGas("/cosmos.gov.v1.MsgVoteWeighted", 2e7),
		NewMsgGasParamsWithFixedGas("/cosmos.oracle.v1.MsgClaim", 1e3),
		NewMsgGasParamsWithFixedGas("/cosmos.slashing.v1beta1.MsgUnjail", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgBeginRedelegate", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgCreateValidator", 2e8),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgDelegate", 1.2e3),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgEditValidator", 2e7),
		NewMsgGasParamsWithFixedGas("/cosmos.staking.v1beta1.MsgUndelegate", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.bridge.MsgTransferOut", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.sp.MsgCreateStorageProvider", 2e8),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.sp.MsgDeposit", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.sp.MsgEditStorageProvider", 2e7),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.sp.MsgUpdateSpStoragePrice", 2e7),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreateBucket", 2.4e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDeleteBucket", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgMirrorBucket", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgUpdateBucketInfo", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreateObject", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgSealObject", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgMirrorObject", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgRejectSealObject", 1.2e4),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDeleteObject", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCopyObject", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCancelCreateObject", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDiscontinueObject", 2.4e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgCreateGroup", 2.4e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDeleteGroup", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgLeaveGroup", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgUpdateGroupMember", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgMirrorGroup", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgPutPolicy", 2.4e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.storage.MsgDeletePolicy", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.payment.MsgCreatePaymentAccount", 2e5),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.payment.MsgDeposit", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.payment.MsgWithdraw", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.payment.MsgDisableRefund", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.challenge.MsgSubmit", 1.2e3),
		NewMsgGasParamsWithFixedGas("/bnbchain.greenfield.challenge.MsgAttest", 1e2),
		NewMsgGasParamsWithDynamicGas(
			"/cosmos.authz.v1beta1.MsgGrant",
			&MsgGasParams_GrantType{
				GrantType: &MsgGasParams_DynamicGasParams{
					FixedGas:   8e2,
					GasPerItem: 8e2,
				},
			},
		),
		NewMsgGasParamsWithDynamicGas(
			"/cosmos.bank.v1beta1.MsgMultiSend",
			&MsgGasParams_MultiSendType{
				MultiSendType: &MsgGasParams_DynamicGasParams{
					FixedGas:   8e2,
					GasPerItem: 8e2,
				},
			},
		),
		NewMsgGasParamsWithDynamicGas(
			"/cosmos.feegrant.v1beta1.MsgGrantAllowance",
			&MsgGasParams_GrantAllowanceType{
				GrantAllowanceType: &MsgGasParams_DynamicGasParams{
					FixedGas:   8e2,
					GasPerItem: 8e2,
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

// validateMaxTxSize performs basic validation of MaxTxSize.
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

// validateMinGasPerByte performs basic validation of MinGasPerByte.
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

// ValidateMsgGasParams performs basic validation of MsgGasParams.
func ValidateMsgGasParams(i interface{}) error {
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
	if err := ValidateMsgGasParams(p.MsgGasParamsSet); err != nil {
		return err
	}

	return nil
}
