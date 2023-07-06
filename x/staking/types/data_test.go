package types_test

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/eth/ethsecp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	pv1, _   = ethsecp256k1.GenPrivKey()
	pv2, _   = ethsecp256k1.GenPrivKey()
	pv3, _   = ethsecp256k1.GenPrivKey()
	pk1      = pv1.PubKey()
	pk1Any   *codectypes.Any
	pk2      = pv2.PubKey()
	pk3      = pv3.PubKey()
	addr1, _ = sdk.Bech32ifyAddressBytes(sdk.Bech32PrefixAccAddr, pk1.Address().Bytes())
	addr2, _ = sdk.Bech32ifyAddressBytes(sdk.Bech32PrefixAccAddr, pk2.Address().Bytes())
	addr3, _ = sdk.Bech32ifyAddressBytes(sdk.Bech32PrefixAccAddr, pk3.Address().Bytes())
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
