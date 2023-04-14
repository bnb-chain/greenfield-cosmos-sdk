package v3_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func compareProps(t *testing.T, oldProp v1beta1.Proposal, newProp v1.Proposal) {
	require.Equal(t, oldProp.ProposalId, newProp.Id)
	require.Equal(t, oldProp.TotalDeposit.String(), sdk.Coins(newProp.TotalDeposit).String())
	require.Equal(t, oldProp.Status.String(), newProp.Status.String())
	require.Equal(t, oldProp.FinalTallyResult.Yes.String(), newProp.FinalTallyResult.YesCount)
	require.Equal(t, oldProp.FinalTallyResult.No.String(), newProp.FinalTallyResult.NoCount)
	require.Equal(t, oldProp.FinalTallyResult.NoWithVeto.String(), newProp.FinalTallyResult.NoWithVetoCount)
	require.Equal(t, oldProp.FinalTallyResult.Abstain.String(), newProp.FinalTallyResult.AbstainCount)

	newContent := newProp.Messages[0].GetCachedValue().(*v1.MsgExecLegacyContent).Content.GetCachedValue().(v1beta1.Content)
	require.Equal(t, oldProp.Content.GetCachedValue().(v1beta1.Content), newContent)

	// Compare UNIX times, as a simple Equal gives difference between Local and
	// UTC times.
	// ref: https://github.com/golang/go/issues/19486#issuecomment-292968278
	require.Equal(t, oldProp.SubmitTime.Unix(), newProp.SubmitTime.Unix())
	require.Equal(t, oldProp.DepositEndTime.Unix(), newProp.DepositEndTime.Unix())
	require.Equal(t, oldProp.VotingStartTime.Unix(), newProp.VotingStartTime.Unix())
	require.Equal(t, oldProp.VotingEndTime.Unix(), newProp.VotingEndTime.Unix())
}
