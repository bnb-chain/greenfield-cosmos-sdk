package types

import "fmt"

const (
	// ModuleName is the name of this module
	ModuleName = "upgrade"

	// StoreKey is the prefix under which we store this module's data
	StoreKey = ModuleName
)

const (
	// DoneByte is a prefix to look up completed upgrade plan by name
	DoneByte = 0x1

	// VersionMapByte is a prefix to look up module names (key) and versions (value)
	VersionMapByte = 0x2

	// KeyUpgradedIBCState is the key under which upgraded ibc state is stored in the upgrade store
	KeyUpgradedIBCState = "upgradedIBCState"

	// KeyUpgradedClient is the sub-key under which upgraded client state will be stored
	KeyUpgradedClient = "upgradedClient"

	// KeyUpgradedConsState is the sub-key under which upgraded consensus state will be stored
	// it's unused, but need to compatible to ibc-go package
	KeyUpgradedConsState = "upgradedConsState"
)

// UpgradedConsStateKey is the key under which the upgraded consensus state is saved
// Connecting IBC chains can verify against the upgraded consensus state in this path before
// upgrading their clients.
func UpgradedConsStateKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s/%d/%s", KeyUpgradedIBCState, height, KeyUpgradedConsState))
}
