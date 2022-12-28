package keeper

import (
	"errors"

	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type Option func(k *Keeper) error

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

func RegisterUpgradePlan(ctx sdk.Context,
	plans []serverconfig.UpgradeConfig, handler map[string]types.UpgradeHandler,
) Option {
	return func(k *Keeper) error {
		for _, plan := range convertUpgradeConfig(ctx, plans) {
			err := k.ScheduleUpgrade(ctx, plan)
			if err != nil &&
				!errors.Is(err, types.ErrUpgradeScheduled) && !errors.Is(err, types.ErrUpgradeCompleted) {
				return err
			}
		}

		k.upgradeHandlers = handler
		return nil
	}
}
