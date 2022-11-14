package testutil

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	ethHd "github.com/evmos/ethermint/crypto/hd"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/testutil/network"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.SigningAlgo = string(ethHd.EthSecp256k1Type)
	cfg.KeyringOptions = []keyring.Option{ethHd.EthSecp256k1Option()}
	cfg.NumValidators = 1
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
