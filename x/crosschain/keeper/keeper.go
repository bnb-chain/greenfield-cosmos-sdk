package keeper

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crosschain/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the cross chain store
type Keeper struct {
	cdc codec.BinaryCodec

	cfg        *crossChainConfig
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   key,
		cfg:        newCrossChainCfg(),
		paramSpace: paramSpace,
	}
}

// Logger inits the logger for cross chain module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// InitGenesis inits the genesis state of cross chain module
func (k Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) {
	k.Logger(ctx).Info("set cross chain genesis state", "params", state.Params.String())
	k.SetParams(ctx, state.Params)
}

// GetRelayerFeeParam returns the default relayer fee for cross chain tx
func (k Keeper) GetRelayerFeeParam(ctx sdk.Context) (relayerFee *big.Int, err error) {
	var relayerFeeParam string
	k.paramSpace.Get(ctx, types.KeyParamRelayerFee, &relayerFeeParam)
	relayerFee, valid := relayerFee.SetString(relayerFeeParam, 10)
	if !valid {
		return nil, fmt.Errorf("invalid relayer fee: %s", relayerFeeParam)
	}
	return relayerFee, nil
}

// SetParams sets the params of cross chain module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// CreateRawIBCPackage creates a cross chain package with default cross chain fee
func (k Keeper) CreateRawIBCPackage(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID,
	packageType sdk.CrossChainPackageType, packageLoad []byte) (uint64, error) {

	relayerFee, err := k.GetRelayerFeeParam(ctx)
	if err != nil {
		return 0, fmt.Errorf("fail to load relayerFee, %v", err)
	}

	return k.CreateRawIBCPackageWithFee(ctx, destChainID, channelID, packageType, packageLoad, *relayerFee)
}

// CreateRawIBCPackageWithFee creates a cross chain package with given cross chain fee
func (k Keeper) CreateRawIBCPackageWithFee(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID,
	packageType sdk.CrossChainPackageType, packageLoad []byte, relayerFee big.Int) (uint64, error) {

	if packageType == sdk.SynCrossChainPackageType && k.GetChannelSendPermission(ctx, destChainID, channelID) != sdk.ChannelAllow {
		return 0, fmt.Errorf("channel %d is not allowed to write syn package", channelID)
	}

	sequence := k.GetSendSequence(ctx, destChainID, channelID)
	key := types.BuildIBCPackageKey(k.GetSrcChainID(), destChainID, channelID, sequence)
	kvStore := ctx.KVStore(k.storeKey)
	if kvStore.Has(key) {
		return 0, fmt.Errorf("duplicated sequence")
	}

	// Assemble the package header
	packageHeader := sdk.EncodePackageHeader(packageType, uint64(ctx.BlockTime().Unix()), relayerFee)

	kvStore.Set(key, append(packageHeader, packageLoad...))

	k.IncrSendSequence(ctx, destChainID, channelID)

	err := ctx.EventManager().EmitTypedEvent(&types.EventCrossChain{
		SrcChainId:  uint32(k.GetSrcChainID()),
		DestChainId: uint32(destChainID),
		ChannelId:   uint32(channelID),
		Sequence:    sequence,
		PackageType: uint32(packageType),
		Timestamp:   uint64(ctx.BlockTime().Unix()),
		PackageLoad: packageLoad,
		RelayerFee:  relayerFee.String(),
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

// RegisterDestChain registers a dest chain
func (k Keeper) RegisterDestChain(chainID sdk.ChainID) error {
	for _, chain := range k.cfg.destChains {
		if chainID == chain {
			return fmt.Errorf("duplicated destination chain chainID")
		}
	}

	k.cfg.destChains = append(k.cfg.destChains, chainID)
	return nil
}

// IsDestChainSupported returns the support status of a dest chain
func (k Keeper) IsDestChainSupported(chainID sdk.ChainID) bool {
	for _, chain := range k.cfg.destChains {
		if chainID == chain {
			return true
		}
	}
	return false
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

// GetIBCPackage returns the ibc package by sequence
func (k Keeper) GetIBCPackage(ctx sdk.Context, destChainID sdk.ChainID, channelId sdk.ChannelID, sequence uint64) ([]byte, error) {
	kvStore := ctx.KVStore(k.storeKey)
	key := types.BuildIBCPackageKey(k.GetSrcChainID(), destChainID, channelId, sequence)
	return kvStore.Get(key), nil
}

// GetSendSequence returns the sending sequence of the channel
func (k Keeper) GetSendSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) uint64 {
	return k.getSequence(ctx, destChainID, channelID, types.PrefixForSendSequenceKey)
}

// IncrSendSequence increases the sending sequence of the channel
func (k Keeper) IncrSendSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) {
	k.incrSequence(ctx, destChainID, channelID, types.PrefixForSendSequenceKey)
}

// GetReceiveSequence returns the receiving sequence of the channel
func (k Keeper) GetReceiveSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) uint64 {
	return k.getSequence(ctx, destChainID, channelID, types.PrefixForReceiveSequenceKey)
}

// IncrReceiveSequence increases the receiving sequence of the channel
func (k Keeper) IncrReceiveSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) {
	k.incrSequence(ctx, destChainID, channelID, types.PrefixForReceiveSequenceKey)
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

// GetParams returns the current params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// GetCrossChainApp returns the cross chain app by channel id
func (k Keeper) GetCrossChainApp(channelID sdk.ChannelID) sdk.CrossChainApplication {
	return k.cfg.channelIDToApp[channelID]
}
