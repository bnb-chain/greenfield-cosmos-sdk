package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseCLIProposal(t *testing.T) {
	data := []byte(`{
			"group_policy_address": "0xa0d45a1aa52d728662cbcf76585ba2216cc4b041",
			"messages": [
			  {
				"@type": "/cosmos.bank.v1beta1.MsgSend",
				"from_address": "0xa0d45a1aa52d728662cbcf76585ba2216cc4b041",
				"to_address": "0xa0d45a1aa52d728662cbcf76585ba2216cc4b041",
				"amount":[{"denom": "stake","amount": "10"}]
			  }
			],
			"metadata": "4pIMOgIGx1vZGU=",
			"proposers": ["0xa0d45a1aa52d728662cbcf76585ba2216cc4b041"]
		}`)

	result, err := parseCLIProposal(data)
	require.NoError(t, err)
	require.Equal(t, result.GroupPolicyAddress, "0xa0d45a1aa52d728662cbcf76585ba2216cc4b041")
	require.NotEmpty(t, result.Metadata)
	require.Equal(t, result.Metadata, "4pIMOgIGx1vZGU=")
	require.Equal(t, result.Proposers, []string{"0xa0d45a1aa52d728662cbcf76585ba2216cc4b041"})
}
