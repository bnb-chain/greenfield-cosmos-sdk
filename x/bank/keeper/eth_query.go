package keeper

import (
	"encoding/json"
	"math/big"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtrpctypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func EthQueryBalanceHandlerGen(srv interface{}) baseapp.EthQueryHandler {
	return func(ctx sdk.Context, req cmtrpctypes.RPCRequest) (abci.ResponseEthQuery, error) {
		in := new(types.QueryBalanceRequest)
		var params interface{}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return abci.ResponseEthQuery{}, err
		}
		in.Denom = "BNB" // only support BNB
		for _, p := range params.([]interface{}) {
			addr := p.(string)
			if _, err := sdk.AccAddressFromHexUnsafe(addr); err == nil {
				in.Address = addr
				break
			}
		}

		res, err := srv.(types.QueryServer).Balance(ctx, in)
		if err != nil {
			return abci.ResponseEthQuery{}, err
		}
		var amtBz []byte
		if res.Balance == nil || res.Balance.Amount.IsZero() {
			amtBz = big.NewInt(0).Bytes()
		} else {
			amtBz = res.Balance.Amount.BigInt().Bytes()
		}
		return abci.ResponseEthQuery{Response: amtBz}, nil
	}
}
