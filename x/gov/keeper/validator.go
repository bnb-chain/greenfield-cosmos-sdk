package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// CheckCreateProposal checks whether create validator proposal has passed and
// match the current create validator operation.
func (k Keeper) CheckCreateProposal(ctx sdk.Context, msg *types.MsgCreateValidator) error {
	_, ok := k.GetProposal(ctx, msg.ProposalId)
	if !ok {
		return fmt.Errorf("proposal %d does not exist", msg.ProposalId)
	}

	// TODO: check the proposal according to the proposal contents.
	return nil
}

// CheckRemoveProposal checks whether remove validator proposal has passed and
// match the current remove validator operation.
func (k Keeper) CheckRemoveProposal(ctx sdk.Context, msg *types.MsgRemoveValidator) error {
	_, ok := k.GetProposal(ctx, msg.ProposalId)
	if !ok {
		return fmt.Errorf("proposal %d does not exist", msg.ProposalId)
	}

	// TODO: check the proposal according to the proposal contents.
	return nil
}
