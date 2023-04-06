//go:build e2e
// +build e2e

package distribution

import (
	"testing"

	"cosmossdk.io/math"
	crosschaintypes "github.com/cosmos/cosmos-sdk/x/crosschain/types"

	"github.com/stretchr/testify/suite"
)

func TestE2ETestSuite(t *testing.T) {
	// don't mint token for crosschain module
	crosschaintypes.DefaultInitModuleBalance = math.ZeroInt()

	suite.Run(t, new(E2ETestSuite))
}

func TestGRPCQueryTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCQueryTestSuite))
}

func TestWithdrawAllSuite(t *testing.T) {
	suite.Run(t, new(WithdrawAllTestSuite))
}
