package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var ParamsKey = []byte{0x01}

const (
	ModuleName   = "oracle"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	// RelayPackagesChannelId is not a communication channel actually, we just use it to record sequence.
	RelayPackagesChannelName               = "relayPackages"
	RelayPackagesChannelId   sdk.ChannelID = 0x00
	MultiMessageChannelId    sdk.ChannelID = 0x08
)
