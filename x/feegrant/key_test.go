package feegrant_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

func TestMarshalAndUnmarshalFeegrantKey(t *testing.T) {
	grantee, err := sdk.AccAddressFromHexUnsafe("0xdF4AA991e455907136662745D576449949110290")
	require.NoError(t, err)
	granter, err := sdk.AccAddressFromHexUnsafe("0xf5Acddf298F45733426E2b4D362413a54ba024Fa")
	require.NoError(t, err)

	key := feegrant.FeeAllowanceKey(granter, grantee)
	require.Len(t, key, len(grantee.Bytes())+len(granter.Bytes())+3)
	require.Equal(t, feegrant.FeeAllowancePrefixByGrantee(grantee), key[:len(grantee.Bytes())+2])

	g1, g2 := feegrant.ParseAddressesFromFeeAllowanceKey(key)
	require.Equal(t, granter, g1)
	require.Equal(t, grantee, g2)
}

func TestMarshalAndUnmarshalFeegrantKeyQueueKey(t *testing.T) {
	grantee, err := sdk.AccAddressFromHexUnsafe("0xdF4AA991e455907136662745D576449949110290")
	require.NoError(t, err)
	granter, err := sdk.AccAddressFromHexUnsafe("0xf5Acddf298F45733426E2b4D362413a54ba024Fa")
	require.NoError(t, err)

	exp := time.Now()
	expBytes := sdk.FormatTimeBytes(exp)

	key := feegrant.FeeAllowancePrefixQueue(&exp, feegrant.FeeAllowanceKey(granter, grantee)[1:])
	require.Len(t, key, len(grantee.Bytes())+len(granter.Bytes())+3+len(expBytes))

	granter1, grantee1 := feegrant.ParseAddressesFromFeeAllowanceQueueKey(key)
	require.Equal(t, granter, granter1)
	require.Equal(t, grantee, grantee1)
}
