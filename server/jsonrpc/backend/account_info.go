package backend

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	rpctypes "github.com/evmos/ethermint/rpc/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (b *Backend) GetBalance(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Big, error) {
	blockNum, err := b.BlockNumberFromTendermint(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	_, err = b.TendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	req := &banktypes.QueryBalanceRequest{
		Address: address.String(),
		Denom:   "BNB",
	}
	queryClient := banktypes.NewQueryClient(b.clientCtx)

	res, err := queryClient.Balance(b.ctx, req)
	if err != nil {
		return (*hexutil.Big)(big.NewInt(0)), err
	}

	return (*hexutil.Big)(res.Balance.Amount.BigInt()), nil
}
