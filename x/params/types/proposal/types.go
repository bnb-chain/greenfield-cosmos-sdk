package proposal

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	SyncParamsChangeChannel                  = "syncParametersChange"
	SyncParamsChangeChannellID sdk.ChannelID = 3
	BridgeSubspace                           = "bridge"
)

var (
	KeySyncParamsChangeRelayerFee = []byte("SyncParamsChangeRelayerFee")
)

// rlp (SyncParamsChangePackage)
type SyncParamsChangePackage struct {
	Key string //
	// new parameter or new smart contract address(es) if is ungraded proposal
	Values []byte // string   // address to bytes
	// smart contract address(es)
	Targets []byte // string
}
