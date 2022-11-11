package simapp_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

func TestSimGenesisAccountValidate(t *testing.T) {
	privKey1, _ := ethsecp256k1.GenerateKey()
	privKey2, _ := ethsecp256k1.GenerateKey()
	pubkey := privKey1.PubKey()
	addr := sdk.AccAddress(pubkey.Address())

	vestingStart := time.Now().UTC()

	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 1000))
	baseAcc := authtypes.NewBaseAccount(addr, pubkey, 0, 0)

	testCases := []struct {
		name    string
		sga     simapp.SimGenesisAccount
		wantErr bool
	}{
		{
			"valid basic account",
			simapp.SimGenesisAccount{
				BaseAccount: baseAcc,
			},
			false,
		},
		{
			"invalid basic account with mismatching address/pubkey",
			simapp.SimGenesisAccount{
				BaseAccount: authtypes.NewBaseAccount(addr, privKey2.PubKey(), 0, 0),
			},
			true,
		},
		{
			"valid basic account with module name",
			simapp.SimGenesisAccount{
				BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(crypto.AddressHash([]byte("testmod"))), nil, 0, 0),
				ModuleName:  "testmod",
			},
			false,
		},
		{
			"valid basic account with invalid module name/pubkey pair",
			simapp.SimGenesisAccount{
				BaseAccount: baseAcc,
				ModuleName:  "testmod",
			},
			true,
		},
		{
			"valid basic account with valid vesting attributes",
			simapp.SimGenesisAccount{
				BaseAccount:     baseAcc,
				OriginalVesting: coins,
				StartTime:       vestingStart.Unix(),
				EndTime:         vestingStart.Add(1 * time.Hour).Unix(),
			},
			false,
		},
		{
			"valid basic account with invalid vesting end time",
			simapp.SimGenesisAccount{
				BaseAccount:     baseAcc,
				OriginalVesting: coins,
				StartTime:       vestingStart.Add(2 * time.Hour).Unix(),
				EndTime:         vestingStart.Add(1 * time.Hour).Unix(),
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.wantErr, tc.sga.Validate() != nil)
		})
	}
}
