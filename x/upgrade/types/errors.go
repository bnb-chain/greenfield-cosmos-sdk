package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/upgrade module sentinel errors
var (
	// ErrUpgradeScheduled error if the upgrade scheduled in the past
	ErrUpgradeScheduled = sdkerrors.Register(ModuleName, 2, "upgrade cannot be scheduled in the past")
	// ErrUpgradeCompleted error if the upgrade has already been completed
	ErrUpgradeCompleted = sdkerrors.Register(ModuleName, 3, "upgrade has already been completed")
)
