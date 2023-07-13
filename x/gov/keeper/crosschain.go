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
	return k.crossChainKeeper.RegisterChannel(types.SyncParamsChannel, types.SyncParamsChannelID, k)
}

func (k Keeper) SyncParams(ctx sdk.Context, cpc govv1.CrossChainParamsChange) error {
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

	encodedPackage := pack.MustSerialize()

	_, err := k.crossChainKeeper.CreateRawIBCPackageWithFee(
		ctx,
		k.crossChainKeeper.GetDestBscChainID(),
		types.SyncParamsChannelID,
		sdk.SynCrossChainPackageType,
		encodedPackage,
		big.NewInt(0),
		big.NewInt(0),
	)
	return err
}

func (k Keeper) ExecuteSynPackage(ctx sdk.Context, _ *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	k.Logger(ctx).Error("received sync params sync package", "payload", hex.EncodeToString(payload))
	return sdk.ExecuteResult{}
}

func (k Keeper) ExecuteAckPackage(ctx sdk.Context, _ *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	k.Logger(ctx).Error("received sync params in ack package", "payload", hex.EncodeToString(payload))
	return sdk.ExecuteResult{}
}

func (k Keeper) ExecuteFailAckPackage(ctx sdk.Context, _ *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	k.Logger(ctx).Error("received sync params fail ack package", "payload", hex.EncodeToString(payload))
	return sdk.ExecuteResult{}
}
