package keeper

import (
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// convertUpgradeConfig converts serverconfig.UpgradeConfig to types.UpgradeConfig
func convertUpgradeConfig(chainID string, plans []serverconfig.UpgradeConfig) (*types.UpgradeConfig, error) {
	upgradeConfig := types.NewUpgradeConfig()
	if chainID == types.MainnetChainID {
		upgradeConfig = types.MainnetConfig
	}

	// override by app config
	for _, plan := range plans {
		nPlan := &types.Plan{
			Name:   plan.Name,
			Height: plan.Height,
			Info:   plan.Info,
		}
		if err := nPlan.ValidateBasic(); err != nil {
			return nil, err
		}
		upgradeConfig.SetPlan(nPlan)
	}

	return upgradeConfig, nil
}

// Option function for Keeper
type KeeperOption func(k *Keeper) error

// RegisterUpgradePlan returns a KeeperOption to set the upgrade plan into the upgrade keeper
func RegisterUpgradePlan(chianID string,
	plans []serverconfig.UpgradeConfig,
) KeeperOption {
	return func(k *Keeper) error {
		c, err := convertUpgradeConfig(chianID, plans)
		if err != nil {
			return err
		}
		k.upgradeConfig = c
		return nil
	}
}

// RegisterUpgradeHandler returns a KeeperOption to set the upgrade handler into the upgrade keeper
func RegisterUpgradeHandler(handlers map[string]types.UpgradeHandler) KeeperOption {
	return func(k *Keeper) error {
		for name, handler := range handlers {
			k.SetUpgradeHandler(name, handler)
		}
		return nil
	}
}

// RegisterUpgradeInitializer returns a KeeperOption to set the upgrade initializer into the upgrade keeper
func RegisterUpgradeInitializer(handlers map[string]types.UpgradeInitializer) KeeperOption {
	return func(k *Keeper) error {
		for name, handler := range handlers {
			k.SetUpgradeInitializer(name, handler)
		}
		return nil
	}
}
