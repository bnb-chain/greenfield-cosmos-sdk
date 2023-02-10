package types

import (
	"fmt"

	"sigs.k8s.io/yaml"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultMaxTxSize     uint64 = 1024
	DefaultMinGasPerByte uint64 = 5
)

var DefaultMsgGasParamsSet = []*MsgGasParams{
	{
		Msg_type_url: "/cosmos.staking.v1beta1.MsgCreateValidator",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgRejectSealObject",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgDeleteGroup",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.authz.v1beta1.MsgRevoke",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.distribution.v1beta1.MsgFundCommunityPool",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.slashing.v1beta1.MsgImpeach",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.staking.v1beta1.MsgEditValidator",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.bridge.MsgTransferOut",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgCreateObject",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.bank.v1beta1.MsgSend",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.gov.v1.MsgVote",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.gov.v1.MsgVoteWeighted",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgSealObject",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgDeleteBucket",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.slashing.v1beta1.MsgUnjail",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.sp.MsgCreateStorageProvider",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.sp.MsgDeposit",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgUpdateGroupMember",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.feegrant.v1beta1.MsgRevokeAllowance",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.gov.v1.MsgDeposit",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.sp.MsgEditStorageProvider",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgCreateBucket",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgLeaveGroup",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.oracle.v1.MsgClaim",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgCopyObject",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.gov.v1.MsgSubmitProposal",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.staking.v1beta1.MsgUndelegate",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.staking.v1beta1.MsgBeginRedelegate",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.staking.v1beta1.MsgDelegate",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgDeleteObject",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/bnbchain.greenfield.storage.MsgCreateGroup",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.authz.v1beta1.MsgExec",
		Params:       []uint64{1e5},
	},
	{
		Msg_type_url: "/cosmos.feegrant.v1beta1.MsgGrantAllowance",
		Params:       []uint64{1e5, 1e5},
	},
	{
		Msg_type_url: "/cosmos.authz.v1beta1.MsgGrant",
		Params:       []uint64{1e5, 1e5},
	},
	{
		Msg_type_url: "/cosmos.bank.v1beta1.MsgMultiSend",
		Params:       []uint64{1e5, 1e5},
	},
}

// Parameter keys
var (
	KeyMaxTxSize       = []byte("MaxTxSize")
	KeyMinGasPerByte   = []byte("MinGasPerByte")
	KeyMsgGasParamsSet = []byte("MsgGasParamsSet")
)

var _ paramtypes.ParamSet = &Params{}

// NewMsgGasParams creates a new MsgGasParams object
func NewMsgGasParams(msgTypeUrl string, params ...uint64) *MsgGasParams {
	return &MsgGasParams{
		Msg_type_url: msgTypeUrl,
		Params:       params,
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
	return Params{
		MaxTxSize:       DefaultMaxTxSize,
		MinGasPerByte:   DefaultMinGasPerByte,
		MsgGasParamsSet: DefaultMsgGasParamsSet,
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

	for _, msgGasParams := range v {
		if len(msgGasParams.Params) == 0 {
			return fmt.Errorf("params cannot be empty")
		}
		for _, param := range msgGasParams.Params {
			if param == 0 {
				return fmt.Errorf("invalid msg gas param: %d", param)
			}
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
