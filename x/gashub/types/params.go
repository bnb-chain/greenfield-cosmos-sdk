package types

import (
	"fmt"

	"sigs.k8s.io/yaml"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultMaxTxSize                   uint64 = 1024
	DefaultMinGasPerByte               uint64 = 5
	DefaultFixedMsgGas                 uint64 = 1e5
	DefaultMsgGrantPerItemGas          uint64 = 1e5
	DefaultMsgMultiSendPerItemGas      uint64 = 1e5
	DefaultMsgGrantAllowancePerItemGas uint64 = 1e5
)

// Parameter keys
var (
	KeyMaxTxSize                   = []byte("MaxTxSize")
	KeyMinGasPerByte               = []byte("MinGasPerByte")
	KeyFixedMsgGas                 = []byte("FixedMsgGas")
	KeyMsgGrantPerItemGas          = []byte("MsgGrantPerItemGas")
	KeyMsgMultiSendPerItemGas      = []byte("MsgMultiSendPerItemGas")
	KeyMsgGrantAllowancePerItemGas = []byte("MsgGrantAllowancePerItemGas")
)

var _ paramtypes.ParamSet = &Params{}

// NewParams creates a new Params object
func NewParams(
	maxTxSize, minGasPerByte, fixedMsgGas, msgGrantPerItemGas, msgMultiSendPerItemGas, msgGrantAllowancePerItemGas uint64,
) Params {
	return Params{
		MaxTxSize:                   maxTxSize,
		MinGasPerByte:               minGasPerByte,
		FixedMsgGas:                 fixedMsgGas,
		MsgGrantPerItemGas:          msgGrantPerItemGas,
		MsgMultiSendPerItemGas:      msgMultiSendPerItemGas,
		MsgGrantAllowancePerItemGas: msgGrantAllowancePerItemGas,
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
		paramtypes.NewParamSetPair(KeyFixedMsgGas, &p.FixedMsgGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgGrantPerItemGas, &p.MsgGrantPerItemGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgMultiSendPerItemGas, &p.MsgMultiSendPerItemGas, validateMsgGas),
		paramtypes.NewParamSetPair(KeyMsgGrantAllowancePerItemGas, &p.MsgGrantAllowancePerItemGas, validateMsgGas),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MaxTxSize:                   DefaultMaxTxSize,
		MinGasPerByte:               DefaultMinGasPerByte,
		FixedMsgGas:                 DefaultFixedMsgGas,
		MsgGrantPerItemGas:          DefaultMsgGrantPerItemGas,
		MsgMultiSendPerItemGas:      DefaultMsgMultiSendPerItemGas,
		MsgGrantAllowancePerItemGas: DefaultMsgGrantAllowancePerItemGas,
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
		return fmt.Errorf("invalid msg send fee: %d", v)
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
	if err := validateMsgGas(p.FixedMsgGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgGrantPerItemGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgMultiSendPerItemGas); err != nil {
		return err
	}
	if err := validateMsgGas(p.MsgGrantAllowancePerItemGas); err != nil {
		return err
	}

	return nil
}
