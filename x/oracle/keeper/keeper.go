package keeper

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/oracle/metrics"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/willf/bitset"
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

// GetRelayerTimeoutParam returns the default relayer timeout for oracle claim
func (k Keeper) GetRelayerTimeoutParam(ctx sdk.Context) uint64 {
	var relayerTimeoutParam uint64
	k.paramSpace.Get(ctx, types.KeyParamRelayerTimeout, &relayerTimeoutParam)
	return relayerTimeoutParam
}

func (k Keeper) isValidatorInTurn(ctx sdk.Context, validators []stakingtypes.Validator, claim *types.MsgClaim) (bool, error) {
	var validatorIndex int64 = -1
	for index, validator := range validators {
		consAddr, err := validator.GetConsAddr()
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

	curTime := ctx.BlockTime().Unix()
	relayerTimeout := k.GetRelayerTimeoutParam(ctx)
	// check block time with package timestamp
	if uint64(curTime)-claim.Timestamp > relayerTimeout {
		return true, nil
	}

	// check inturn validator index
	inturnValidatorIndex := claim.Timestamp % uint64(len(validators))
	return uint64(validatorIndex) == inturnValidatorIndex, nil
}

// ProcessClaim checks the bls signature
func (k Keeper) ProcessClaim(ctx sdk.Context, claim *types.MsgClaim) error {
	validators := k.StakingKeeper.GetLastValidators(ctx)

	inturn, err := k.isValidatorInTurn(ctx, validators, claim)
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

// SendCoinsFromAccountToFeeCollector transfers amt to the fee collector account.
func (k Keeper) SendCoinsFromAccountToFeeCollector(ctx sdk.Context, senderAddr sdk.AccAddress, amt sdk.Coins) error {
	return k.BankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, k.feeCollectorName, amt)
}
