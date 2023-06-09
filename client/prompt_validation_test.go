package client_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/require"
)

func TestValidatePromptNotEmpty(t *testing.T) {
	require := require.New(t)

	require.NoError(client.ValidatePromptNotEmpty("foo"))
	require.ErrorContains(client.ValidatePromptNotEmpty(""), "input cannot be empty")
}

func TestValidatePromptURL(t *testing.T) {
	require := require.New(t)

	require.NoError(client.ValidatePromptURL("https://example.com"))
	require.ErrorContains(client.ValidatePromptURL("foo"), "invalid URL")
}

func TestValidatePromptAddress(t *testing.T) {
	require := require.New(t)

	require.NoError(client.ValidatePromptAddress("0x319D057ce294319bA1fa5487134608727e1F3e29"))
	require.NoError(client.ValidatePromptAddress("0x319D057ce294319bA1fa5487134608727e1F3e29"))
	require.NoError(client.ValidatePromptAddress("0x319D057ce294319bA1fa5487134608727e1F3e29"))
	require.ErrorContains(client.ValidatePromptAddress("foo"), "invalid address")
}

func TestValidatePromptCoins(t *testing.T) {
	require := require.New(t)

	require.NoError(client.ValidatePromptCoins("100stake"))
	require.ErrorContains(client.ValidatePromptCoins("foo"), "invalid coins")
}
