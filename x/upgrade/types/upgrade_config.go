package types

import (
	"math"

	"github.com/cosmos/cosmos-sdk/types"
)

const (
	// EnablePublicDelegationUpgrade is the upgrade name for enabling public delegation
	EnablePublicDelegationUpgrade = types.EnablePublicDelegationUpgrade

	// Nagqu is the upgrade name for Nagqu upgrade
	Nagqu = types.Nagqu

	// Pampas is the upgrade name for Pampas upgrade
	Pampas = types.Pampas

	// Manchurian is the upgrade name for Manchurian upgrade
	Manchurian = types.Manchurian

	// Hulunbeier is the upgrade name for Hulunbeier upgrade
	Hulunbeier = types.Hulunbeier

	HulunbeierPatch = types.HulunbeierPatch

	// Ural is the upgrade name for Ural upgrade
	Ural = types.Ural

	// Pawnee is the upgrade name for Pawnee upgrade
	Pawnee = types.Pawnee

	// Serengeti is the upgrade name for Serengeti upgrade
	Serengeti = types.Serengeti

	// Erdos is the upgrade name for Erdos upgrade
	Erdos = types.Erdos

	// Veld is the upgrade name for Veld upgrade
	Veld = types.Veld
)

// The default upgrade config for networks
var (
	MainnetChainID = "greenfield_1017-1"
	MainnetConfig  = NewUpgradeConfig().SetPlan(&Plan{
		Name:   Nagqu,
		Height: 1,
		Info:   "Nagqu hardfork",
	}).SetPlan(&Plan{
		Name:   Pampas,
		Height: 2006197,
		Info:   "Pampas hardfork",
	}).SetPlan(&Plan{
		Name:   Manchurian,
		Height: 3426973,
		Info:   "Manchurian hardfork",
	}).SetPlan(&Plan{
		Name:   Hulunbeier,
		Height: 4653883,
		Info:   "Hulunbeier hardfork",
	}).SetPlan(&Plan{
		Name:   HulunbeierPatch,
		Height: 4653883,
		Info:   "Hulunbeier hardfork",
	}).SetPlan(&Plan{
		Name:   Ural,
		Height: 5347231,
		Info:   "Ural hardfork",
	}).SetPlan(&Plan{
		Name:   Pawnee,
		Height: 6239520,
		Info:   "Pawnee hardfork",
	}).SetPlan(&Plan{
		Name:   Serengeti,
		Height: 6863285,
		Info:   "Serengeti hardfork",
	}).SetPlan(&Plan{
		Name:   Erdos,
		Height: 7861456,
		Info:   "Erdos hardfork",
	}).SetPlan(&Plan{
		Name:   Veld,
		Height: 9030588,
		Info:   "Veld hardfork",
	})

	TestnetChainID = "greenfield_5600-1"
	TestnetConfig  = NewUpgradeConfig().SetPlan(&Plan{
		Name:   Nagqu,
		Height: 471350,
		Info:   "Nagqu hardfork",
	}).SetPlan(&Plan{
		Name:   Pampas,
		Height: 2427233,
		Info:   "Pampas hardfork",
	}).SetPlan(&Plan{
		Name:   Manchurian,
		Height: 3922485,
		Info:   "Manchurian hardfork",
	}).SetPlan(&Plan{
		Name:   Hulunbeier,
		Height: 4849568,
		Info:   "Hulunbeier hardfork",
	}).SetPlan(&Plan{
		Name:   HulunbeierPatch,
		Height: 4849568,
		Info:   "Hulunbeier hardfork",
	}).SetPlan(&Plan{
		Name:   Ural,
		Height: 5761391,
		Info:   "Ural hardfork",
	}).SetPlan(&Plan{
		Name:   Pawnee,
		Height: 6623127,
		Info:   "Pawnee hardfork",
	}).SetPlan(&Plan{
		Name:   Serengeti,
		Height: 7354695,
		Info:   "Serengeti hardfork",
	}).SetPlan(&Plan{
		Name:   Erdos,
		Height: 8116724,
		Info:   "Erdos hardfork",
	}).SetPlan(&Plan{
		Name:   Veld,
		Height: 9379516,
		Info:   "Veld hardfork",
	})
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

// UpgradeConfig is a list of upgrade plans
type UpgradeConfig struct {
	keys     map[string]*key
	elements map[int64][]*Plan
}

// SetPlan sets a new upgrade plan
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

// Clear removes all upgrade plans at a given height
func (c *UpgradeConfig) Clear(height int64) {
	for _, plan := range c.elements[height] {
		delete(c.keys, plan.Name)
	}
	c.elements[height] = nil
}

// GetPlan returns the upgrade plan at a given height
func (c *UpgradeConfig) GetPlan(height int64) []*Plan {
	plans, exist := c.elements[height]
	if exist && len(plans) != 0 {
		return plans
	}

	// get recent upgrade plan
	recentHeight := int64(math.MaxInt64)
	for vHeight, vPlans := range c.elements {
		if vHeight > height && vHeight < recentHeight {
			plans = vPlans
			recentHeight = vHeight
		}
	}
	return plans
}
