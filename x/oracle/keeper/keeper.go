package keeper

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/willf/bitset"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace

	StakingKeeper    types.StakingKeeper
	CrossChainKeeper types.CrossChainKeeper
	BankKeeper       types.BankKeeper

	feeCollectorName string // name of the FeeCollector ModuleAccount
}

func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace, feeCollector string,
	crossChainKeeper types.CrossChainKeeper, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace,
		feeCollectorName: feeCollector,

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
	k.Logger(ctx).Info("set oracle genesis state", "params", state.Params.String())
	k.SetParams(ctx, state.Params)
}

// SetParams sets the params of oarcle module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetRelayerParam returns the relayer timeout and backoff time for oracle claim
func (k Keeper) GetRelayerParam(ctx sdk.Context) (uint64, uint64) {
	var relayerTimeoutParam uint64
	var relayerBackoffTimeParam uint64
	k.paramSpace.Get(ctx, types.KeyParamRelayerTimeout, &relayerTimeoutParam)
	k.paramSpace.Get(ctx, types.KeyParamRelayerBackoffTime, &relayerBackoffTimeParam)
	return relayerTimeoutParam, relayerBackoffTimeParam
}

// GetRelayerRewardShare returns the relayer reward share
func (k Keeper) GetRelayerRewardShare(ctx sdk.Context) uint32 {
	var relayerRewardShare uint32
	k.paramSpace.Get(ctx, types.KeyParamRelayerRewardShare, &relayerRewardShare)
	return relayerRewardShare
}

// IsRelayerValid returns true if the relayer is valid and allowed to send the claim message
func (k Keeper) IsRelayerValid(ctx sdk.Context, relayer sdk.AccAddress, validators []stakingtypes.Validator, claimTimestamp uint64) (bool, error) {
	var validatorIndex int64 = -1
	for index, validator := range validators {
		if validator.RelayerAddress == relayer.String() {
			validatorIndex = int64(index)
			break
		}
	}

	if validatorIndex < 0 {
		return false, sdkerrors.Wrapf(types.ErrNotRelayer, fmt.Sprintf("sender(%s) is not a relayer", relayer.String()))
	}

	// check inturn validator index
	inturnValidatorIndex := claimTimestamp % uint64(len(validators))

	curTime := ctx.BlockTime().Unix()
	relayerTimeout, relayerBackoffTime := k.GetRelayerParam(ctx)

	// inturn validator can always relay package
	if uint64(validatorIndex) == inturnValidatorIndex {
		return true, nil
	}

	// not inturn validators can not relay in the timeout duration
	if uint64(curTime)-claimTimestamp <= relayerTimeout {
		return false, nil
	}
	validatorDistance := (validatorIndex - int64(inturnValidatorIndex) + int64(len(validators))) % int64(len(validators))
	return curTime > int64(claimTimestamp+relayerTimeout)+(validatorDistance-1)*int64(relayerBackoffTime), nil
}

// CheckClaim checks the bls signature
func (k Keeper) CheckClaim(ctx sdk.Context, claim *types.MsgClaim) (sdk.AccAddress, []string, error) {
	relayer, err := sdk.AccAddressFromHexUnsafe(claim.FromAddress)
	if err != nil {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrInvalidAddress, fmt.Sprintf("from address (%s) is invalid", claim.FromAddress))
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
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrRelayerNotInTurn, fmt.Sprintf("relayer(%s) is not in turn", claim.FromAddress))
	}

	validatorsBitSet := bitset.From(claim.VoteAddressSet)
	if validatorsBitSet.Count() > uint(len(validators)) {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrValidatorSet, "number of validator set is larger than validators")
	}

	signedRelayers := make([]string, 0, validatorsBitSet.Count())
	votedPubKeys := make([]bls.PublicKey, 0, validatorsBitSet.Count())
	for index, val := range validators {
		if !validatorsBitSet.Test(uint(index)) {
			continue
		}

		signedRelayers = append(signedRelayers, val.RelayerAddress)

		votePubKey, err := bls.PublicKeyFromBytes(val.RelayerBlsKey)
		if err != nil {
			return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrBlsPubKey, fmt.Sprintf("BLS public key converts failed: %v", err))
		}
		votedPubKeys = append(votedPubKeys, votePubKey)
	}

	// The valid voted validators should be no less than 2/3 validators.
	if len(votedPubKeys) <= len(validators)*2/3 {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrBlsVotesNotEnough, fmt.Sprintf("not enough validators voted, need: %d, voted: %d", len(validators)*2/3, len(votedPubKeys)))
	}

	// Verify the aggregated signature.
	aggSig, err := bls.SignatureFromBytes(claim.AggSignature)
	if err != nil {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrInvalidBlsSignature, fmt.Sprintf("BLS signature converts failed: %v", err))
	}

	if !aggSig.FastAggregateVerify(votedPubKeys, claim.GetBlsSignBytes()) {
		return sdk.AccAddress{}, nil, sdkerrors.Wrapf(types.ErrInvalidBlsSignature, "signature verify failed")
	}

	return relayer, signedRelayers, nil
}

// GetParams returns the current params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}
