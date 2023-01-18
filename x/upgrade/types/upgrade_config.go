package types

import (
	"math"
)

// ex.
// const (
// 	BEP111 =	"BEP111"
// )

const (
	EnablePublicDelegationUpgrade = "EnablePublicDelegationUpgrade"
)

var (
	MainnetChainID = "inscription_9000-1"
	MainnetConfig  = NewUpgradeConfig()
	// .SetPlan(&Plan{
	// 	Name: BEP111,
	// 	Height: 100,
	// 	Info: "https://github.com/bnb-chain/BEPs/pull/111",
	// })
)

func NewUpgradeConfig() *UpgradeConfig {
	return &UpgradeConfig{
		keys:     make(map[string]*key),
		elements: make(map[int64][]*Plan),
	}
}

type key struct {
	index  int
	height int64
}

type UpgradeConfig struct {
	keys     map[string]*key
	elements map[int64][]*Plan
}

func (c *UpgradeConfig) SetPlan(plan *Plan) *UpgradeConfig {
	if key, ok := c.keys[plan.Name]; ok {
		if c.elements[key.height][key.index].Height == plan.Height {
			*c.elements[key.height][key.index] = *plan
			return c
		}

		c.elements[key.height] = append(c.elements[key.height][:key.index], c.elements[key.height][key.index+1:]...)
	}

	c.elements[plan.Height] = append(c.elements[plan.Height], plan)
	c.keys[plan.Name] = &key{height: plan.Height, index: len(c.elements[plan.Height]) - 1}

	return c
}

func (c *UpgradeConfig) Clear(height int64) {
	for _, plan := range c.elements[height] {
		delete(c.keys, plan.Name)
	}
	c.elements[height] = nil
}

func (c *UpgradeConfig) GetPlan(height int64) []*Plan {
	plans, exist := c.elements[height]
	if exist && len(plans) != 0 {
		return plans
	}

	// get recent upgrade plan
	recentHeight := int64(math.MaxInt64)
	for vHeight, vPlans := range c.elements {
		if vHeight > height {
			if vHeight < recentHeight {
				plans = vPlans
				recentHeight = vHeight
			}
		}
	}
	return plans
}
