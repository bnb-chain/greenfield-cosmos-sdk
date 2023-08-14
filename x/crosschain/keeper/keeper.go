package keeper

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crosschain/types"
)

// Keeper of the cross chain store
type Keeper struct {
	cdc codec.BinaryCodec

	cfg      *crossChainConfig
	storeKey storetypes.StoreKey

	authority string

	stakingKeeper types.StakingKeeper
	bankKeeper    types.BankKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, authority string,
	stakingKeeper types.StakingKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      key,
		cfg:           newCrossChainCfg(),
		authority:     authority,
		stakingKeeper: stakingKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Logger inits the logger for cross chain module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

// InitGenesis inits the genesis state of cross chain module
func (k Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper) {
	k.Logger(ctx).Info("set cross chain genesis state", "params", state.Params.String())
	k.SetParams(ctx, state.Params)

	params := k.GetParams(ctx)

	// for testing
	if !params.InitModuleBalance.IsNil() && params.InitModuleBalance.GT(sdk.ZeroInt()) {
		bondDenom := stakingKeeper.BondDenom(ctx)

		err := bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.Coin{
			Denom:  bondDenom,
			Amount: params.InitModuleBalance,
		}})
		if err != nil {
			panic(fmt.Sprintf("mint initial cross chain module balance error, err=%s", err.Error()))
		}
	}
}

// GetParams returns the current x/crosschain module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return p
	}

	k.cdc.MustUnmarshal(bz, &p)
	return p
}

// SetParams sets the params of cross chain module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}

// UpdatePermissions updates the permission of channels
func (k Keeper) UpdatePermissions(ctx sdk.Context, permissions []*types.ChannelPermission) error {
	for _, permission := range permissions {
		if !k.IsDestChainSupported(sdk.ChainID(permission.DestChainId)) {
			return fmt.Errorf("dest chain %d is not supported", permission.DestChainId)
		}
		if !k.IsChannelSupported(sdk.ChannelID(permission.ChannelId)) {
			return fmt.Errorf("channel %d is not supported", permission.ChannelId)
		}
		if sdk.ChannelPermission(permission.Permission) != sdk.ChannelAllow && sdk.ChannelPermission(permission.Permission) != sdk.ChannelForbidden {
			return fmt.Errorf("permission %d is not supported", permission.Permission)
		}

		k.SetChannelSendPermission(ctx, sdk.ChainID(permission.DestChainId), sdk.ChannelID(permission.ChannelId), sdk.ChannelPermission(permission.Permission))
	}
	return nil
}

// CreateRawIBCPackageWithFee creates a cross chain package with given cross chain fee
func (k Keeper) CreateRawIBCPackageWithFee(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID,
	packageType sdk.CrossChainPackageType, packageLoad []byte, relayerFee, ackRelayerFee *big.Int,
) (uint64, error) {
	if packageType == sdk.SynCrossChainPackageType && k.GetChannelSendPermission(ctx, destChainId, channelID) != sdk.ChannelAllow {
		return 0, fmt.Errorf("channel %d is not allowed to write syn package", channelID)
	}

	sequence := k.GetSendSequence(ctx, destChainId, channelID)
	key := types.BuildCrossChainPackageKey(k.GetSrcChainID(), destChainId, channelID, sequence)
	kvStore := ctx.KVStore(k.storeKey)
	if kvStore.Has(key) {
		return 0, fmt.Errorf("duplicated sequence")
	}

	// Assemble the package header
	packageHeader := sdk.EncodePackageHeader(sdk.PackageHeader{
		PackageType:   packageType,
		Timestamp:     uint64(ctx.BlockTime().Unix()),
		RelayerFee:    relayerFee,
		AckRelayerFee: ackRelayerFee,
	})

	kvStore.Set(key, append(packageHeader, packageLoad...))

	k.IncrSendSequence(ctx, destChainId, channelID)

	err := ctx.EventManager().EmitTypedEvent(&types.EventCrossChain{
		SrcChainId:    uint32(k.GetSrcChainID()),
		DestChainId:   uint32(destChainId),
		ChannelId:     uint32(channelID),
		Sequence:      sequence,
		PackageType:   uint32(packageType),
		Timestamp:     uint64(ctx.BlockTime().Unix()),
		PackageLoad:   hex.EncodeToString(packageLoad),
		RelayerFee:    relayerFee.String(),
		AckRelayerFee: ackRelayerFee.String(),
	})
	if err != nil {
		return 0, err
	}

	return sequence, nil
}

// RegisterChannel register a channel to the cross chain module with the cross chain app
func (k Keeper) RegisterChannel(name string, id sdk.ChannelID, app sdk.CrossChainApplication) error {
	_, ok := k.cfg.nameToChannelID[name]
	if ok {
		return fmt.Errorf("duplicated channel name")
	}
	_, ok = k.cfg.channelIDToName[id]
	if ok {
		return fmt.Errorf("duplicated channel id")
	}
	if app == nil {
		return fmt.Errorf("nil cross chain app")
	}
	k.cfg.nameToChannelID[name] = id
	k.cfg.channelIDToName[id] = name
	k.cfg.channelIDToApp[id] = app
	return nil
}

