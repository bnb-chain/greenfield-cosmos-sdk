package types

import "math"

func NewUpgradeConfig() UpgradeConfig {
	return UpgradeConfig(map[int64][]*Plan{})
}

type UpgradeConfig map[int64][]*Plan

var (
	MainnetChainID = "inscription_9000-1"
	MainnetConfig  = UpgradeConfig(map[int64][]*Plan{})
)

func (c UpgradeConfig) SetPlan(plan *Plan) {
	c[plan.Height] = append(c[plan.Height], plan)
}

func (c UpgradeConfig) ClearPlan(height int64) {
	c[height] = nil
}

func (c UpgradeConfig) GetPlan(height int64) []*Plan {
	plans, exist := c[height]
	if exist {
		return c[height]
	}

	// get recent upgrade plan
	var (
		recentHeight = int64(math.MaxInt64)
	)
	for vHeight, vPlans := range c {
		if vHeight > height {
			if vHeight < recentHeight {
				plans = vPlans
				recentHeight = vHeight
			}
		}
	}
	return plans
}
