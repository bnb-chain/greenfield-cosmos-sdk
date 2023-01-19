package server

import (
	"github.com/ethereum/go-ethereum/rpc"
	rpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server/jsonrpc/backend"
	"github.com/cosmos/cosmos-sdk/server/jsonrpc/namespaces/ethereum/eth"
	"github.com/cosmos/cosmos-sdk/server/jsonrpc/namespaces/ethereum/net"
)

// RPC namespaces and API version
const (
	// Ethereum namespaces

	EthNamespace = "eth"
	NetNamespace = "net"

	apiVersion = "1.0"
)

// APICreator creates the JSON-RPC API implementations.
type APICreator = func(
	ctx *Context,
	clientCtx client.Context,
	tendermintWebsocketClient *rpcclient.WSClient,
) []rpc.API

// apiCreators defines the JSON-RPC API namespaces.
var apiCreators map[string]APICreator

func init() {
	apiCreators = map[string]APICreator{
		EthNamespace: func(ctx *Context,
			clientCtx client.Context,
			tmWSClient *rpcclient.WSClient,
		) []rpc.API {
			evmBackend := backend.NewBackend(ctx.Viper, ctx.Logger, clientCtx)
			return []rpc.API{
				{
					Namespace: EthNamespace,
					Version:   apiVersion,
					Service:   eth.NewPublicAPI(ctx.Logger, evmBackend),
					Public:    true,
				},
			}
		},
		NetNamespace: func(_ *Context, clientCtx client.Context, _ *rpcclient.WSClient) []rpc.API {
			return []rpc.API{
				{
					Namespace: NetNamespace,
					Version:   apiVersion,
					Service:   net.NewPublicAPI(clientCtx),
					Public:    true,
				},
			}
		},
	}
}

// GetRPCAPIs returns the list of all APIs
func GetRPCAPIs(ctx *Context,
	clientCtx client.Context,
	tmWSClient *rpcclient.WSClient,
	selectedAPIs []string,
) []rpc.API {
	var apis []rpc.API

	for _, ns := range selectedAPIs {
		if creator, ok := apiCreators[ns]; ok {
			apis = append(apis, creator(ctx, clientCtx, tmWSClient)...)
		} else {
			ctx.Logger.Error("invalid namespace value", "namespace", ns)
		}
	}

	return apis
}
