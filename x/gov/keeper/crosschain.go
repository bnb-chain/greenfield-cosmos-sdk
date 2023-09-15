package keeper

import (
	"encoding/hex"
	"math/big"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func (k Keeper) RegisterCrossChainSyncParamsApp() error {
	if err := k.crossChainKeeper.RegisterChannel(types.SyncParamsChannel, types.SyncParamsChannelID, SyncParamsApp{keeper: k}); err != nil {
		return err
	}

	return nil
}

type SyncParamsApp struct {
	keeper Keeper
}

func (app SyncParamsApp) ExecuteSynPackage(ctx sdk.Context, _ *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	app.keeper.Logger(ctx).Error("received sync params sync package", "payload", hex.EncodeToString(payload))
	return sdk.ExecuteResult{}
}

func (app SyncParamsApp) ExecuteAckPackage(ctx sdk.Context, _ *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	app.keeper.Logger(ctx).Error("received sync params in ack package", "payload", hex.EncodeToString(payload))
	return sdk.ExecuteResult{}
}

func (app SyncParamsApp) ExecuteFailAckPackage(ctx sdk.Context, _ *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	app.keeper.Logger(ctx).Error("received sync params fail ack package", "payload", hex.EncodeToString(payload))
	return sdk.ExecuteResult{}
}

func (k Keeper) SyncParams(ctx sdk.Context, destChainId sdk.ChainID, cpc govv1.CrossChainParamsChange) error {
	if err := cpc.ValidateBasic(); err != nil {
		return err
	}
	values := make([]byte, 0)
	addresses := make([]byte, 0)

	for i, v := range cpc.Values {
		var value []byte
		var err error
		if cpc.Key == types.KeyUpgrade {
			value = sdk.MustAccAddressFromHex(v)
		} else {
			value, err = hex.DecodeString(v)
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidValue, "value is not valid %s", v)
			}
		}
		values = append(values, value...)
		addr := sdk.MustAccAddressFromHex(cpc.Targets[i])
		addresses = append(addresses, addr.Bytes()...)
	}

	pack := types.SyncParamsPackage{
		Key:    cpc.Key,
		Value:  values,
		Target: addresses,
	}

	encodedPackage, err := pack.Serialize()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidSyncParamPackage, "fail to serialize, err: %s", err.Error())
	}

	if !k.crossChainKeeper.IsDestChainSupported(destChainId) {
		return sdkerrors.Wrapf(types.ErrChainNotSupported, "destination chain (%d) is not supported", destChainId)
	}

	_, err = k.crossChainKeeper.CreateRawIBCPackageWithFee(
		ctx,
		destChainId,
		types.SyncParamsChannelID,
		sdk.SynCrossChainPackageType,
		encodedPackage,
		big.NewInt(0),
		big.NewInt(0),
	)
	return err
}
