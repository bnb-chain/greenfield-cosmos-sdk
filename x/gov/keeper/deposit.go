package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// GetDeposit gets the deposit of a specific depositor on a specific proposal
func (k Keeper) GetDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) (deposit v1.Deposit, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DepositKey(proposalID, depositorAddr))
	if bz == nil {
		return deposit, false
	}

	k.cdc.MustUnmarshal(bz, &deposit)

	return deposit, true
}

// SetDeposit sets a Deposit to the gov store
func (k Keeper) SetDeposit(ctx sdk.Context, deposit v1.Deposit) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&deposit)
	depositor := sdk.MustAccAddressFromHex(deposit.Depositor)

	store.Set(types.DepositKey(deposit.ProposalId, depositor), bz)
}

// GetAllDeposits returns all the deposits from the store
func (k Keeper) GetAllDeposits(ctx sdk.Context) (deposits v1.Deposits) {
	k.IterateAllDeposits(ctx, func(deposit v1.Deposit) bool {
		deposits = append(deposits, &deposit)
		return false
	})

	return
}

// GetDeposits returns all the deposits of a proposal
func (k Keeper) GetDeposits(ctx sdk.Context, proposalID uint64) (deposits v1.Deposits) {
	k.IterateDeposits(ctx, proposalID, func(deposit v1.Deposit) bool {
		deposits = append(deposits, &deposit)
		return false
	})

	return
}

// DeleteAndBurnDeposits deletes and burns all the deposits on a specific proposal.
func (k Keeper) DeleteAndBurnDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)

	k.IterateDeposits(ctx, proposalID, func(deposit v1.Deposit) bool {
		err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, deposit.Amount)
		if err != nil {
			panic(err)
		}

		depositor := sdk.MustAccAddressFromHex(deposit.Depositor)

		store.Delete(types.DepositKey(proposalID, depositor))
		return false
	})
}

// IterateAllDeposits iterates over all the stored deposits and performs a callback function.
func (k Keeper) IterateAllDeposits(ctx sdk.Context, cb func(deposit v1.Deposit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositsKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var deposit v1.Deposit

		k.cdc.MustUnmarshal(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// IterateDeposits iterates over all the proposals deposits and performs a callback function
func (k Keeper) IterateDeposits(ctx sdk.Context, proposalID uint64, cb func(deposit v1.Deposit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositsKey(proposalID))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var deposit v1.Deposit

		k.cdc.MustUnmarshal(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal.
// Activates voting period when appropriate and returns true in that case, else returns false.
func (k Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (bool, error) {
	// Checks to see if proposal exists
	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		return false, sdkerrors.Wrapf(types.ErrUnknownProposal, "%d", proposalID)
	}

	// Check if proposal is still depositable
	if (proposal.Status != v1.StatusDepositPeriod) && (proposal.Status != v1.StatusVotingPeriod) {
		return false, sdkerrors.Wrapf(types.ErrInactiveProposal, "%d", proposalID)
	}

	// update the governance module's account coins pool
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, types.ModuleName, depositAmount)
	if err != nil {
		return false, err
	}

	// Update proposal
	proposal.TotalDeposit = sdk.NewCoins(proposal.TotalDeposit...).Add(depositAmount...)
	k.SetProposal(ctx, proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false

	if proposal.Status == v1.StatusDepositPeriod && sdk.NewCoins(proposal.TotalDeposit...).IsAllGTE(k.GetParams(ctx).MinDeposit) {
		k.ActivateVotingPeriod(ctx, proposal)

		activatedVotingPeriod = true
	}

	// Add or update deposit object
	deposit, found := k.GetDeposit(ctx, proposalID, depositorAddr)

	if found {
		deposit.Amount = sdk.NewCoins(deposit.Amount...).Add(depositAmount...)
	} else {
		deposit = v1.NewDeposit(proposalID, depositorAddr, depositAmount)
	}

	// called when deposit has been added to a proposal, however the proposal may not be active
	k.Hooks().AfterProposalDeposit(ctx, proposalID, depositorAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	k.SetDeposit(ctx, deposit)

	return activatedVotingPeriod, nil
}

// RefundAndDeleteDeposits refunds and deletes all the deposits on a specific proposal.
func (k Keeper) RefundAndDeleteDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)

	k.IterateDeposits(ctx, proposalID, func(deposit v1.Deposit) bool {
		depositor := sdk.MustAccAddressFromHex(deposit.Depositor)

		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, deposit.Amount)
		if err != nil {
			panic(err)
		}

		store.Delete(types.DepositKey(proposalID, depositor))
		return false
	})
}

// validateInitialDeposit validates if initial deposit is greater than or equal to the minimum
// required at the time of proposal submission. This threshold amount is determined by
// the deposit parameters. Returns nil on success, error otherwise.
func (k Keeper) validateInitialDeposit(ctx sdk.Context, initialDeposit sdk.Coins) error {
	params := k.GetParams(ctx)
	minInitialDepositRatio, err := sdk.NewDecFromStr(params.MinInitialDepositRatio)
	if err != nil {
		return err
	}
	if minInitialDepositRatio.IsZero() {
		return nil
	}
	minDepositCoins := params.MinDeposit
	for i := range minDepositCoins {
		minDepositCoins[i].Amount = sdk.NewDecFromInt(minDepositCoins[i].Amount).Mul(minInitialDepositRatio).RoundInt()
	}
	if !initialDeposit.IsAllGTE(minDepositCoins) {
		return sdkerrors.Wrapf(types.ErrMinDepositTooSmall, "was (%s), need (%s)", initialDeposit, minDepositCoins)
	}
	return nil
}
