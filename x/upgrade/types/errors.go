package types

import (
	"cosmossdk.io/errors"
)

// x/upgrade module sentinel errors
var (
	// ErrUpgradeScheduled error if the upgrade scheduled in the past
	ErrUpgradeScheduled = errors.Register(ModuleName, 2, "upgrade cannot be scheduled in the past")
	// ErrUpgradeCompleted error if the upgrade has already been completed
	ErrUpgradeCompleted = errors.Register(ModuleName, 3, "upgrade has already been completed")
)
