package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/evidence/types"
)

func TestEquivocation_Valid(t *testing.T) {
	n, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	addr := sdk.ConsAddress("foo_________________")

	e := types.Equivocation{
		Height:           100,
		Time:             n,
		Power:            1000000,
		ConsensusAddress: addr.String(),
	}

	require.Equal(t, e.GetTotalPower(), int64(0))
	require.Equal(t, e.GetValidatorPower(), e.Power)
	require.Equal(t, e.GetTime(), e.Time)
	require.Equal(t, e.GetConsensusAddress().String(), e.ConsensusAddress)
	require.Equal(t, e.GetHeight(), e.Height)
	require.Equal(t, e.Type(), types.TypeEquivocation)
	require.Equal(t, e.Route(), types.RouteEquivocation)
	require.Equal(t, e.Hash().String(), "93707E4C05DB40E8F061301C3911BE4933B5E40AC043D5BD86444C6EBDA964EB")
	require.Equal(t, e.String(), "consensus_address: 0x666F6F5F5F5F5f5f5f5f5f5F5f5f5f5F5f5F5f5f\nheight: 100\npower: 1000000\ntime: \"2006-01-02T15:04:05Z\"\n")
	require.NoError(t, e.ValidateBasic())

	require.Equal(t, int64(0), e.GetTotalPower())
	require.Equal(t, e.Power, e.GetValidatorPower())
	require.Equal(t, e.Time, e.GetTime())
	require.Equal(t, e.ConsensusAddress, e.GetConsensusAddress().String())
	require.Equal(t, e.Height, e.GetHeight())
	require.Equal(t, types.TypeEquivocation, e.Type())
	require.Equal(t, types.RouteEquivocation, e.Route())
	require.Equal(t, "93707E4C05DB40E8F061301C3911BE4933B5E40AC043D5BD86444C6EBDA964EB", e.Hash().String())
	require.Equal(t, "consensus_address: 0x666F6F5F5F5F5f5f5f5f5f5F5f5f5f5F5f5F5f5f\nheight: 100\npower: 1000000\ntime: \"2006-01-02T15:04:05Z\"\n", e.String())
	require.NoError(t, e.ValidateBasic())
}

func TestEquivocationValidateBasic(t *testing.T) {
	var zeroTime time.Time
	addr := sdk.ConsAddress("foo_________________")

	n, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	testCases := []struct {
		name      string
		e         types.Equivocation
		expectErr bool
	}{
		{"valid", types.Equivocation{100, n, 1000000, addr.String()}, false},
		{"invalid time", types.Equivocation{100, zeroTime, 1000000, addr.String()}, true},
		{"invalid height", types.Equivocation{0, n, 1000000, addr.String()}, true},
		{"invalid power", types.Equivocation{100, n, 0, addr.String()}, true},
		{"invalid address", types.Equivocation{100, n, 1000000, ""}, true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectErr, tc.e.ValidateBasic() != nil)
		})
	}
}

func TestEvidenceAddressConversion(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForConsensusNode("testcnclcons", "testcnclconspub")
	tmEvidence := abci.Evidence{
		Type: abci.EvidenceType_DUPLICATE_VOTE,
		Validator: abci.Validator{
			Address: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Power:   100,
		},
		Height:           1,
		Time:             time.Now(),
		TotalVotingPower: 100,
	}
	evidence := types.FromABCIEvidence(tmEvidence).(*types.Equivocation)
	consAddr := evidence.GetConsensusAddress()
	// Check the address is the same after conversion
	require.Equal(t, tmEvidence.Validator.Address, consAddr.Bytes())
	sdk.GetConfig().SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
}
