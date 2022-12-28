package types

import (
	"fmt"

	"sigs.k8s.io/yaml"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultMaxTxSize       uint64 = 1024
	DefaultMinGasPerByte   uint64 = 5
	DefaultMsgSendGas      uint64 = 1e6
	DefaultMsgMultiSendGas uint64 = 8e5
)

// Parameter keys
var (
	KeyMaxTxSize       = []byte("MaxTxSize")
	KeyMinGasPerByte   = []byte("MinGasPerByte")
	KeyMsgSendGas      = []byte("MsgSendGas")
	KeyMsgMultiSendGas = []byte("MsgMultiSendGas")
)

var _ paramtypes.ParamSet = &Params{}

// NewParams creates a new Params object
func NewParams(
	maxTxSize, minGasPerByte, msgSendGas, msgMultiSendGas uint64,
) Params {
	return Params{
		MaxTxSize:       maxTxSize,
		MinGasPerByte:   minGasPerByte,
		MsgSendGas:      msgSendGas,
		MsgMultiSendGas: msgMultiSendGas,
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
		paramtypes.NewParamSetPair(KeyMsgSendGas, &p.MsgSendGas, validateMsgSendGas),
		paramtypes.NewParamSetPair(KeyMsgMultiSendGas, &p.MsgMultiSendGas, validateMsgMultiSendGas),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MaxTxSize:       DefaultMaxTxSize,
		MinGasPerByte:   DefaultMinGasPerByte,
		MsgSendGas:      DefaultMsgSendGas,
		MsgMultiSendGas: DefaultMsgMultiSendGas,
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

func validateMsgSendGas(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid msg send fee: %d", v)
	}

	return nil
}

func validateMsgMultiSendGas(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid msg multi send fee: %d", v)
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
	if err := validateMsgSendGas(p.MsgSendGas); err != nil {
		return err
	}
	if err := validateMsgMultiSendGas(p.MsgMultiSendGas); err != nil {
		return err
	}

	return nil
}
