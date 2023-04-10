package keeper

import (
	"encoding/hex"
	"math/big"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func (k Keeper) RegisterCrossChainSyncParamsApp() error {
	return k.crossChainKeeper.RegisterChannel(types.SyncParamsChannel, types.SyncParamsChannelID, k)
}

func (k Keeper) SyncParams(ctx sdk.Context, cpc govv1.CrossChainParamsChange) error {
	// this validates content and size of changes is not empty
	if err := cpc.ValidateBasic(); err != nil {
		return err
	}
	values := make([]byte, 0)
	addresses := make([]byte, 0)

	for i, v := range cpc.Values {
		var value []byte
		var err error
		if cpc.Key == types.KeyUpgrade {
			value, err = sdk.AccAddressFromHexUnsafe(v)
			if err != nil {
				return sdkerrors.Wrapf(types.ErrAddressNotValid, "smart contract address is not valid %s", v)
			}
		} else {
			value, err = hex.DecodeString(v)
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidValue, "value is not valid %s", v)
			}
		}
		values = append(values, value...)

		addr, err := sdk.AccAddressFromHexUnsafe(cpc.Targets[i])
		if err != nil {
			return sdkerrors.Wrapf(types.ErrAddressNotValid, "smart contract address is not valid %s", cpc.Targets[i])
		}
		addresses = append(addresses, addr.Bytes()...)
	}

	pack := types.SyncParamsPackage{
		Key:    cpc.Key,
		Value:  values,
		Target: addresses,
	}

	encodedPackage, err := rlp.EncodeToBytes(pack)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidUpgradeProposal, "encode sync params package error")
	}
	_, err = k.crossChainKeeper.CreateRawIBCPackageWithFee(
		ctx,
		types.SyncParamsChannelID,
		sdk.SynCrossChainPackageType,
		encodedPackage,
		big.NewInt(0),
		big.NewInt(0),
	)
	return err
}

// Need these in order to register paramsKeeper to be a CrosschainApp so that it can register channel(3)

func (k Keeper) ExecuteSynPackage(ctx sdk.Context, appCtx *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	k.Logger(ctx).Error("received sync params sync package", "payload", hex.EncodeToString(payload))
	return sdk.ExecuteResult{}
}

func (k Keeper) ExecuteAckPackage(ctx sdk.Context, header *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	k.Logger(ctx).Error("received sync params in ack package", "payload", hex.EncodeToString(payload))
	return sdk.ExecuteResult{}
}

func (k Keeper) ExecuteFailAckPackage(ctx sdk.Context, header *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	k.Logger(ctx).Error("received sync params fail ack package", "payload", hex.EncodeToString(payload))
	return sdk.ExecuteResult{}
}
