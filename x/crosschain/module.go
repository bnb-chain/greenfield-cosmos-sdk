package crosschain

import (
	"context"
	"encoding/json"
	"fmt"

	modulev1 "cosmossdk.io/api/cosmos/crosschain/module/v1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	abci "github.com/cometbft/cometbft/abci/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/crosschain/keeper"
	"github.com/cosmos/cosmos-sdk/x/crosschain/types"
	"github.com/cosmos/cosmos-sdk/x/params/client/cli"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the params module.
type AppModuleBasic struct{}

// Name returns the params module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the params module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// DefaultGenesis returns default genesis state as raw bytes for the params
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the params module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the params module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns no root tx command for the params module.
func (AppModuleBasic) GetTxCmd() *cobra.Command { return nil }

// GetQueryCmd returns no root query command for the params module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.NewQueryCmd()
}

func (am AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// AppModule implements an application module for the distribution module.
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		bankKeeper:     bankKeeper,
		stakingKeeper:  stakingKeeper,
	}
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs a no-op.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, &genesisState, am.bankKeeper, am.stakingKeeper)
	return []abci.ValidatorUpdate{}
}

// GenerateGenesisState performs a no-op.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {}

// RegisterServices registers a gRPC query service to respond to the
// module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterStoreDecoder doesn't register any type.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

// ExportGenesis performs a no-op.
func (am AppModule) ExportGenesis(_ sdk.Context, _ codec.JSONCodec) json.RawMessage {
	return nil
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

//
// App Wiring Setup
//

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type CrossChainInputs struct {
	depinject.In

	Config *modulev1.Module
	Key    *store.KVStoreKey
	Cdc    codec.Codec

	BankKeeper    types.BankKeeper
	StakingKeeper types.StakingKeeper
}

type CrossChainOutputs struct {
	depinject.Out

	CrossChainKeeper keeper.Keeper
	Module           appmodule.AppModule
}

func ProvideModule(in CrossChainInputs) CrossChainOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrHexAddress(in.Config.Authority)
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.Key,
		authority.String(),
		in.StakingKeeper,
		in.BankKeeper,
	)

	m := NewAppModule(k, in.BankKeeper, in.StakingKeeper)

	return CrossChainOutputs{CrossChainKeeper: k, Module: m}
}