// IsDestChainSupported returns the support status of a dest chain
func (k Keeper) IsDestChainSupported(chainID sdk.ChainID) bool {
	if chainID == k.cfg.destBscChainId {
		return true
	}
	if k.cfg.destOpChainId != 0 && chainID == k.cfg.destOpChainId {
		return true
	}
	return false
}

// IsChannelSupported returns the support status of a channel
func (k Keeper) IsChannelSupported(channelId sdk.ChannelID) bool {
	_, ok := k.cfg.channelIDToName[channelId]
	return ok
}

// SetChannelSendPermission sets the channel send permission
func (k Keeper) SetChannelSendPermission(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID, permission sdk.ChannelPermission) {
	kvStore := ctx.KVStore(k.storeKey)
	kvStore.Set(types.BuildChannelPermissionKey(destChainID, channelID), []byte{byte(permission)})
}

// GetChannelSendPermission gets the channel send permission by channel id
func (k Keeper) GetChannelSendPermission(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) sdk.ChannelPermission {
	kvStore := ctx.KVStore(k.storeKey)
	bz := kvStore.Get(types.BuildChannelPermissionKey(destChainID, channelID))
	if bz == nil {
		return sdk.ChannelForbidden
	}
	return sdk.ChannelPermission(bz[0])
}

// SetSrcChainID sets the current chain id
func (k Keeper) SetSrcChainID(srcChainID sdk.ChainID) {
	k.cfg.srcChainID = srcChainID
}

// GetSrcChainID gets the current  chain id
func (k Keeper) GetSrcChainID() sdk.ChainID {
	return k.cfg.srcChainID
}

// SetDestBscChainID sets the destination chain id
func (k Keeper) SetDestBscChainID(destChainId sdk.ChainID) {
	k.cfg.destBscChainId = destChainId
}

// GetDestBscChainID gets the destination chain id of bsc
func (k Keeper) GetDestBscChainID() sdk.ChainID {
	return k.cfg.destBscChainId
}

// SetDestOpChainID sets the destination chain id of op chain
func (k Keeper) SetDestOpChainID(destChainId sdk.ChainID) {
	k.cfg.destOpChainId = destChainId
}

// GetDestOpChainID gets the destination chain id of op chain
func (k Keeper) GetDestOpChainID() sdk.ChainID {
	return k.cfg.destOpChainId
}

// GetCrossChainPackage returns the ibc package by sequence
func (k Keeper) GetCrossChainPackage(ctx sdk.Context, destChainId sdk.ChainID, channelId sdk.ChannelID, sequence uint64) ([]byte, error) {
	kvStore := ctx.KVStore(k.storeKey)
	key := types.BuildCrossChainPackageKey(k.GetSrcChainID(), destChainId, channelId, sequence)
	return kvStore.Get(key), nil
}

// GetSendSequence returns the sending sequence of the channel
func (k Keeper) GetSendSequence(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID) uint64 {
	return k.getSequence(ctx, destChainId, channelID, types.PrefixForSendSequenceKey)
}

// IncrSendSequence increases the sending sequence of the channel
func (k Keeper) IncrSendSequence(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID) {
	k.incrSequence(ctx, destChainId, channelID, types.PrefixForSendSequenceKey)
}

// GetReceiveSequence returns the receiving sequence of the channel
func (k Keeper) GetReceiveSequence(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID) uint64 {
	return k.getSequence(ctx, destChainId, channelID, types.PrefixForReceiveSequenceKey)
}

// IncrReceiveSequence increases the receiving sequence of the channel
func (k Keeper) IncrReceiveSequence(ctx sdk.Context, destChainId sdk.ChainID, channelID sdk.ChannelID) {
	k.incrSequence(ctx, destChainId, channelID, types.PrefixForReceiveSequenceKey)
}

// getSequence returns the sequence with a prefix
func (k Keeper) getSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID, prefix []byte) uint64 {
	kvStore := ctx.KVStore(k.storeKey)
	bz := kvStore.Get(types.BuildChannelSequenceKey(destChainID, channelID, prefix))
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// incrSequence increases the sequence with a prefix
func (k Keeper) incrSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID, prefix []byte) {
	var sequence uint64
	kvStore := ctx.KVStore(k.storeKey)
	bz := kvStore.Get(types.BuildChannelSequenceKey(destChainID, channelID, prefix))
	if bz == nil {
		sequence = 0
	} else {
		sequence = binary.BigEndian.Uint64(bz)
	}

	sequenceBytes := make([]byte, types.SequenceLength)
	binary.BigEndian.PutUint64(sequenceBytes, sequence+1)
	kvStore.Set(types.BuildChannelSequenceKey(destChainID, channelID, prefix), sequenceBytes)
}

// GetCrossChainApp returns the cross chain app by channel id
func (k Keeper) GetCrossChainApp(channelID sdk.ChannelID) sdk.CrossChainApplication {
	return k.cfg.channelIDToApp[channelID]
}

func (k Keeper) MintModuleAccountTokens(ctx sdk.Context, amount math.Int) error {
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.Coin{
		Denom:  bondDenom,
		Amount: amount,
	}})
	if err != nil {
		return fmt.Errorf("mint cross chain module amount error, err=%s", err.Error())
	}
	return nil
}
