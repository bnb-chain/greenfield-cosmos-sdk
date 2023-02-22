package keeper

import (
	"encoding/hex"
	"testing"

	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

func TestSyncParams(t *testing.T) {
	var paramChanges []types.ParamChange
	ug1 := types.ParamChange{
		Subspace: "BSC",
		Key:      "upgrade",
		Value:    "0x8f86403A4DE0BB5791fa46B8e795C547942fE4Cf",
	}
	paramChanges = append(paramChanges, ug1)
	ug2 := types.ParamChange{
		Subspace: "BSC",
		Key:      "upgrade",
		Value:    "0x9d4454B023096f34B160D6B654540c56A1F81688",
	}
	paramChanges = append(paramChanges, ug2)
	content := types.NewCrossChainParameterChangeProposal(
		"upgrade GovHub and CrossChain",
		"upgrade GovHub and CrossChain",
		paramChanges,
		[]string{"0x6c615C766EE6b7e69275b0D070eF50acc93ab880", "0x04ED4ad3cDe36FE8ba944E3D6CFC54f7Fe6c3C72"},
	)
	values := make([]byte, 0)
	addresses := make([]byte, 0)

	for i, c := range paramChanges {
		// params value key = "abc" value = "value"
		// 1. address hex.DecodeString()
		//decodeString, err := hex.DecodeString(c.Value)
		//if err != nil {
		//	return
		//}
		// parameter
		hex.DecodeString(c.Value)

		// if value address
		value, _ := sdk.AccAddressFromHexUnsafe(c.Value)
		values = append(values, value.Bytes()...)

		addr, _ := sdk.AccAddressFromHexUnsafe(content.Addresses[i])
		addresses = append(addresses, addr.Bytes()...)
	}

	pack := types.SyncParamsPackage{
		Key:    content.Changes[0].Key,
		Value:  values,
		Target: addresses,
	}
	encodedPackage, _ := rlp.EncodeToBytes(pack)
	t.Log(hex.EncodeToString(encodedPackage))
}

// f8878775706772616465b854307838663836343033413444453042423537393166613436423865373935433534373934326645344366307839643434353442303233303936663334423136304436423635343534306335364131463831363838a86c615c766ee6b7e69275b0d070ef50acc93ab88004ed4ad3cde36fe8ba944e3d6cfc54f7fe6c3c72
