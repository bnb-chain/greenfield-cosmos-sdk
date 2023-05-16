package types

import (
	"fmt"
)

// Default parameter values
const (
	DefaultMaxTxSize     uint64 = 64 * 1024 // 64kb
	DefaultMinGasPerByte uint64 = 5
)

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
	maxTxSize, minGasPerByte uint64,
) Params {
	return Params{
		MaxTxSize:     maxTxSize,
		MinGasPerByte: minGasPerByte,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MaxTxSize:     DefaultMaxTxSize,
		MinGasPerByte: DefaultMinGasPerByte,
	}
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

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateMaxTxSize(p.MaxTxSize); err != nil {
		return err
	}
	if err := validateMinGasPerByte(p.MinGasPerByte); err != nil {
		return err
	}

	return nil
}

// Validate gets any errors with this MsgGasParams entry.
func (mgp MsgGasParams) Validate() error {
	if mgp.MsgTypeUrl == "" {
		return fmt.Errorf("invalid msg type url. cannot be empty")
	}

	switch p := mgp.GasParams.(type) {
	case *MsgGasParams_FixedType:
		return nil
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

	return nil
}
