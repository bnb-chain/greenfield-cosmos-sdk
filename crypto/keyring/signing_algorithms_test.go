package keyring

import (
	"fmt"
	"testing"

	ethHd "github.com/evmos/ethermint/crypto/hd"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
)

func TestNewSigningAlgoByString(t *testing.T) {
	tests := []struct {
		name         string
		algoStr      string
		isSupported  bool
		expectedAlgo SignatureAlgo
		expectedErr  error
	}{
		{
			"supported algorithm",
			"eth_secp256k1",
			true,
			ethHd.EthSecp256k1,
			nil,
		},
		{
			"not supported",
			"notsupportedalgo",
			false,
			nil,
			fmt.Errorf("provided algorithm \"notsupportedalgo\" is not supported"),
		},
	}

	list := SigningAlgoList{ethHd.EthSecp256k1}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			algorithm, err := NewSigningAlgoFromString(tt.algoStr, list)
			if tt.isSupported {
				require.Equal(t, ethHd.EthSecp256k1, algorithm)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func TestAltSigningAlgoList_Contains(t *testing.T) {
	list := SigningAlgoList{ethHd.EthSecp256k1}

	require.True(t, list.Contains(ethHd.EthSecp256k1))
	require.False(t, list.Contains(notSupportedAlgo{}))
}

func TestAltSigningAlgoList_String(t *testing.T) {
	list := SigningAlgoList{ethHd.EthSecp256k1, notSupportedAlgo{}}
	require.Equal(t, fmt.Sprintf("%s,notSupported", hd.Secp256k1Type), list.String())
}

type notSupportedAlgo struct{}

func (n notSupportedAlgo) Name() hd.PubKeyType {
	return "notSupported"
}

func (n notSupportedAlgo) Derive() hd.DeriveFn {
	return ethHd.EthSecp256k1.Derive()
}

func (n notSupportedAlgo) Generate() hd.GenerateFn {
	return ethHd.EthSecp256k1.Generate()
}
