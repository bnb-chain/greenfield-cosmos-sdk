package keeper

import (
	"encoding/hex"
	"math/big"

	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	types "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

func (k Keeper) SyncParams(ctx sdk.Context, p *types.ParameterChangeProposal) error {
	// this validates content and size of changes is not empty
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	if len(p.Changes) != len(p.Addresses) {
		return sdkerrors.Wrap(types.ErrAddressSizeNotMatch, "number of addresses not match")
	}

	key := p.Changes[0].Key
	values := make([]byte, 0)
	addresses := make([]byte, 0)

	if key != types.KeyUpgrade && len(p.Changes) > 1 {
		return sdkerrors.Wrap(types.ErrExceedParamsChangeLimit, "only single parameter change allowed")
	}

	for i, c := range p.Changes {
		if c.Key != key {
			return sdkerrors.Wrap(types.ErrInvalidPackage, "all changes key should be 'ungrade'")
		}
		values = append(values, []byte(c.Value)...)
		adr, err := sdk.AccAddressFromHexUnsafe(p.Addresses[i])
		if err != nil {
			return sdkerrors.Wrapf(types.ErrAddressNotValid, "smart contract address is not valid %s", p.Addresses[i])
		}
		addresses = append(addresses, adr.Bytes()...)
	}

	pack := types.SyncParamsPackage{
		Key:    key,
		Value:  values,
		Target: addresses,
	}

	encodedPackage, err := rlp.EncodeToBytes(pack)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidPackage, "encode sync params package error")
	}
	_, err = (*k.crossChainKeeper).CreateRawIBCPackageWithFee(
		ctx,
		types.SyncParamsChannelID,
		sdk.SynCrossChainPackageType,
		encodedPackage,
		big.NewInt(0),
		big.NewInt(0),
	)

	if err != nil {
		return err
	}
	return nil
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
