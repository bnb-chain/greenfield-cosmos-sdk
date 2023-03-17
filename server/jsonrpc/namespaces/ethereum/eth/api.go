package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/server/jsonrpc/backend"
)

type EthereumAPI interface {
	BlockNumber() (hexutil.Uint64, error)
	GetBlockByNumber(ethBlockNum rpctypes.BlockNumber, fullTx bool) (map[string]interface{}, error)
	GetBalance(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Big, error)
	ChainId() (*hexutil.Big, error)
}

var _ EthereumAPI = (*PublicAPI)(nil)

// PublicAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicAPI struct {
	ctx     context.Context
	logger  log.Logger
	backend backend.EVMBackend
}

// NewPublicAPI creates an instance of the public ETH Web3 API.
func NewPublicAPI(logger log.Logger, backend backend.EVMBackend) *PublicAPI {
	api := &PublicAPI{
		ctx:     context.Background(),
		logger:  logger.With("json-rpc", "public-api"),
		backend: backend,
	}

	return api
}

// BlockNumber returns the current block number.
func (e *PublicAPI) BlockNumber() (hexutil.Uint64, error) {
	e.logger.Debug("eth_blockNumber")
	return e.backend.BlockNumber()
}

// GetBlockByNumber returns the block identified by number.
func (e *PublicAPI) GetBlockByNumber(ethBlockNum rpctypes.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	e.logger.Debug("eth_getBlockByNumber", "number", ethBlockNum, "full", fullTx)
	return e.backend.GetBlockByNumber(ethBlockNum, fullTx)
}

// GetBalance returns the provided account's balance up to the provided block number.
func (e *PublicAPI) GetBalance(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Big, error) {
	e.logger.Debug("eth_getBalance", "address", address.String(), "block number or hash", blockNrOrHash)
	return e.backend.GetBalance(address, blockNrOrHash)
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (e *PublicAPI) ChainId() (*hexutil.Big, error) {
	e.logger.Debug("eth_chainId")
	return e.backend.ChainID()
}
