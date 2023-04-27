package baseapp

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtrpctypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// EthQueryHandler defines a function type which handles EVM json-rpc requests
type EthQueryHandler = func(ctx sdk.Context, req cmtrpctypes.RPCRequest) (abci.ResponseEthQuery, error)

// EthQueryRouter routes eth Query requests to handlers
type EthQueryRouter struct {
	routes map[string]EthQueryHandler
}

// NewEthQueryRouter creates a new EthQueryRouter
func NewEthQueryRouter() *EthQueryRouter {
	return &EthQueryRouter{
		routes: map[string]EthQueryHandler{},
	}
}

// Route returns the EthQueryHandler for a given query route path or nil
// if not found
func (e *EthQueryRouter) Route(path string) EthQueryHandler {
	handler, found := e.routes[path]
	if !found {
		return nil
	}
	return handler
}

// AddRoute adds a query path to the router with a given Querier. It will panic
// if a duplicate route is given. The route must be alphanumeric.
func (e *EthQueryRouter) AddRoute(route string, h EthQueryHandler) {
	if e.routes[route] != nil {
		panic(fmt.Sprintf("route %s has already been initialized", route))
	}

	e.routes[route] = h
}

// RegisterEthQueryBalanceHandler adds router for EthGetBalance with a given handlerGen and server. It will panic
// if a duplicate route is given. The route must be alphanumeric.
func (e *EthQueryRouter) RegisterEthQueryBalanceHandler(srv interface{}, handlerGen func(interface{}) EthQueryHandler) {
	if e.routes[EthGetBalance] != nil {
		panic(fmt.Sprintf("route %s has already been initialized", EthGetBalance))
	}

	e.AddRoute(EthGetBalance, handlerGen(srv))
}

// RegisterConstHandler adds router for constant eth query.
func (e *EthQueryRouter) RegisterConstHandler() {
	e.AddRoute(EthBlockNumber, blockNumberHandler)
	e.AddRoute(EthGetBlockByNumber, blockNumberHandler)
	e.AddRoute(EthNetworkID, chainIdHandler)
	e.AddRoute(EthChainID, chainIdHandler)
	e.AddRoute(NetVersion, chainIdHandler)
}

func blockNumberHandler(ctx sdk.Context, req cmtrpctypes.RPCRequest) (abci.ResponseEthQuery, error) {
	var res abci.ResponseEthQuery
	res.Response = big.NewInt(ctx.BlockHeight()).Bytes()
	return res, nil
}

func chainIdHandler(ctx sdk.Context, req cmtrpctypes.RPCRequest) (abci.ResponseEthQuery, error) {
	var res abci.ResponseEthQuery
	eip155ChainID, err := sdk.ParseChainID(ctx.ChainID())
	if err != nil {
		return res, errorsmod.Wrap(sdkerrors.ErrInvalidChainID, err.Error())
	}
	res.Response = eip155ChainID.Bytes()
	return res, nil
}
