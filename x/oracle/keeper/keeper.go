package keeper

import (
	"encoding/hex"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	StakingKeeper    types.StakingKeeper
	CrossChainKeeper types.CrossChainKeeper
	BankKeeper       types.BankKeeper

	feeCollectorName string // name of the FeeCollector ModuleAccount
	authority        string
}

func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, feeCollector, authority string,
	crossChainKeeper types.CrossChainKeeper, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper,
) Keeper {
	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		feeCollectorName: feeCollector,
		authority:        authority,

		CrossChainKeeper: crossChainKeeper,
		BankKeeper:       bankKeeper,
		StakingKeeper:    stakingKeeper,
	}
}

// Logger inits the logger for cross chain module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// InitGenesis inits the genesis state of oracle module
func (k Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) {
	k.SetParams(ctx, state.Params)
}

// SetParams sets the params of oarcle module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetRelayerParams returns the relayer timeout, relayer interval for oracle claim
func (k Keeper) GetRelayerParams(ctx sdk.Context) (uint64, uint64) {
	params := k.GetParams(ctx)
	return params.RelayerTimeout, params.RelayerInterval
}

// GetRelayerRewardShare returns the relayer reward share
func (k Keeper) GetRelayerRewardShare(ctx sdk.Context) uint32 {
	params := k.GetParams(ctx)
	return params.RelayerRewardShare
}

// IsRelayerValid returns true if the relayer is valid and allowed to send the claim message
func (k Keeper) IsRelayerValid(ctx sdk.Context, relayer sdk.AccAddress, validators []stakingtypes.Validator, claimTimestamp uint64) (bool, error) {
	// todo: implement this
	return true, nil
}

// CheckClaim checks the bls signature
func (k Keeper) CheckClaim(ctx sdk.Context, claim *types.MsgClaim) (sdk.AccAddress, []sdk.AccAddress, error) {
	// todo: implement this

	return sdk.AccAddress{}, []sdk.AccAddress{}, nil
}

// GetParams returns the current params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) getInturnRelayer(ctx sdk.Context, relayerInterval uint64) ([]byte, *types.RelayInterval, error) {
	// todo: implement this
	return nil, nil, nil
}

func (k Keeper) GetInturnRelayer(ctx sdk.Context, relayerInterval uint64) (*types.QueryInturnRelayerResponse, error) {
	blsKey, interval, err := k.getInturnRelayer(ctx, relayerInterval)
	if err != nil {
		return nil, err
	}
	res := &types.QueryInturnRelayerResponse{
		BlsPubKey:     hex.EncodeToString(blsKey),
		RelayInterval: interval,
	}
	return res, nil
}
