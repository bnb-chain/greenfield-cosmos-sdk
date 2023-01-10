package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeccak256Hash(t *testing.T) {
	input := []byte("test input")
	hash := Keccak256Hash(input)
	require.Equal(t, "1df1102036c102fbc689e6f72a64a9162ae0b1ea151932530deb8cd186c36c01", hex.EncodeToString(hash[:]))
}
