package types

func NewUpgradeConfig() UpgradeConfig {
	return UpgradeConfig(map[int64]Plan{})
}

type UpgradeConfig map[int64]Plan

var (
	MainnetChainID = "inscription_9000-1"
	MainnetConfig  = UpgradeConfig(map[int64]Plan{})
)

func (c UpgradeConfig) SetPlan(plan Plan) {
	c[plan.Height] = plan
}
