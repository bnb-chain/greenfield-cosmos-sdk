package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/server/jsonrpc/backend"
)

type EthereumAPI interface {
	// Getting Blocks
	//
	// Retrieves information from a particular block in the blockchain.
	BlockNumber() (hexutil.Uint64, error)
	GetBlockByNumber(ethBlockNum rpctypes.BlockNumber, fullTx bool) (map[string]interface{}, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error)
	GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint
	GetBlockTransactionCountByNumber(blockNum rpctypes.BlockNumber) *hexutil.Uint

	// Reading Transactions
	//
	// Retrieves information on the state data for addresses regardless of whether
	// it is a user or a smart contract.
	GetTransactionByHash(hash common.Hash) (*rpctypes.RPCTransaction, error)
	GetTransactionCount(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Uint64, error)
	GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error)
	GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*rpctypes.RPCTransaction, error)
	GetTransactionByBlockNumberAndIndex(blockNum rpctypes.BlockNumber, idx hexutil.Uint) (*rpctypes.RPCTransaction, error)
	// eth_getBlockReceipts

	// Writing Transactions
	//
	// Allows developers to both send ETH from one address to another, write data
	// on-chain, and interact with smart contracts.
	SendRawTransaction(data hexutil.Bytes) (common.Hash, error)
	SendTransaction(args evmtypes.TransactionArgs) (common.Hash, error)
	// eth_sendPrivateTransaction
	// eth_cancel	PrivateTransaction

	// Account Information
	//
	// Returns information regarding an address's stored on-chain data.
	Accounts() ([]common.Address, error)
	GetBalance(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Big, error)
	GetStorageAt(address common.Address, key string, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error)
	GetCode(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error)
	GetProof(address common.Address, storageKeys []string, blockNrOrHash rpctypes.BlockNumberOrHash) (*rpctypes.AccountResult, error)

	// EVM/Smart Contract Execution
	//
	// Allows developers to read data from the blockchain which includes executing
	// smart contracts. However, no data is published to the Ethereum network.
	Call(args evmtypes.TransactionArgs, blockNrOrHash rpctypes.BlockNumberOrHash, _ *rpctypes.StateOverride) (hexutil.Bytes, error)

	// Chain Information
	//
	// Returns information on the Ethereum network and internal settings.
	ProtocolVersion() hexutil.Uint
	GasPrice() (*hexutil.Big, error)
	EstimateGas(args evmtypes.TransactionArgs, blockNrOptional *rpctypes.BlockNumber) (hexutil.Uint64, error)
	FeeHistory(blockCount rpc.DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*rpctypes.FeeHistoryResult, error)
	MaxPriorityFeePerGas() (*hexutil.Big, error)
	ChainId() (*hexutil.Big, error)

	// Getting Uncles
	//
	// Returns information on uncle blocks are which are network rejected blocks and replaced by a canonical block instead.
	GetUncleByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) map[string]interface{}
	GetUncleByBlockNumberAndIndex(number, idx hexutil.Uint) map[string]interface{}
	GetUncleCountByBlockHash(hash common.Hash) hexutil.Uint
	GetUncleCountByBlockNumber(blockNum rpctypes.BlockNumber) hexutil.Uint

	// Proof of Work
	Hashrate() hexutil.Uint64
	Mining() bool

	// Other
	Syncing() (interface{}, error)
	Coinbase() (string, error)
	Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error)
	GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error)
	SignTypedData(address common.Address, typedData apitypes.TypedData) (hexutil.Bytes, error)
	FillTransaction(args evmtypes.TransactionArgs) (*rpctypes.SignTransactionResult, error)
	Resend(ctx context.Context, args evmtypes.TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error)
	GetPendingTransactions() ([]*rpctypes.RPCTransaction, error)
	// eth_signTransaction (on Ethereum.org)
	// eth_getCompilers (on Ethereum.org)
	// eth_compileSolidity (on Ethereum.org)
	// eth_compileLLL (on Ethereum.org)
	// eth_compileSerpent (on Ethereum.org)
	// eth_getWork (on Ethereum.org)
	// eth_submitWork (on Ethereum.org)
	// eth_submitHashrate (on Ethereum.org)
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

// /////////////////////////////////////////////////////////////////////////////
// /                              Blocks						                ///
// /////////////////////////////////////////////////////////////////////////////

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

// GetBlockByHash returns the block identified by hash.
func (e *PublicAPI) GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	e.logger.Debug("eth_getBlockByHash", "hash", hash.Hex(), "full", fullTx)
	return e.backend.GetBlockByHash(hash, fullTx)
}

// /////////////////////////////////////////////////////////////////////////////
// /                             Read Txs					                ///
// /////////////////////////////////////////////////////////////////////////////

