package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// AddVote adds a vote on a specific proposal
func (k Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options v1.WeightedVoteOptions, metadata string) error {
	// Check if proposal is in voting period.
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.VotingPeriodProposalKey(proposalID)) {
		return sdkerrors.Wrapf(types.ErrInactiveProposal, "%d", proposalID)
	}

	err := k.assertMetadataLength(metadata)
	if err != nil {
		return err
	}

	for _, option := range options {
		if !v1.ValidWeightedVoteOption(*option) {
			return sdkerrors.Wrap(types.ErrInvalidVote, option.String())
		}
	}

	vote := v1.NewVote(proposalID, voterAddr, options, metadata)
	k.SetVote(ctx, vote)

	// called after a vote on a proposal is cast
	k.Hooks().AfterProposalVote(ctx, proposalID, voterAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalVote,
			sdk.NewAttribute(types.AttributeKeyOption, options.String()),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

// GetAllVotes returns all the votes from the store
func (k Keeper) GetAllVotes(ctx sdk.Context) (votes v1.Votes) {
	k.IterateAllVotes(ctx, func(vote v1.Vote) bool {
		votes = append(votes, &vote)
		return false
	})
	return
}

// GetVotes returns all the votes from a proposal
func (k Keeper) GetVotes(ctx sdk.Context, proposalID uint64) (votes v1.Votes) {
	k.IterateVotes(ctx, proposalID, func(vote v1.Vote) bool {
		votes = append(votes, &vote)
		return false
	})
	return
}

// GetVote gets the vote from an address on a specific proposal
func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (vote v1.Vote, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VoteKey(proposalID, voterAddr))
	if bz == nil {
		return vote, false
	}

	k.cdc.MustUnmarshal(bz, &vote)

	return vote, true
}

// SetVote sets a Vote to the gov store
func (k Keeper) SetVote(ctx sdk.Context, vote v1.Vote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&vote)
	addr := sdk.MustAccAddressFromHex(vote.Voter)

	store.Set(types.VoteKey(vote.ProposalId, addr), bz)
}

// IterateAllVotes iterates over all the stored votes and performs a callback function
func (k Keeper) IterateAllVotes(ctx sdk.Context, cb func(vote v1.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote v1.Vote
		k.cdc.MustUnmarshal(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// IterateVotes iterates over all the proposals votes and performs a callback function
func (k Keeper) IterateVotes(ctx sdk.Context, proposalID uint64, cb func(vote v1.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.VotesKey(proposalID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote v1.Vote
		k.cdc.MustUnmarshal(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// deleteVote deletes a vote from a given proposalID and voter from the store
func (k Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.VoteKey(proposalID, voterAddr))
}
