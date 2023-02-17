package keeper

import (
	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	types "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"math/big"
)

func (k Keeper) SyncParams(ctx sdk.Context, p *types.ParameterChangeProposal) error {

	key := p.Changes[0].Key
	values := make([]byte, 0)
	addresses := make([]byte, 0)

	if len(p.Changes) != len(p.Addresses) {
		return sdkerrors.Wrap(types.ErrAddressSizeNotMatch, "number of addresses not match")
	}
	if key != types.KeyUpgrade && len(p.Changes) > 1 {
		return sdkerrors.Wrap(types.ErrExceedParamsChangeLimit, "only single parameter change allowed")
	}

	for i, c := range p.Changes {
		values = append(values, []byte(c.Value)...)
		adr, err := sdk.AccAddressFromHexUnsafe(p.Addresses[i])
		if err != nil {
			return sdkerrors.Wrapf(types.ErrAddressNotValid, "smart contract address is not valid %s", p.Addresses[i])
		}
		addresses = append(addresses, adr.Bytes()...)
	}

	relayerFeeAmount, err := k.GetSyncParamsRelayerFee(ctx)
	if err != nil {
		return err
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
	_, err = (*k.CrossChainKeeper).CreateRawIBCPackageWithFee(
		ctx,
		types.SyncParamsChannelID,
		sdk.SynCrossChainPackageType,
		encodedPackage,
		relayerFeeAmount,
		big.NewInt(0),
	)

	if err != nil {
		return err
	}
	return nil
}

// GetSyncParamsRelayerFee gets the sync params change relayer fee params
func (k Keeper) GetSyncParamsRelayerFee(ctx sdk.Context) (*big.Int, error) {
	var syncParamsRelayerFeeParam string
	ss, ok := k.GetSubspace(types.BridgeSubspace)
	if !ok {
		return nil, sdkerrors.Wrap(types.ErrUnknownSubspace, types.BridgeSubspace)
	}
	ss.Get(ctx, types.KeySyncParamsRelayerFee, &syncParamsRelayerFeeParam)
	relayerFee, valid := big.NewInt(0).SetString(syncParamsRelayerFeeParam, 10)
	if !valid {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRelayerFee, "invalid relayer fee: %s", syncParamsRelayerFeeParam)
	}
	return relayerFee, nil
}

// Need these in order to register paramsKeeper to be a CrosschainApp so that it can register channel(3)

func (k Keeper) ExecuteSynPackage(ctx sdk.Context, appCtx *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	return sdk.ExecuteResult{}
}

func (k Keeper) ExecuteAckPackage(ctx sdk.Context, header *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	return sdk.ExecuteResult{}
}

func (k Keeper) ExecuteFailAckPackage(ctx sdk.Context, header *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	return sdk.ExecuteResult{}
}
