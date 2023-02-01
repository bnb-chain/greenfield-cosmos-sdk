package types

import (
	"fmt"

	"sigs.k8s.io/yaml"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultMaxTxSize                         uint64 = 1024
	DefaultMinGasPerByte                     uint64 = 5
	DefaultMsgGrantFixedGas                  uint64 = 1e5
	DefaultMsgGrantPerItemGas                uint64 = 1e5
	DefaultMsgRevokeGas                      uint64 = 1e5
	DefaultMsgExecGas                        uint64 = 1e5
	DefaultMsgSendGas                        uint64 = 1e5
	DefaultMsgMultiSendFixedGas              uint64 = 1e5
	DefaultMsgMultiSendPerItemGas            uint64 = 1e5
	DefaultMsgWithdrawDelegatorRewardGas     uint64 = 1e5
	DefaultMsgWithdrawValidatorCommissionGas uint64 = 1e5
	DefaultMsgSetWithdrawAddressGas          uint64 = 1e5
	DefaultMsgFundCommunityPoolGas           uint64 = 1e5
	DefaultMsgGrantAllowanceFixedGas         uint64 = 1e5
	DefaultMsgGrantAllowancePerItemGas       uint64 = 1e5
	DefaultMsgRevokeAllowanceGas             uint64 = 1e5
	DefaultMsgSubmitProposalGas              uint64 = 1e5
	DefaultMsgVoteGas                        uint64 = 1e5
	DefaultMsgVoteWeightedGas                uint64 = 1e5
	DefaultMsgDepositGas                     uint64 = 1e5
	DefaultMsgUnjailGas                      uint64 = 1e5
	DefaultMsgImpeachGas                     uint64 = 1e5
	DefaultMsgEditValidatorGas               uint64 = 1e5
	DefaultMsgDelegateGas                    uint64 = 1e5
	DefaultMsgUndelegateGas                  uint64 = 1e5
	DefaultMsgBeginRedelegateGas             uint64 = 1e5
	DefaultMsgCancelUnbondingDelegationGas   uint64 = 1e5
	DefaultMsgCreateValidatorGas             uint64 = 1e5
	DefaultMsgClaimGas                       uint64 = 1e5
	DefaultMsgTransferOutGas                 uint64 = 1e5
	DefaultMsgCreateStorageProviderGas       uint64 = 1e5
	DefaultMsgEditStorageProviderGas         uint64 = 1e5
	DefaultMsgSpDepositGas                   uint64 = 1e5
)

// Parameter keys
var (
	KeyMaxTxSize                         = []byte("MaxTxSize")
	KeyMinGasPerByte                     = []byte("MinGasPerByte")
	KeyMsgGrantFixedGas                  = []byte("MsgGrantFixedGas")
	KeyMsgGrantPerItemGas                = []byte("MsgGrantPerItemGas")
	KeyMsgRevokeGas                      = []byte("MsgRevokeGas")
	KeyMsgExecGas                        = []byte("MsgExecGas")
	KeyMsgSendGas                        = []byte("MsgSendGas")
	KeyMsgMultiSendFixedGas              = []byte("MsgMultiSendFixedGas")
	KeyMsgMultiSendPerItemGas            = []byte("MsgMultiSendPerItemGas")
	KeyMsgWithdrawDelegatorRewardGas     = []byte("MsgWithdrawDelegatorRewardGas")
	KeyMsgWithdrawValidatorCommissionGas = []byte("MsgWithdrawValidatorCommissionGas")
	KeyMsgSetWithdrawAddressGas          = []byte("MsgSetWithdrawAddressGas")
	KeyMsgFundCommunityPoolGas           = []byte("MsgFundCommunityPoolGas")
	KeyMsgGrantAllowanceFixedGas         = []byte("MsgGrantAllowanceFixedGas")
	KeyMsgGrantAllowancePerItemGas       = []byte("MsgGrantAllowancePerItemGas")
	KeyMsgRevokeAllowanceGas             = []byte("MsgRevokeAllowanceGas")
	KeyMsgSubmitProposalGas              = []byte("MsgSubmitProposalGas")
	KeyMsgVoteGas                        = []byte("MsgVoteGas")
	KeyMsgVoteWeightedGas                = []byte("MsgVoteWeightedGas")
	KeyMsgDepositGas                     = []byte("MsgDepositGas")
	KeyMsgUnjailGas                      = []byte("MsgUnjailGas")
	KeyMsgImpeachGas                     = []byte("MsgImpeachGas")
	KeyMsgEditValidatorGas               = []byte("MsgEditValidatorGas")
	KeyMsgDelegateGas                    = []byte("MsgDelegateGas")
	KeyMsgUndelegateGas                  = []byte("MsgUndelegateGas")
	KeyMsgBeginRedelegateGas             = []byte("MsgBeginRedelegateGas")
	KeyMsgCancelUnbondingDelegationGas   = []byte("MsgCancelUnbondingDelegationGas")
	KeyMsgCreateValidatorGas             = []byte("MsgCreateValidatorGas")
	KeyMsgClaimGas                       = []byte("MsgClaimGas")
	KeyMsgTransferOutGas                 = []byte("MsgTransferOutGas")
	KeyMsgCreateStorageProviderGas       = []byte("MsgCreateStorageProviderGas")
	KeyMsgEditStorageProviderGas         = []byte("MsgEditStorageProviderGas")
	KeyMsgSpDepositGas                   = []byte("MsgSpDepositGas")
)

