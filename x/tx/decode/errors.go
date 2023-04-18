package decode

import "github.com/cosmos/cosmos-sdk/errors"

const (
	txCodespace = "tx"
)

var ErrUnknownField = errors.Register(txCodespace, 2, "unknown protobuf field")
