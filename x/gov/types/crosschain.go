package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	SyncParamsChannel                 = "syncParametersChange"
	SyncParamsChannelID sdk.ChannelID = 3
	KeyUpgrade                        = "upgrade"
)

// SyncParamsPackage is the payload to be encoded for cross-chain IBC package
type SyncParamsPackage struct {
	// Key is the parameter to be changed
	Key string
	// Value is either new parameter or new smart contract address(es) if it is an upgrade proposal
	Value []byte
	// Target is the smart contract address(es)
	Target []byte
}