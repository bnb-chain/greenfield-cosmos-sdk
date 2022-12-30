package keeper

import (
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func convertUpgradeConfig(ctx sdk.Context, plans []serverconfig.UpgradeConfig) types.UpgradeConfig {
	upgradeConfig := types.NewUpgradeConfig()
	if ctx.ChainID() == types.MainnetChainID {
		upgradeConfig = types.MainnetConfig
	}

	// override by app config
	for _, plan := range plans {
		upgradeConfig.SetPlan(types.Plan{
			Name:   plan.Name,
			Height: plan.Height,
			Info:   plan.Info,
		})
	}

	return upgradeConfig
}

// Option function for Keeper
type KeeperOption func(k *Keeper)

// RegisterUpgradePlan returns a KeeperOption to set the upgrade plan into the upgrade keeper
func RegisterUpgradePlan(ctx sdk.Context,
	plans []serverconfig.UpgradeConfig,
) KeeperOption {
	return func(k *Keeper) {
		k.upgradeConfig = convertUpgradeConfig(ctx, plans)
	}
}

// RegisterUpgradeHandler returns a KeeperOption to set the upgrade handler into the upgrade keeper
func RegisterUpgradeHandler(handlers map[string]types.UpgradeHandler) KeeperOption {
	return func(k *Keeper) {
		for name, handler := range handlers {
			k.SetUpgradeHandler(name, handler)
		}
	}
}
