package codec

import (
	"testing"

	"github.com/stretchr/testify/require"

	gogotypes "github.com/cosmos/gogoproto/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/collections/colltest"
)

func TestCollectionsCorrectness(t *testing.T) {
	cdc := NewProtoCodec(codectypes.NewInterfaceRegistry())
	t.Run("CollValue", func(t *testing.T) {
		colltest.TestValueCodec(t, CollValue[gogotypes.UInt64Value](cdc), gogotypes.UInt64Value{
			Value: 500,
		})
	})

	t.Run("BoolValue", func(t *testing.T) {
		colltest.TestValueCodec(t, BoolValue, true)
		colltest.TestValueCodec(t, BoolValue, false)

		// asserts produced bytes are equal
		valueAssert := func(b bool) {
			wantBytes, err := (&gogotypes.BoolValue{Value: b}).Marshal()
			require.NoError(t, err)
			gotBytes, err := BoolValue.Encode(b)
			require.NoError(t, err)
			require.Equal(t, wantBytes, gotBytes)
		}

		valueAssert(true)
		valueAssert(false)
	})
}
