package errors

import "cosmossdk.io/errors"

// gashubCodespace is the codespace for all errors defined in gashub package
const gashubCodespace = "gashub"

var ErrInvalidMsgGasParams = errors.Register(gashubCodespace, 2, "msg gas params are invalid")