// GetTransactionByHash returns the transaction identified by hash.
func (e *PublicAPI) GetTransactionByHash(hash common.Hash) (*rpctypes.RPCTransaction, error) {
	e.logger.Debug("eth_getTransactionByHash", "hash", hash.Hex())
	return e.backend.GetTransactionByHash(hash)
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (e *PublicAPI) GetTransactionCount(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Uint64, error) {
	e.logger.Debug("eth_getTransactionCount", "address", address.Hex(), "block number or hash", blockNrOrHash)
	blockNum, err := e.backend.BlockNumberFromTendermint(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	return e.backend.GetTransactionCount(address, blockNum)
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (e *PublicAPI) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
	hexTx := hash.Hex()
	e.logger.Debug("eth_getTransactionReceipt", "hash", hexTx)
	return e.backend.GetTransactionReceipt(hash)
}

// GetBlockTransactionCountByHash returns the number of transactions in the block identified by hash.
func (e *PublicAPI) GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint {
	e.logger.Debug("eth_getBlockTransactionCountByHash", "hash", hash.Hex())
	return e.backend.GetBlockTransactionCountByHash(hash)
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block identified by number.
func (e *PublicAPI) GetBlockTransactionCountByNumber(blockNum rpctypes.BlockNumber) *hexutil.Uint {
	e.logger.Debug("eth_getBlockTransactionCountByNumber", "height", blockNum.Int64())
	return e.backend.GetBlockTransactionCountByNumber(blockNum)
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
func (e *PublicAPI) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*rpctypes.RPCTransaction, error) {
	e.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash.Hex(), "index", idx)
	return e.backend.GetTransactionByBlockHashAndIndex(hash, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (e *PublicAPI) GetTransactionByBlockNumberAndIndex(blockNum rpctypes.BlockNumber, idx hexutil.Uint) (*rpctypes.RPCTransaction, error) {
	e.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "number", blockNum, "index", idx)
	return e.backend.GetTransactionByBlockNumberAndIndex(blockNum, idx)
}

// /////////////////////////////////////////////////////////////////////////////
// /                             Write Txs					                ///
// /////////////////////////////////////////////////////////////////////////////

// SendRawTransaction send a raw Ethereum transaction.
func (e *PublicAPI) SendRawTransaction(data hexutil.Bytes) (common.Hash, error) {
	e.logger.Debug("eth_sendRawTransaction", "length", len(data))
	return e.backend.SendRawTransaction(data)
}

// SendTransaction sends an Ethereum transaction.
func (e *PublicAPI) SendTransaction(args evmtypes.TransactionArgs) (common.Hash, error) {
	e.logger.Debug("eth_sendTransaction", "args", args.String())
	return e.backend.SendTransaction(args)
}

// /////////////////////////////////////////////////////////////////////////////
// /                           Account Information				            ///
// /////////////////////////////////////////////////////////////////////////////

// Accounts returns the list of accounts available to this node.
func (e *PublicAPI) Accounts() ([]common.Address, error) {
	e.logger.Debug("eth_accounts")
	return e.backend.Accounts()
}

// GetBalance returns the provided account's balance up to the provided block number.
func (e *PublicAPI) GetBalance(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Big, error) {
	e.logger.Debug("eth_getBalance", "address", address.String(), "block number or hash", blockNrOrHash)
	return e.backend.GetBalance(address, blockNrOrHash)
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (e *PublicAPI) GetStorageAt(address common.Address, key string, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	e.logger.Debug("eth_getStorageAt", "address", address.Hex(), "key", key, "block number or hash", blockNrOrHash)
	return e.backend.GetStorageAt(address, key, blockNrOrHash)
}

// GetCode returns the contract code at the given address and block number.
func (e *PublicAPI) GetCode(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	e.logger.Debug("eth_getCode", "address", address.Hex(), "block number or hash", blockNrOrHash)
	return e.backend.GetCode(address, blockNrOrHash)
}

// GetProof returns an account object with proof and any storage proofs
func (e *PublicAPI) GetProof(address common.Address,
	storageKeys []string,
	blockNrOrHash rpctypes.BlockNumberOrHash,
) (*rpctypes.AccountResult, error) {
	e.logger.Debug("eth_getProof", "address", address.Hex(), "keys", storageKeys, "block number or hash", blockNrOrHash)
	return e.backend.GetProof(address, storageKeys, blockNrOrHash)
}

// Call performs a raw contract call.
func (e *PublicAPI) Call(args evmtypes.TransactionArgs,
	blockNrOrHash rpctypes.BlockNumberOrHash,
	_ *rpctypes.StateOverride,
) (hexutil.Bytes, error) {
	e.logger.Debug("eth_call", "args", args.String(), "block number or hash", blockNrOrHash)
	return nil, fmt.Errorf("should not be called")
}

// /////////////////////////////////////////////////////////////////////////////
// /                            Chain Information					        ///
// /////////////////////////////////////////////////////////////////////////////

// ProtocolVersion returns the supported Ethereum protocol version.
func (e *PublicAPI) ProtocolVersion() hexutil.Uint {
	e.logger.Debug("eth_protocolVersion")
	return hexutil.Uint(ethermint.ProtocolVersion)
}

// GasPrice returns the current gas price based on Ethermint's gas price oracle.
func (e *PublicAPI) GasPrice() (*hexutil.Big, error) {
	e.logger.Debug("eth_gasPrice")
	return e.backend.GasPrice()
}

// EstimateGas returns an estimate of gas usage for the given smart contract call.
func (e *PublicAPI) EstimateGas(args evmtypes.TransactionArgs, blockNrOptional *rpctypes.BlockNumber) (hexutil.Uint64, error) {
	e.logger.Debug("eth_estimateGas")
	return e.backend.EstimateGas(args, blockNrOptional)
}

func (e *PublicAPI) FeeHistory(blockCount rpc.DecimalOrHex,
	lastBlock rpc.BlockNumber,
	rewardPercentiles []float64,
) (*rpctypes.FeeHistoryResult, error) {
	e.logger.Debug("eth_feeHistory")
	return e.backend.FeeHistory(blockCount, lastBlock, rewardPercentiles)
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.
func (e *PublicAPI) MaxPriorityFeePerGas() (*hexutil.Big, error) {
	e.logger.Debug("eth_maxPriorityFeePerGas")
	return nil, fmt.Errorf("should not be called")
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (e *PublicAPI) ChainId() (*hexutil.Big, error) { // nolint
	e.logger.Debug("eth_chainId")
	return e.backend.ChainID()
}

// /////////////////////////////////////////////////////////////////////////////
// /                              Uncles										///
// /////////////////////////////////////////////////////////////////////////////

// GetUncleByBlockHashAndIndex returns the uncle identified by hash and index. Always returns nil.
func (e *PublicAPI) GetUncleByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// GetUncleByBlockNumberAndIndex returns the uncle identified by number and index. Always returns nil.
func (e *PublicAPI) GetUncleByBlockNumberAndIndex(number, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// GetUncleCountByBlockHash returns the number of uncles in the block identified by hash. Always zero.
func (e *PublicAPI) GetUncleCountByBlockHash(hash common.Hash) hexutil.Uint {
	return 0
}

// GetUncleCountByBlockNumber returns the number of uncles in the block identified by number. Always zero.
func (e *PublicAPI) GetUncleCountByBlockNumber(blockNum rpctypes.BlockNumber) hexutil.Uint {
	return 0
}

// /////////////////////////////////////////////////////////////////////////////
// /                             Proof of Work							    ///
// /////////////////////////////////////////////////////////////////////////////

// Hashrate returns the current node's hashrate. Always 0.
func (e *PublicAPI) Hashrate() hexutil.Uint64 {
	e.logger.Debug("eth_hashrate")
	return 0
}

// Mining returns whether this node is currently mining. Always false.
func (e *PublicAPI) Mining() bool {
	e.logger.Debug("eth_mining")
	return false
}

// /////////////////////////////////////////////////////////////////////////////
// /                               Other 						            ///
// /////////////////////////////////////////////////////////////////////////////

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronize from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (e *PublicAPI) Syncing() (interface{}, error) {
	e.logger.Debug("eth_syncing")
	return e.backend.Syncing()
}

// Coinbase is the address that staking rewards will be sent to (alias for Etherbase).
func (e *PublicAPI) Coinbase() (string, error) {
	e.logger.Debug("eth_coinbase")
	return "", fmt.Errorf("should not be called")
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (e *PublicAPI) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	e.logger.Debug("eth_sign", "address", address.Hex(), "data", common.Bytes2Hex(data))
	return e.backend.Sign(address, data)
}

// GetTransactionLogs returns the logs given a transaction hash.
func (e *PublicAPI) GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error) {
	e.logger.Debug("eth_getTransactionLogs", "hash", txHash)
	return nil, fmt.Errorf("should not be called")
}

// SignTypedData signs EIP-712 conformant typed data
func (e *PublicAPI) SignTypedData(address common.Address, typedData apitypes.TypedData) (hexutil.Bytes, error) {
	e.logger.Debug("eth_signTypedData", "address", address.Hex(), "data", typedData)
	return e.backend.SignTypedData(address, typedData)
}

// FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
// on a given unsigned transaction, and returns it to the caller for further
// processing (signing + broadcast).
func (e *PublicAPI) FillTransaction(args evmtypes.TransactionArgs) (*rpctypes.SignTransactionResult, error) {
	return nil, fmt.Errorf("should not be called")
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (e *PublicAPI) Resend(_ context.Context,
	args evmtypes.TransactionArgs,
	gasPrice *hexutil.Big,
	gasLimit *hexutil.Uint64,
) (common.Hash, error) {
	e.logger.Debug("eth_resend", "args", args.String())
	return e.backend.Resend(args, gasPrice, gasLimit)
}

// GetPendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (e *PublicAPI) GetPendingTransactions() ([]*rpctypes.RPCTransaction, error) {
	return nil, fmt.Errorf("should not be called")
}
