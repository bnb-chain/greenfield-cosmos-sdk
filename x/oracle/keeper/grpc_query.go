package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
)

var _ types.QueryServer = Keeper{}
