package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlsClaim(t *testing.T) {
	claim := &BlsClaim{
		ChainId:  1,
		Sequence: 1,
		Payload:  []byte("test payload"),
	}

	signBytes := claim.GetSignBytes()

	require.Equal(t, "3b0858e23a9ca1335fff8539c8a27037ed29a4e5c2258a92c590ca9ad319abe0",
		hex.EncodeToString(signBytes[:]))
}
