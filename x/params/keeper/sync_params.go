package keeper

import (
	"encoding/hex"
	"math/big"

	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	types "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

func (k Keeper) RegisterCrossChainSyncParamsApp() error {
	return (*k.crossChainKeeper).RegisterChannel(types.SyncParamsChannel, types.SyncParamsChannelID, k)
}

func (k Keeper) SyncParams(ctx sdk.Context, p *types.ParameterChangeProposal) error {
	// validates if change(s) is/are present, proposal content is valid for params change or contract upgrade
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	values := make([]byte, 0)
	addresses := make([]byte, 0)

	for i, c := range p.Changes {
		var value []byte
		var err error
		if c.Key == types.KeyUpgrade {
			value, err = sdk.AccAddressFromHexUnsafe(c.Value)
			if err != nil {
				return sdkerrors.Wrapf(types.ErrAddressNotValid, "smart contract address format is not valid, address=%s", c.Value)
			}
		} else {
			value, err = hex.DecodeString(c.Value)
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidValue, "ParamChange value is not valid, should be in hex format, value=%s", c.Value)
			}
		}
		values = append(values, value...)
		addr, err := sdk.AccAddressFromHexUnsafe(p.Addresses[i])
		if err != nil {
			return sdkerrors.Wrapf(types.ErrAddressNotValid, "smart contract address format is not valid, address=%s", p.Addresses[i])
		}
		addresses = append(addresses, addr.Bytes()...)
	}

	pack := types.SyncParamsPackage{
		Key:    p.Changes[0].Key,
		Value:  values,
		Target: addresses,
	}

	encodedPackage, err := rlp.EncodeToBytes(pack)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidUpgradeProposal, "encode sync params package error")
	}
	_, err = (*k.crossChainKeeper).CreateRawIBCPackageWithFee(
		ctx,
		types.SyncParamsChannelID,
		sdk.SynCrossChainPackageType,
		encodedPackage,
		big.NewInt(0),
		big.NewInt(0),
	)
	return err
}

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
