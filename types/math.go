package types

import (
	"math"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Type aliases to the SDK's math submodule
//
// Deprecated: Functionality of this package has been moved to its own module:
// cosmossdk.io/math
//
// Please use the above module instead of this package.
type (
	Int  = sdkmath.Int
	Uint = sdkmath.Uint
)

var (
	NewIntFromBigInt = sdkmath.NewIntFromBigInt
	OneInt           = sdkmath.OneInt
	NewInt           = sdkmath.NewInt
	ZeroInt          = sdkmath.ZeroInt
	IntEq            = sdkmath.IntEq
	NewIntFromString = sdkmath.NewIntFromString
	NewUint          = sdkmath.NewUint
	NewIntFromUint64 = sdkmath.NewIntFromUint64
	MaxInt           = sdkmath.MaxInt
	MinInt           = sdkmath.MinInt
)

const (
	MaxBitLen = sdkmath.MaxBitLen
)

func (ip IntProto) String() string {
	return ip.Int.String()
}

// SafeInt64 checks for overflows while casting an uint64 to int64 value.
func SafeInt64(value uint64) (int64, error) {
	if value > uint64(math.MaxInt64) {
		return 0, errors.Wrapf(sdkerrors.ErrInvalidHeight, "uint64 value %v cannot exceed %v", value, int64(math.MaxInt64))
	}

	return int64(value), nil
}
