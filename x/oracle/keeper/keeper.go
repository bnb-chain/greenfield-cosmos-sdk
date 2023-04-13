package keeper

import (
	"bytes"
	"encoding/hex"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/willf/bitset"

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
	var validatorIndex int64 = -1
	var vldr stakingtypes.Validator
	for index, validator := range validators {
		if validator.RelayerAddress == relayer.String() {
			validatorIndex = int64(index)
			vldr = validator
			break
		}
	}

	if validatorIndex < 0 {
		return false, sdkerrors.Wrapf(types.ErrNotRelayer, "sender(%s) is not a relayer", relayer.String())
	}

	inturnRelayerTimeout, relayerInterval := k.GetRelayerParams(ctx)

	// check whether submitter of msgClaim is an in-turn relayer
	inturnRelayerBlsKey, _, err := k.getInturnRelayer(ctx, relayerInterval)
	if err != nil {
		return false, err
	}

	if bytes.Equal(inturnRelayerBlsKey, vldr.BlsKey) {
		return true, nil
	}

	// It is possible that claim comes from out-turn relayers when exceeding the inturnRelayerTimeout, all other
	// relayers can relay within the in-turn relayer's current interval
	curTime := ctx.BlockTime().Unix()
	if uint64(curTime) < claimTimestamp {
		return false, nil
	}

	return uint64(curTime)-claimTimestamp >= inturnRelayerTimeout, nil
}

// CheckClaim checks the bls signature
func (k Keeper) CheckClaim(ctx sdk.Context, claim *types.MsgClaim) (sdk.AccAddress, []sdk.AccAddress, error) {
	relayer, err := sdk.AccAddressFromHexUnsafe(claim.FromAddress)
	if err != nil {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrInvalidAddress, "from address (%s) is invalid", claim.FromAddress)
	}

	historicalInfo, ok := k.StakingKeeper.GetHistoricalInfo(ctx, ctx.BlockHeight())
	if !ok {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrValidatorSet, "get historical validators failed")
	}
	validators := historicalInfo.Valset

	isValid, err := k.IsRelayerValid(ctx, relayer, validators, claim.Timestamp)
	if err != nil {
		return sdk.AccAddress{}, nil, err
	}

	if !isValid {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrRelayerNotInTurn, "relayer(%s) is not in turn", claim.FromAddress)
	}

	validatorsBitSet := bitset.From(claim.VoteAddressSet)
	if validatorsBitSet.Count() > uint(len(validators)) {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrValidatorSet, "number of validator set is larger than validators")
	}

	signedRelayers := make([]sdk.AccAddress, 0, validatorsBitSet.Count())
	votedPubKeys := make([]bls.PublicKey, 0, validatorsBitSet.Count())
	for index, val := range validators {
		if !validatorsBitSet.Test(uint(index)) {
			continue
		}

		signedRelayers = append(signedRelayers, sdk.MustAccAddressFromHex(val.RelayerAddress))

		votePubKey, err := bls.PublicKeyFromBytes(val.BlsKey)
		if err != nil {
			return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrBlsPubKey, "BLS public key converts failed: %v", err)
		}
		votedPubKeys = append(votedPubKeys, votePubKey)
	}

	// The valid voted validators should be no less than 2/3 validators.
	if len(votedPubKeys) <= len(validators)*2/3 {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrBlsVotesNotEnough, "not enough validators voted, need: %d, voted: %d", len(validators)*2/3, len(votedPubKeys))
	}

	// Verify the aggregated signature.
	aggSig, err := bls.SignatureFromBytes(claim.AggSignature)
	if err != nil {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrInvalidBlsSignature, "BLS signature converts failed: %v", err)
	}

	if !aggSig.FastAggregateVerify(votedPubKeys, claim.GetBlsSignBytes()) {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrInvalidBlsSignature, "signature verify failed")
	}

	return relayer, signedRelayers, nil
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
	historicalInfo, ok := k.StakingKeeper.GetHistoricalInfo(ctx, ctx.BlockHeight())
	if !ok {
		return nil, nil, sdkerrors.Wrapf(types.ErrValidatorSet, "get historical validators failed")
	}
	validators := historicalInfo.Valset

	validatorsSize := len(validators)

	// totalIntervals is sum of intervals from all relayers
	totalIntervals := relayerInterval * uint64(validatorsSize)

	curTimeStamp := uint64(ctx.BlockTime().Unix())

	// remainder is used to locate inturn relayer.
	remainder := curTimeStamp % totalIntervals
	inTurnRelayerIndex := remainder / relayerInterval

	start := curTimeStamp - (remainder - inTurnRelayerIndex*relayerInterval)
	end := start + relayerInterval

	inturnRelayer := validators[inTurnRelayerIndex]

	return inturnRelayer.BlsKey, &types.RelayInterval{
		Start: start,
		End:   end,
	}, nil
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
