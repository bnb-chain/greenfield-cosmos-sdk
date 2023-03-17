package backend

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	rpctypes "github.com/cosmos/cosmos-sdk/server/jsonrpc/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// BlockNumber returns the current block number in abci app state. Because abci
// app state could lag behind from tendermint the latest block, it's more stable for
// the client to use the latest block number in abci app state than tendermint
// rpc.
func (b *Backend) BlockNumber() (hexutil.Uint64, error) {
	// do any grpc query, ignore the response and use the returned block height
	var header metadata.MD
	queryClient := banktypes.NewQueryClient(b.clientCtx)

	_, err := queryClient.Balance(b.ctx, &banktypes.QueryBalanceRequest{Address: "0x0000000000000000000000000000000000000000", Denom: "BNB"}, grpc.Header(&header))
	if err != nil {
		return hexutil.Uint64(0), err
	}

	blockHeightHeader := header.Get(grpctypes.GRPCBlockHeightHeader)
	if headerLen := len(blockHeightHeader); headerLen != 1 {
		return 0, fmt.Errorf("unexpected '%s' gRPC header length; got %d, expected: %d", grpctypes.GRPCBlockHeightHeader, headerLen, 1)
	}

	height, err := strconv.ParseUint(blockHeightHeader[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block height: %w", err)
	}

	return hexutil.Uint64(height), nil
}

// GetBlockByNumber returns the JSON-RPC compatible Ethereum block identified by
// block number. Depending on fullTx it either returns the full transaction
// objects or if false only the hashes of the transactions.
func (b *Backend) GetBlockByNumber(blockNum rpctypes.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	resBlock, err := b.TendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, nil
	}

	// return if requested block height is greater than the current one
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	res, err := b.RPCBlockFromTendermintBlock(resBlock, nil, false)
	if err != nil {
		b.logger.Debug("GetEthBlockFromTendermint failed", "height", blockNum, "error", err.Error())
		return nil, err
	}

	return res, nil
}

// TendermintBlockByNumber returns a Tendermint-formatted block for a given
// block number
func (b *Backend) TendermintBlockByNumber(blockNum rpctypes.BlockNumber) (*tmrpctypes.ResultBlock, error) {
	height := blockNum.Int64()
	if height <= 0 {
		// fetch the latest block number from the app state, more accurate than the tendermint block store state.
		n, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}
		height = int64(n)
	}
	resBlock, err := b.clientCtx.Client.Block(b.ctx, &height)
	if err != nil {
		b.logger.Debug("tendermint client failed to get block", "height", height, "error", err.Error())
		return nil, err
	}

	if resBlock.Block == nil {
		b.logger.Debug("TendermintBlockByNumber block not found", "height", height)
		return nil, nil
	}

	return resBlock, nil
}

// BlockNumberFromTendermint returns the BlockNumber from BlockNumberOrHash
func (b *Backend) BlockNumberFromTendermint(blockNrOrHash rpctypes.BlockNumberOrHash) (rpctypes.BlockNumber, error) {
	switch {
	case blockNrOrHash.BlockHash == nil && blockNrOrHash.BlockNumber == nil:
		return rpctypes.EthEarliestBlockNumber, fmt.Errorf("types BlockHash and BlockNumber cannot be both nil")
	case blockNrOrHash.BlockNumber != nil:
		return *blockNrOrHash.BlockNumber, nil
	default:
		return rpctypes.EthEarliestBlockNumber, nil
	}
}

// HeaderByNumber returns the block header identified by height.
func (b *Backend) HeaderByNumber(blockNum rpctypes.BlockNumber) (*ethtypes.Header, error) {
	resBlock, err := b.TendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		return nil, errors.Errorf("block not found for height %d", blockNum)
	}

	ethHeader := EthHeaderFromTendermint(resBlock.Block.Header)
	return ethHeader, nil
}

// RPCBlockFromTendermintBlock returns a JSON-RPC compatible Ethereum block from a
// given Tendermint block and its block result.
func (b *Backend) RPCBlockFromTendermintBlock(
	resBlock *tmrpctypes.ResultBlock,
	blockRes *tmrpctypes.ResultBlockResults,
	fullTx bool,
) (map[string]interface{}, error) {
	block := resBlock.Block

	formattedBlock := rpctypes.FormatBlock(
		block.Header, block.Size(),
		0, new(big.Int).SetUint64(0),
		[]interface{}{}, ethtypes.Bloom{}, common.Address{}, nil,
	)
	return formattedBlock, nil
}

func EthHeaderFromTendermint(header tmtypes.Header) *ethtypes.Header {
	txHash := ethtypes.EmptyRootHash
	if len(header.DataHash) == 0 {
		txHash = common.BytesToHash(header.DataHash)
	}

	return &ethtypes.Header{
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		UncleHash:   ethtypes.EmptyUncleHash,
		Coinbase:    common.BytesToAddress(header.ProposerAddress),
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      txHash,
		ReceiptHash: ethtypes.EmptyRootHash,
		Difficulty:  big.NewInt(0),
		Number:      big.NewInt(header.Height),
		GasLimit:    0,
		GasUsed:     0,
		Time:        uint64(header.Time.UTC().Unix()),
		Extra:       []byte{},
		MixDigest:   common.Hash{},
		Nonce:       ethtypes.BlockNonce{},
	}
}
