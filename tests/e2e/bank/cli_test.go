//go:build e2e
// +build e2e

package client

import (
	"testing"

	"cosmossdk.io/math"
	"cosmossdk.io/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	crosschaintypes "github.com/cosmos/cosmos-sdk/x/crosschain/types"

	"github.com/stretchr/testify/suite"
)

func TestE2ETestSuite(t *testing.T) {
	// don't mint token for crosschain module
	crosschaintypes.DefaultInitModuleBalance = math.ZeroInt()

	cfg := network.DefaultConfig(simapp.NewTestNetworkFixture)
	cfg.NumValidators = 1
	suite.Run(t, NewE2ETestSuite(cfg))
}
