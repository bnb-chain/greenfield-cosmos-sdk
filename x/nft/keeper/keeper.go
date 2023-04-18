package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/nft"

	store "github.com/cosmos/cosmos-sdk/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Keeper of the nft store
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	bk           nft.BankKeeper
}

// NewKeeper creates a new nft Keeper instance
func NewKeeper(storeService store.KVStoreService,
	cdc codec.BinaryCodec, ak nft.AccountKeeper, bk nft.BankKeeper,
) Keeper {
	// ensure nft module account is set
	if addr := ak.GetModuleAddress(nft.ModuleName); addr == nil {
		panic("the nft module account has not been set")
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		bk:           bk,
	}
}
