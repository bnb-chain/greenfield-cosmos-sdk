package proposal

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	SyncParamsChannel                 = "syncParametersChange"
	SyncParamsChannelID sdk.ChannelID = 3
	KeyUpgrade                        = "upgrade"
)

var KeySyncParamsRelayerFee = []byte("SyncParamsRelayerFee")

// SyncParamsPackage is the payload be relayed to BSC
type SyncParamsPackage struct {
	Key string //
	// new parameter or new smart contract address(es) if is ungraded proposal
	Value []byte // string   // address to bytes
	// smart contract address(es)
	Target []byte // string
}
