package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	SyncParamsChannel = "syncParametersChange"

	SyncParamsChannelID sdk.ChannelID = 3

	KeyUpgrade = "upgrade"
)

// SyncParamsPackage is the payload to be encoded for cross-chain IBC package
type SyncParamsPackage struct {
	// Key is the parameter to be changed
	Key string
	// Value is either new parameter or new smart contract address(es) if it is an upgrade proposal
	Value []byte
	// Target is the smart contract address(es)
	Target []byte
}

var (
	syncParamsPackageType, _ = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "Key", Type: "string"},
		{Name: "Value", Type: "bytes"},
		{Name: "Target", Type: "bytes"},
	})

	syncParamsPackageArgs = abi.Arguments{
		{Type: syncParamsPackageType},
	}
)

func (p SyncParamsPackage) Serialize() ([]byte, error) {
	encodedBytes, err := syncParamsPackageArgs.Pack(&p)
	if err != nil {
		return nil, err
	}
	return encodedBytes, nil
}
