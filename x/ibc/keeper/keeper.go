package keeper

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/bsc"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the ibc store
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
	return Keeper{
		cdc:        cdc,
		storeKey:   key,
		paramSpace: paramSpace,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k *Keeper) GetRelayerFeeParam(ctx sdk.Context) (relayerFee *big.Int, err error) {
	var relayerFeeParam int64
	k.paramSpace.Get(ctx, types.KeyParamRelayerFee, &relayerFeeParam)
	relayerFee = bsc.ConvertBCAmountToBSCAmount(relayerFeeParam)
	return
}

func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k *Keeper) CreateRawIBCPackage(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID,
	packageType sdk.CrossChainPackageType, packageLoad []byte) (uint64, error) {

	relayerFee, err := k.GetRelayerFeeParam(ctx)
	if err != nil {
		return 0, fmt.Errorf("fail to load relayerFee, %v", err)
	}

	return k.CreateRawIBCPackageWithFee(ctx, destChainID, channelID, packageType, packageLoad, *relayerFee)
}

func (k *Keeper) GetChannelSendPermission(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) sdk.ChannelPermission {
	kvStore := ctx.KVStore(k.storeKey)
	bz := kvStore.Get(types.BuildChannelPermissionKey(destChainID, channelID))
	if bz == nil {
		return sdk.ChannelForbidden
	}
	return sdk.ChannelPermission(bz[0])
}

func (k *Keeper) CreateRawIBCPackageWithFee(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID,
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

	err := ctx.EventManager().EmitTypedEvent(&types.EventIBCPackage{
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

func (k *Keeper) RegisterChannel(name string, id sdk.ChannelID, app sdk.CrossChainApplication) error {
	_, ok := k.cfg.nameToChannelID[name]
	if ok {
		return fmt.Errorf("duplicated channel name")
	}
	_, ok = k.cfg.channelIDToName[id]
	if ok {
		return fmt.Errorf("duplicated channel id")
	}
	k.cfg.nameToChannelID[name] = id
	k.cfg.channelIDToName[id] = name
	k.cfg.channelIDToApp[id] = app
	return nil
}

func (k *Keeper) RegisterDestChain(name string, chainID sdk.ChainID) error {
	_, ok := k.cfg.destChainNameToID[name]
	if ok {
		return fmt.Errorf("duplicated destination chain name")
	}
	_, ok = k.cfg.destChainIDToName[chainID]
	if ok {
		return fmt.Errorf("duplicated destination chain chainID")
	}
	k.cfg.destChainNameToID[name] = chainID
	k.cfg.destChainIDToName[chainID] = name
	return nil
}

func (k *Keeper) SetChannelSendPermission(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID, permission sdk.ChannelPermission) {
	kvStore := ctx.KVStore(k.storeKey)
	kvStore.Set(types.BuildChannelPermissionKey(destChainID, channelID), []byte{byte(permission)})
}

func (k *Keeper) SetSrcChainID(srcChainID sdk.ChainID) {
	k.cfg.srcChainID = srcChainID
}

func (k *Keeper) GetSrcChainID() sdk.ChainID {
	return k.cfg.srcChainID
}

func (k *Keeper) GetDestChainID(name string) (sdk.ChainID, error) {
	destChainID, exist := k.cfg.destChainNameToID[name]
	if !exist {
		return sdk.ChainID(0), fmt.Errorf("non-existing destination chainName ")
	}
	return destChainID, nil
}

func (k *Keeper) GetIBCPackage(ctx sdk.Context, destChainID sdk.ChainID, channelId sdk.ChannelID, sequence uint64) ([]byte, error) {
	kvStore := ctx.KVStore(k.storeKey)
	key := types.BuildIBCPackageKey(k.GetSrcChainID(), destChainID, channelId, sequence)
	return kvStore.Get(key), nil
}

func (k *Keeper) GetSendSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) uint64 {
	return k.getSequence(ctx, destChainID, channelID, types.PrefixForSendSequenceKey)
}

func (k *Keeper) IncrSendSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) {
	k.incrSequence(ctx, destChainID, channelID, types.PrefixForSendSequenceKey)
}

func (k *Keeper) GetReceiveSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) uint64 {
	return k.getSequence(ctx, destChainID, channelID, types.PrefixForReceiveSequenceKey)
}

func (k *Keeper) IncrReceiveSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID) {
	k.incrSequence(ctx, destChainID, channelID, types.PrefixForReceiveSequenceKey)
}

func (k *Keeper) getSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID, prefix []byte) uint64 {
	kvStore := ctx.KVStore(k.storeKey)
	bz := kvStore.Get(types.BuildChannelSequenceKey(destChainID, channelID, prefix))
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func (k *Keeper) incrSequence(ctx sdk.Context, destChainID sdk.ChainID, channelID sdk.ChannelID, prefix []byte) {
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
