package keeper

import (
	"fmt"

	crosschaintypes "github.com/cosmos/cosmos-sdk/x/crosschain/types"

	sdkerrors "cosmossdk.io/errors"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/willf/bitset"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/metrics"
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
	Metrics          *metrics.Metrics

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

		Metrics: metrics.PrometheusMetrics(),
	}
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

func (k Keeper) IsValidatorInturn(ctx sdk.Context, validators []stakingtypes.Validator, claim *types.MsgClaim) (bool, error) {
	var validatorIndex int64 = -1
	for index, validator := range validators {
		consAddr, err := validator.GetConsAddr() // TODO: update this
		if err != nil {
			return false, err
		}
		if consAddr.String() == claim.FromAddress {
			validatorIndex = int64(index)
			break
		}
	}

	if validatorIndex < 0 {
		return false, sdkerrors.Wrapf(types.ErrNotValidator, fmt.Sprintf("sender is not validator"))
	}

	// check inturn validator index
	inturnValidatorIndex := claim.Timestamp % uint64(len(validators))

	curTime := ctx.BlockTime().Unix()
	relayerTimeout, relayerBackoffTime := k.GetRelayerParam(ctx)

	// check block time with package timestamp
	if uint64(curTime)-claim.Timestamp <= relayerTimeout {
		if uint64(validatorIndex) == inturnValidatorIndex {
			return true, nil
		}
		return false, nil
	}

	backoffIndex := (uint64(curTime)-claim.Timestamp-relayerTimeout-1)/relayerBackoffTime + 1

	return uint64(validatorIndex) == (inturnValidatorIndex+backoffIndex)%uint64(len(validators)), nil
}

// ProcessClaim checks the bls signature
func (k Keeper) ProcessClaim(ctx sdk.Context, claim *types.MsgClaim) error {
	validators := k.StakingKeeper.GetLastValidators(ctx)

	inturn, err := k.IsValidatorInturn(ctx, validators, claim)
	if err != nil {
		return err
	}
	if !inturn {
		return sdkerrors.Wrapf(types.ErrValidatorNotInTurn, fmt.Sprintf("validator is not in turn"))
	}

	validatorsBitSet := bitset.From(claim.VoteAddressSet)
	if validatorsBitSet.Count() > uint(len(validators)) {
		return sdkerrors.Wrapf(types.ErrValidatorSet, fmt.Sprintf("number of validator set is larger than validators"))
	}

	votedPubKeys := make([]bls.PublicKey, 0, validatorsBitSet.Count())
	for index, val := range validators {
		if !validatorsBitSet.Test(uint(index)) {
			continue
		}

		// TODO: confirm the pub key
		voteAddr, err := bls.PublicKeyFromBytes(val.ConsensusPubkey.Value)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrBlsPubKey, fmt.Sprintf("BLS public key converts failed: %v", err))

		}
		votedPubKeys = append(votedPubKeys, voteAddr)
	}

	// The valid voted validators should be no less than 2/3 validators.
	if len(votedPubKeys) <= len(validators)*2/3 {
		return sdkerrors.Wrapf(types.ErrBlsVotesNotEnough, fmt.Sprintf("not enough validators voted, need: %d, voted: %d", len(validators)*2/3, len(votedPubKeys)))
	}

	// Verify the aggregated signature.
	aggSig, err := bls.SignatureFromBytes(claim.AggSignature)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidBlsSignature, fmt.Sprintf("BLS signature converts failed: %v", err))
	}

	if !aggSig.FastAggregateVerify(votedPubKeys, claim.GetBlsSignBytes()) {
		return sdkerrors.Wrapf(types.ErrInvalidBlsSignature, fmt.Sprintf("signature verify failed"))
	}

	return nil
}

// SendCoinsToFeeCollector transfers amt to the fee collector account.
func (k Keeper) SendCoinsToFeeCollector(ctx sdk.Context, amt sdk.Coins) error {
	return k.BankKeeper.SendCoinsFromModuleToModule(ctx, crosschaintypes.ModuleName, k.feeCollectorName, amt)
}
