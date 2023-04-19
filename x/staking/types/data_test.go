package types_test

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	pk1      = ed25519.GenPrivKey().PubKey()
	pk1Any   *codectypes.Any
	pk2      = ed25519.GenPrivKey().PubKey()
	pk3      = ed25519.GenPrivKey().PubKey()
	valAddr1 = sdk.AccAddress(pk1.Address())
	valAddr2 = sdk.AccAddress(pk2.Address())
	valAddr3 = sdk.AccAddress(pk3.Address())

	emptyAddr   sdk.AccAddress
	emptyPubkey cryptotypes.PubKey
)

func init() {
	var err error
	pk1Any, err = codectypes.NewAnyWithValue(pk1)
	if err != nil {
		panic(fmt.Sprintf("Can't pack pk1 %t as Any", pk1))
	}
}
