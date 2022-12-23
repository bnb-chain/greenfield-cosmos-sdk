package ante

import (
	"fmt"

	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/feehub/types"
)

// ConsumeMsgGasDecorator will take in parameters and consume gas depending on
// the size of tx and msg type before calling next AnteHandler.
//
// CONTRACT: If simulate=true, then signatures must either be completely filled
// in or empty.
type ConsumeMsgGasDecorator struct {
	ak  AccountKeeper
	fhk FeehubKeeper
}

func NewConsumeMsgGasDecorator(ak AccountKeeper, fhk FeehubKeeper) ConsumeMsgGasDecorator {
	return ConsumeMsgGasDecorator{
		ak:  ak,
		fhk: fhk,
	}
}

func (cmfg ConsumeMsgGasDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, errors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}

	params := cmfg.fhk.GetParams(ctx)
	txBytesLength := uint64(len(ctx.TxBytes()))

	msgGas, err := cmfg.getMsgGas(params, sigTx)
	if err != nil {
		return ctx, err
	}

	// simulate gas cost for signatures in simulate mode
	if simulate {
		// in simulate mode, each element should be a nil signature
		sigs, err := sigTx.GetSignaturesV2()
		if err != nil {
			return ctx, err
		}
		n := len(sigs)

		for i, signer := range sigTx.GetSigners() {
			// if signature is already filled in, no need to simulate gas cost
			if i < n && !isIncompleteSignature(sigs[i].Data) {
				continue
			}

			var pubkey cryptotypes.PubKey

			acc := cmfg.ak.GetAccount(ctx, signer)

			// use placeholder simSecp256k1Pubkey if sig is nil
			if acc == nil || acc.GetPubKey() == nil {
				pubkey = simSecp256k1Pubkey
			} else {
				pubkey = acc.GetPubKey()
			}

			// use stdsignature to mock the size of a full signature
			simSig := legacytx.StdSignature{ //nolint:staticcheck // this will be removed when proto is ready
				Signature: simSecp256k1Sig[:],
				PubKey:    pubkey,
			}

			sigBz := legacy.Cdc.MustMarshal(simSig)
			txBytesLength = txBytesLength + uint64(len(sigBz)) + 6
		}
	}

	if txBytesLength > params.GetMaxTxSize() {
		return ctx, errors.Wrapf(sdkerrors.ErrTxTooLarge, "tx length: %d, limit: %d", txBytesLength, params.GetMaxTxSize())
	}

	txBytesGas := txBytesLength * params.GetMinGasPerByte()
	if txBytesGas >= msgGas {
		ctx.GasMeter().ConsumeGas(txBytesGas, "gas by tx byte length")
	} else {
		ctx.GasMeter().ConsumeGas(msgGas, "msg gas")
	}

	return next(ctx, tx, simulate)
}

func (cmfg ConsumeMsgGasDecorator) getMsgGas(params types.Params, tx sdk.Tx) (uint64, error) {
	msg := tx.GetMsgs()[0]
	feeCalcGen := types.GetCalculatorGen(sdk.MsgTypeURL(msg))
	if feeCalcGen == nil {
		return 0, fmt.Errorf("failed to find fee calculator")
	}
	feeCalc := feeCalcGen(params)
	return feeCalc(msg), nil
}
