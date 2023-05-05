package types

import (
	"fmt"
)

const (
	DefaultRelayerTimeout     uint64 = 40  // in s
	DefaultRelayerRewardShare uint32 = 50  // in s
	DefaultRealyerInterval    uint64 = 600 // in s
)

func DefaultParams() Params {
	return Params{
		RelayerTimeout:     DefaultRelayerTimeout,
		RelayerRewardShare: DefaultRelayerRewardShare,
		RelayerInterval:    DefaultRealyerInterval,
	}
}

func (p *Params) Validate() error {
	if err := validateRelayerTimeout(p.RelayerTimeout); err != nil {
		return err
	}

	if err := validateRelayerRewardShare(p.RelayerRewardShare); err != nil {
		return err
	}

	if err := validateRelayerInterval(p.RelayerInterval); err != nil {
		return err
	}

	return nil
}

func validateRelayerTimeout(timeout uint64) error {
	if timeout <= 0 {
		return fmt.Errorf("the relayer timeout must be positive: %d", timeout)
	}

	return nil
}

func validateRelayerRewardShare(share uint32) error {
	if share <= 0 {
		return fmt.Errorf("the relayer reward share should be positive: %d", share)
	}

	if share > 100 {
		return fmt.Errorf("the relayer reward share should not be larger than 100")
	}

	return nil
}

func validateRelayerInterval(interval uint64) error {
	if interval <= 0 {
		return fmt.Errorf("the relayer relay interval should be positive: %d", interval)
	}

	return nil
}
