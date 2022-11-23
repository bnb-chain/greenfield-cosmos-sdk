package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName   = "oracle"
	QuerierRoute = ModuleName

	// RelayPackagesChannelId is not a communication channel actually, we just use it to record sequence.
	RelayPackagesChannelName               = "relayPackages"
	RelayPackagesChannelId   sdk.ChannelID = 0x00
)
