package std

import (
	"github.com/cosmos/cosmos-sdk/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/tx/signing/direct"
	"github.com/cosmos/cosmos-sdk/x/tx/signing/textual"
)

// SignModeOptions are options for configuring the standard sign mode handler map.
type SignModeOptions struct {
	// Textual are options for SIGN_MODE_TEXTUAL
	Textual textual.SignModeOptions
}

// HandlerMap returns a sign mode handler map that Cosmos SDK apps can use out
// of the box to support all "standard" sign modes.
func (s SignModeOptions) HandlerMap() (*signing.HandlerMap, error) {
	txt, err := textual.NewSignModeHandler(s.Textual)
	if err != nil {
		return nil, err
	}

	return signing.NewHandlerMap(
		direct.SignModeHandler{},
		txt,
	), nil
}
