package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/collections/colltest"
)

func TestIntValue(t *testing.T) {
	colltest.TestValueCodec(t, IntValue, NewInt(10005994859))
}
