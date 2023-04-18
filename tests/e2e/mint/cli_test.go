//go:build e2e
// +build e2e

package mint

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/math"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	crosschaintypes "github.com/cosmos/cosmos-sdk/x/crosschain/types"
)

func TestE2ETestSuite(t *testing.T) {
	// don't mint token for crosschain module
	crosschaintypes.DefaultInitModuleBalance = math.ZeroInt()

	cfg := network.DefaultConfig(simapp.NewTestNetworkFixture)
	cfg.NumValidators = 1
	suite.Run(t, NewE2ETestSuite(cfg))
}
