package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/x/params/client/utils"
)

func TestParseProposal(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Staking Param Change",
  "description": "Update max validators",
  "changes": [
    {
      "subspace": "staking",
      "key": "MaxValidators",
      "value": 1
    }
  ],
  "deposit": "1000stake"
}
`)
	proposal, err := utils.ParseParamChangeProposalJSON(cdc, okJSON.Name())
	require.NoError(t, err)

	require.Equal(t, "Staking Param Change", proposal.Title)
	require.Equal(t, "Update max validators", proposal.Description)
	require.Equal(t, "1000stake", proposal.Deposit)
	require.Equal(t, utils.ParamChangesJSON{
		{
			Subspace: "staking",
			Key:      "MaxValidators",
			Value:    []byte{0x31},
		},
	}, proposal.Changes)
}

func TestParseCrossChainProposal(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "BSC smart contract upgrade",
  "description": "BSC smart contract upgrade",
  "changes": [
    {
      "subspace": "BSC",
      "key": "upgrade",
      "value": "0x8f86403A4DE0BB5791fa46B8e795C547942fE4Cf"
    }
  ],
  "deposit": "1000000000000BNB",
  "cross_chain": true,
  "addresses": ["0x6c615C766EE6b7e69275b0D070eF50acc93ab880"]
}
`)
	proposal, err := utils.ParseParamChangeProposalJSON(cdc, okJSON.Name())
	require.NoError(t, err)

	require.Equal(t, "BSC smart contract upgrade", proposal.Title)
	require.Equal(t, "BSC smart contract upgrade", proposal.Description)
	require.Equal(t, "1000000000000BNB", proposal.Deposit)

	require.Equal(t, true, proposal.CrossChain)
	require.Equal(t, "0x6c615C766EE6b7e69275b0D070eF50acc93ab880", proposal.Addresses[0])
	require.Equal(t, utils.ParamChangesJSON{
		{
			Subspace: "BSC",
			Key:      "upgrade",
			Value:    []byte{0x22, 0x30, 0x78, 0x38, 0x66, 0x38, 0x36, 0x34, 0x30, 0x33, 0x41, 0x34, 0x44, 0x45, 0x30, 0x42, 0x42, 0x35, 0x37, 0x39, 0x31, 0x66, 0x61, 0x34, 0x36, 0x42, 0x38, 0x65, 0x37, 0x39, 0x35, 0x43, 0x35, 0x34, 0x37, 0x39, 0x34, 0x32, 0x66, 0x45, 0x34, 0x43, 0x66, 0x22},
		},
	}, proposal.Changes)
}
