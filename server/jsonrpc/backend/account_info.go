package backend

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	rpctypes "github.com/evmos/ethermint/rpc/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

	paramsReq := stakingtypes.QueryParamsRequest{}
	queryStakingClient := stakingtypes.NewQueryClient(b.clientCtx)

	paramsRes, err := queryStakingClient.Params(b.ctx, &paramsReq)
	if err != nil {
		return (*hexutil.Big)(big.NewInt(0)), err
	}

	balanceReq := banktypes.QueryBalanceRequest{
		Address: address.String(),
		Denom:   paramsRes.Params.BondDenom,
	}
	queryBankClient := banktypes.NewQueryClient(b.clientCtx)
	balanceRes, err := queryBankClient.Balance(b.ctx, &balanceReq)
	if err != nil {
		return (*hexutil.Big)(big.NewInt(0)), err
	}

	return (*hexutil.Big)(balanceRes.Balance.Amount.BigInt()), nil
}
