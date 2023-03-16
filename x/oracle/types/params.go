package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultRelayerTimeout     uint64 = 40  // in s
	DefaultRelayerRewardShare uint32 = 50  // in s
	DefaultRealyerInterval    uint64 = 600 // in s
)

var (
	KeyParamRelayerTimeout     = []byte("RelayerTimeout")
	KeyParamRelayerBackoffTime = []byte("RelayerBackoffTime")
	KeyParamRelayerRewardShare = []byte("RelayerRewardShare")
	KeyParamRelayerInterval    = []byte("RelayerInterval")
)

func DefaultParams() Params {
	return Params{
		RelayerTimeout:     DefaultRelayerTimeout,
		RelayerRewardShare: DefaultRelayerRewardShare,
		RelayerInterval:    DefaultRealyerInterval,
	}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyParamRelayerTimeout, &p.RelayerTimeout, validateRelayerTimeout),
		paramtypes.NewParamSetPair(KeyParamRelayerRewardShare, &p.RelayerRewardShare, validateRelayerRewardShare),
		paramtypes.NewParamSetPair(KeyParamRelayerInterval, &p.RelayerInterval, validateRelayerInterval),
	}
}

func validateRelayerTimeout(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("the relayer timeout must be positive: %d", v)
	}

	return nil
}

func validateRelayerRewardShare(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("the relayer reward share should be positive: %d", v)
	}

	if v > 100 {
		return fmt.Errorf("the relayer reward share should not be larger than 100")
	}

	return nil
}

func validateRelayerInterval(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("the relayer relay interval should be positive: %d", v)
	}

	return nil
}
