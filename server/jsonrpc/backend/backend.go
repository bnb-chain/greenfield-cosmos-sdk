package backend

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server/config"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EVMBackend implements the functionality shared within ethereum namespaces
// as defined by EIP-1474: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1474.md
// Implemented by Backend.
// This server is exclusively designed for wallet connections to sign EIP712 typed messages.
// It does not offer any other unnecessary EVM RPC APIs, such as eth_call.
type EVMBackend interface {
	BlockNumber() (hexutil.Uint64, error)
	GetBlockByNumber(blockNum rpctypes.BlockNumber, fullTx bool) (map[string]interface{}, error)
	GetBalance(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Big, error)
	ChainID() (*hexutil.Big, error)
}

var _ EVMBackend = (*Backend)(nil)

// Backend implements the BackendI interface
type Backend struct {
	ctx       context.Context
	clientCtx client.Context
	logger    log.Logger
	chainID   *big.Int
	cfg       config.Config
}

// NewBackend creates a new Backend instance for cosmos and ethereum namespaces
func NewBackend(
	viper *viper.Viper,
	logger log.Logger,
	clientCtx client.Context,
) *Backend {
	chainID, err := sdk.ParseChainID(clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	appConf, err := config.GetConfig(viper)
	if err != nil {
		panic(err)
	}

	return &Backend{
		ctx:       context.Background(),
		clientCtx: clientCtx,
		logger:    logger.With("json-rpc", "backend"),
		chainID:   chainID,
		cfg:       appConf,
	}
}