var _ paramtypes.ParamSet = &Params{}

// NewParams creates a new Params object
func NewParams(
	maxTxSize, minGasPerByte, msgGrantFixedGas, msgGrantPerItemGas, msgRevokeGas, msgExecGas, msgSendGas, msgMultiSendFixedGas,
	msgMultiSendPerItemGas, msgWithdrawDelegatorRewardGas, msgWithdrawValidatorCommissionGas, msgSetWithdrawAddressGas,
	msgFundCommunityPoolGas, msgGrantAllowanceFixedGas, msgGrantAllowancePerItemGas, msgRevokeAllowanceGas, msgSubmitProposalGas,
	msgVoteGas, msgVoteWeightedGas, msgDepositGas, msgUnjailGas, msgImpeachGas, msgEditValidatorGas, msgDelegateGas,
	msgUndelegateGas, msgBeginRedelegateGas, msgCancelUnbondingDelegationGas, msgCreateValidatorGas, msgClaimGas,
	msgTransferOutGas, msgCreateStorageProviderGas, msgEditStorageProviderGas, msgSpDepositGas uint64,
) Params {
	return Params{
		MaxTxSize:                         maxTxSize,
		MinGasPerByte:                     minGasPerByte,
		MsgGrantFixedGas:                  msgGrantFixedGas,
		MsgGrantPerItemGas:                msgGrantPerItemGas,
		MsgRevokeGas:                      msgRevokeGas,
		MsgExecGas:                        msgExecGas,
		MsgSendGas:                        msgSendGas,
		MsgMultiSendFixedGas:              msgMultiSendFixedGas,
		MsgMultiSendPerItemGas:            msgMultiSendPerItemGas,
		MsgWithdrawDelegatorRewardGas:     msgWithdrawDelegatorRewardGas,
		MsgWithdrawValidatorCommissionGas: msgWithdrawValidatorCommissionGas,
		MsgSetWithdrawAddressGas:          msgSetWithdrawAddressGas,
		MsgFundCommunityPoolGas:           msgFundCommunityPoolGas,
		MsgGrantAllowanceFixedGas:         msgGrantAllowanceFixedGas,
		MsgGrantAllowancePerItemGas:       msgGrantAllowancePerItemGas,
		MsgRevokeAllowanceGas:             msgRevokeAllowanceGas,
		MsgSubmitProposalGas:              msgSubmitProposalGas,
		MsgVoteGas:                        msgVoteGas,
		MsgVoteWeightedGas:                msgVoteWeightedGas,
		MsgDepositGas:                     msgDepositGas,
		MsgUnjailGas:                      msgUnjailGas,
		MsgImpeachGas:                     msgImpeachGas,
		MsgEditValidatorGas:               msgEditValidatorGas,
		MsgDelegateGas:                    msgDelegateGas,
		MsgUndelegateGas:                  msgUndelegateGas,
		MsgBeginRedelegateGas:             msgBeginRedelegateGas,
		MsgCancelUnbondingDelegationGas:   msgCancelUnbondingDelegationGas,
		MsgCreateValidatorGas:             msgCreateValidatorGas,
		MsgClaimGas:                       msgClaimGas,
		MsgTransferOutGas:                 msgTransferOutGas,
		MsgCreateStorageProviderGas:       msgCreateStorageProviderGas,
		MsgEditStorageProviderGas:         msgEditStorageProviderGas,
		MsgSpDepositGas:                   msgSpDepositGas,
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
		paramtypes.NewParamSetPair(KeyMsgGrantFixedGas, &p.MsgGrantFixedGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgGrantPerItemGas, &p.MsgGrantPerItemGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgRevokeGas, &p.MsgRevokeGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgExecGas, &p.MsgExecGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgSendGas, &p.MsgSendGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgMultiSendFixedGas, &p.MsgMultiSendFixedGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgMultiSendPerItemGas, &p.MsgMultiSendPerItemGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgWithdrawDelegatorRewardGas, &p.MsgWithdrawDelegatorRewardGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgWithdrawValidatorCommissionGas, &p.MsgWithdrawValidatorCommissionGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgSetWithdrawAddressGas, &p.MsgSetWithdrawAddressGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgFundCommunityPoolGas, &p.MsgFundCommunityPoolGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgGrantAllowanceFixedGas, &p.MsgGrantAllowanceFixedGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgGrantAllowancePerItemGas, &p.MsgGrantAllowancePerItemGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgRevokeAllowanceGas, &p.MsgRevokeAllowanceGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgSubmitProposalGas, &p.MsgSubmitProposalGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgVoteGas, &p.MsgVoteGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgVoteWeightedGas, &p.MsgVoteWeightedGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgDepositGas, &p.MsgDepositGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgUnjailGas, &p.MsgUnjailGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgImpeachGas, &p.MsgImpeachGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgEditValidatorGas, &p.MsgEditValidatorGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgDelegateGas, &p.MsgDelegateGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgUndelegateGas, &p.MsgUndelegateGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgBeginRedelegateGas, &p.MsgBeginRedelegateGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgCancelUnbondingDelegationGas, &p.MsgCancelUnbondingDelegationGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgCreateValidatorGas, &p.MsgCreateValidatorGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgClaimGas, &p.MsgClaimGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgTransferOutGas, &p.MsgTransferOutGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgCreateStorageProviderGas, &p.MsgCreateStorageProviderGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgEditStorageProviderGas, &p.MsgEditStorageProviderGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgSpDepositGas, &p.MsgSpDepositGas, validateMsgGas),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MaxTxSize:                         DefaultMaxTxSize,
		MinGasPerByte:                     DefaultMinGasPerByte,
		MsgGrantFixedGas:                  DefaultMsgGrantFixedGas,
		MsgGrantPerItemGas:                DefaultMsgGrantPerItemGas,
		MsgRevokeGas:                      DefaultMsgRevokeGas,
		MsgExecGas:                        DefaultMsgExecGas,
		MsgSendGas:                        DefaultMsgSendGas,
		MsgMultiSendFixedGas:              DefaultMsgMultiSendFixedGas,
		MsgMultiSendPerItemGas:            DefaultMsgMultiSendPerItemGas,
		MsgWithdrawDelegatorRewardGas:     DefaultMsgWithdrawDelegatorRewardGas,
		MsgWithdrawValidatorCommissionGas: DefaultMsgWithdrawValidatorCommissionGas,
		MsgSetWithdrawAddressGas:          DefaultMsgSetWithdrawAddressGas,
		MsgFundCommunityPoolGas:           DefaultMsgFundCommunityPoolGas,
		MsgGrantAllowanceFixedGas:         DefaultMsgGrantAllowanceFixedGas,
		MsgGrantAllowancePerItemGas:       DefaultMsgGrantAllowancePerItemGas,
		MsgRevokeAllowanceGas:             DefaultMsgRevokeAllowanceGas,
		MsgSubmitProposalGas:              DefaultMsgSubmitProposalGas,
		MsgVoteGas:                        DefaultMsgVoteGas,
		MsgVoteWeightedGas:                DefaultMsgVoteWeightedGas,
		MsgDepositGas:                     DefaultMsgDepositGas,
		MsgUnjailGas:                      DefaultMsgUnjailGas,
		MsgImpeachGas:                     DefaultMsgImpeachGas,
		MsgEditValidatorGas:               DefaultMsgEditValidatorGas,
		MsgDelegateGas:                    DefaultMsgDelegateGas,
		MsgUndelegateGas:                  DefaultMsgUndelegateGas,
		MsgBeginRedelegateGas:             DefaultMsgBeginRedelegateGas,
		MsgCancelUnbondingDelegationGas:   DefaultMsgCancelUnbondingDelegationGas,
		MsgCreateValidatorGas:             DefaultMsgCreateValidatorGas,
		MsgClaimGas:                       DefaultMsgClaimGas,
		MsgTransferOutGas:                 DefaultMsgTransferOutGas,
		MsgCreateStorageProviderGas:       DefaultMsgCreateStorageProviderGas,
		MsgEditStorageProviderGas:         DefaultMsgEditStorageProviderGas,
		MsgSpDepositGas:                   DefaultMsgSpDepositGas,
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

func validateMsgGas(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid msg gas: %d", v)
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
	if err := validateMsgGas(p.MsgGrantFixedGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgGrantPerItemGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgRevokeGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgExecGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgSendGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgMultiSendFixedGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgMultiSendPerItemGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgWithdrawDelegatorRewardGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgWithdrawValidatorCommissionGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgSetWithdrawAddressGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgFundCommunityPoolGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgGrantAllowanceFixedGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgGrantAllowancePerItemGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgRevokeAllowanceGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgSubmitProposalGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgVoteGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgVoteWeightedGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgDepositGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgUnjailGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgImpeachGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgEditValidatorGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgDelegateGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgUndelegateGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgBeginRedelegateGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgCancelUnbondingDelegationGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgCreateValidatorGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgClaimGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgTransferOutGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgCreateStorageProviderGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgEditStorageProviderGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgSpDepositGas); err != nil {
		return err
	}

	return nil
}
