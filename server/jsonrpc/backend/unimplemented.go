package backend

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (b *Backend) TendermintBlockResultByNumber(height *int64) (*tmrpctypes.ResultBlockResults, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) EthMsgsFromTendermintBlock(block *tmrpctypes.ResultBlock, blockRes *tmrpctypes.ResultBlockResults) []*evmtypes.MsgEthereumTx {
	panic("should not be called")
}

func (b *Backend) ChainConfig() *params.ChainConfig {
	panic("should not be called")
}

func (b *Backend) GetCoinbase() (sdk.AccAddress, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) FeeHistory(blockCount rpc.DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*rpctypes.FeeHistoryResult, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) SuggestGasTipCap(baseFee *big.Int) (*big.Int, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) RPCFilterCap() int32 {
	panic("should not be called")
}

func (b *Backend) RPCLogsCap() int32 {
	panic("should not be called")
}

func (b *Backend) RPCBlockRangeCap() int32 {
	panic("should not be called")
}

func (b *Backend) Accounts() ([]common.Address, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) Syncing() (interface{}, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) SetEtherbase(etherbase common.Address) bool {
	panic("should not be called")
}

func (b *Backend) SetGasPrice(gasPrice hexutil.Big) bool {
	panic("should not be called")
}

func (b *Backend) ImportRawKey(privkey, password string) (common.Address, error) {
	return common.Address{}, fmt.Errorf("should not be called")
}

func (b *Backend) ListAccounts() ([]common.Address, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) NewMnemonic(uid string, language keyring.Language, hdPath, bip39Passphrase string, algo keyring.SignatureAlgo) (*keyring.Record, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) UnprotectedAllowed() bool {
	panic("should not be called")
}

func (b *Backend) RPCGasCap() uint64 {
	panic("should not be called")
}

func (b *Backend) RPCEVMTimeout() time.Duration {
	panic("should not be called")
}

func (b *Backend) RPCTxFeeCap() float64 {
	panic("should not be called")
}

func (b *Backend) RPCMinGasPrice() int64 {
	panic("should not be called")
}

func (b *Backend) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) SendTransaction(args evmtypes.TransactionArgs) (common.Hash, error) {
	return common.Hash{}, fmt.Errorf("should not be called")
}

func (b *Backend) SignTypedData(address common.Address, typedData apitypes.TypedData) (hexutil.Bytes, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetCode(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetStorageAt(address common.Address, key string, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetProof(address common.Address, storageKeys []string, blockNrOrHash rpctypes.BlockNumberOrHash) (*rpctypes.AccountResult, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetTransactionCount(address common.Address, blockNum rpctypes.BlockNumber) (*hexutil.Uint64, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetTransactionByHash(txHash common.Hash) (*rpctypes.RPCTransaction, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetTxByEthHash(txHash common.Hash) (*ethermint.TxResult, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetTxByTxIndex(height int64, txIndex uint) (*ethermint.TxResult, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetTransactionByBlockAndIndex(block *tmrpctypes.ResultBlock, idx hexutil.Uint) (*rpctypes.RPCTransaction, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*rpctypes.RPCTransaction, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetTransactionByBlockNumberAndIndex(blockNum rpctypes.BlockNumber, idx hexutil.Uint) (*rpctypes.RPCTransaction, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetLogs(hash common.Hash) ([][]*ethtypes.Log, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetLogsByHeight(height *int64) ([][]*ethtypes.Log, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) BloomStatus() (uint64, uint64) {
	panic("should not be called")
}

func (b *Backend) TraceTransaction(hash common.Hash, config *evmtypes.TraceConfig) (interface{}, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) TraceBlock(height rpctypes.BlockNumber, config *evmtypes.TraceConfig, block *tmrpctypes.ResultBlock) ([]*evmtypes.TxTraceResult, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint {
	panic("should not be called")
}

func (b *Backend) GetBlockTransactionCountByNumber(blockNum rpctypes.BlockNumber) *hexutil.Uint {
	panic("should not be called")
}

func (b *Backend) GetBlockTransactionCount(block *tmrpctypes.ResultBlock) *hexutil.Uint {
	panic("should not be called")
}

func (b *Backend) TendermintBlockByHash(blockHash common.Hash) (*tmrpctypes.ResultBlock, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) BlockNumberFromTendermintByHash(blockHash common.Hash) (*big.Int, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) BlockBloom(blockRes *tmrpctypes.ResultBlockResults) (ethtypes.Bloom, error) {
	return ethtypes.Bloom{}, fmt.Errorf("should not be called")
}

func (b *Backend) PendingTransactions() ([]*sdk.Tx, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) GlobalMinGasPrice() (sdk.Dec, error) {
	return sdk.Dec{}, fmt.Errorf("should not be called")
}

func (b *Backend) BaseFee(blockRes *tmrpctypes.ResultBlockResults) (*big.Int, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) EstimateGas(args evmtypes.TransactionArgs, blockNrOptional *rpctypes.BlockNumber) (hexutil.Uint64, error) {
	return hexutil.Uint64(0), fmt.Errorf("should not be called")
}

func (b *Backend) DoCall(
	args evmtypes.TransactionArgs, blockNr rpctypes.BlockNumber,
) (*evmtypes.MsgEthereumTxResponse, error) {
	return nil, fmt.Errorf("should not be called")
}

func (b *Backend) Resend(args evmtypes.TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	return common.Hash{}, fmt.Errorf("should not be called")
}

func (b *Backend) SendRawTransaction(data hexutil.Bytes) (common.Hash, error) {
	return common.Hash{}, fmt.Errorf("should not be called")
}

func (b *Backend) SetTxDefaults(args evmtypes.TransactionArgs) (evmtypes.TransactionArgs, error) {
	return evmtypes.TransactionArgs{}, fmt.Errorf("should not be called")
}

func (b *Backend) GasPrice() (*hexutil.Big, error) {
	return (*hexutil.Big)(big.NewInt(0)), fmt.Errorf("should not be called")
}
