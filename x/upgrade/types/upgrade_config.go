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

	// Mongolian is the upgrade name for Mongolian upgrade
	Mongolian = types.Mongolian

	// Altai is the upgrade name for Altai upgrade
	Altai = types.Altai

	// Savanna is the upgrade name for Savanna upgrade
	Savanna = types.Savanna

	// Tundra is the upgrade name for Tundra upgrade
	Tundra = types.Tundra

	// Prairie is the upgrade name for Prairie upgrade
	Prairie = types.Prairie

	// Taiga is the upgrade name for Taiga upgrade
	Taiga = types.Taiga

	// Steppe is the upgrade name for Steppe upgrade
	Steppe = types.Steppe

	// Cerrado is the upgrade name for Cerrado upgrade
	Cerrado = types.Cerrado
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
		Height: 9269910,
		Info:   "Veld hardfork",
	}).SetPlan(&Plan{
		Name:   Mongolian,
		Height: 10314605,
		Info:   "Mongolian hardfork",
	}).SetPlan(&Plan{
		Name:   Altai,
		Height: 11917971,
		Info:   "Altai hardfork",
	}).SetPlan(&Plan{
		Name:   Savanna,
		Height: 14667574,
		Info:   "Savanna hardfork",
	}).SetPlan(&Plan{
		Name:   Tundra,
		Height: 25213033,
		Info:   "Tundra hardfork",
	}).SetPlan(&Plan{
		Name:   Prairie,
		Height: 31432100,
		Info:   "Prairie hardfork",
	}).SetPlan(&Plan{
		Name:   Taiga,
		Height: 31877400,
		Info:   "Taiga hardfork",
	}).SetPlan(&Plan{
		Name:   Steppe,
		Height: 32750500,
		Info:   "Steppe hardfork",
	}).SetPlan(&Plan{
		Name:   Cerrado,
		Height: 32752300,
		Info:   "Cerrado hardfork",
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
		Height: 9581218,
		Info:   "Veld hardfork",
	}).SetPlan(&Plan{
		Name:   Mongolian,
		Height: 10780238,
		Info:   "Mongolian hardfork",
	}).SetPlan(&Plan{
		Name:   Altai,
		Height: 12513708,
		Info:   "Altai hardfork",
	}).SetPlan(&Plan{
		Name:   Savanna,
		Height: 14691233,
		Info:   "Savanna hardfork",
	}).SetPlan(&Plan{
		Name:   Tundra,
		Height: 24007800,
		Info:   "Tundra hardfork",
	}).SetPlan(&Plan{
		Name:   Prairie,
		Height: 29798400,
		Info:   "Prairie hardfork",
	}).SetPlan(&Plan{
		Name:   Taiga,
		Height: 30310800,
		Info:   "Taiga hardfork",
	}).SetPlan(&Plan{
		Name:   Steppe,
		Height: 31275600,
		Info:   "Steppe hardfork",
	}).SetPlan(&Plan{
		Name:   Cerrado,
		Height: 31279400,
		Info:   "Cerrado hardfork",
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
