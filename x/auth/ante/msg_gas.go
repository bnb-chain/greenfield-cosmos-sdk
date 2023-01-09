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
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

// ValidateTxSizeDecorator will validate tx bytes length given the parameters passed in
// If tx is too large decorator returns with error, otherwise call next AnteHandler
//
// CONTRACT: If simulate=true, then signatures must either be completely filled
// in or empty.
type ValidateTxSizeDecorator struct {
	ak  AccountKeeper
	fhk GashubKeeper
}

func NewValidateTxSizeDecorator(ak AccountKeeper, fhk GashubKeeper) ValidateTxSizeDecorator {
	return ValidateTxSizeDecorator{
		ak:  ak,
		fhk: fhk,
	}
}

func (vtsd ValidateTxSizeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, errors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}

	newCtx := ctx
	txSize := newCtx.TxSize()
	// simulate signatures in simulate mode
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

			acc := vtsd.ak.GetAccount(ctx, signer)

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
			txSize = txSize + uint64(len(sigBz)) + 14
		}

		newCtx = ctx.WithTxSize(txSize)
	}

	params := vtsd.fhk.GetParams(ctx)
	if txSize > params.GetMaxTxSize() {
		return ctx, errors.Wrapf(sdkerrors.ErrTxTooLarge, "tx length: %d, limit: %d", txSize, params.GetMaxTxSize())
	}

	return next(newCtx, tx, simulate)
}

// ConsumeMsgGasDecorator will take in parameters and consume gas depending on
// the size of tx and msg type before calling next AnteHandler.
type ConsumeMsgGasDecorator struct {
	ak  AccountKeeper
	fhk GashubKeeper
}

func NewConsumeMsgGasDecorator(ak AccountKeeper, fhk GashubKeeper) ConsumeMsgGasDecorator {
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
	gasByTxSize := cmfg.getTxSizeGas(params, ctx)
	gasByMsgType, err := cmfg.getMsgGas(params, sigTx)
	if err != nil {
		return ctx, err
	}

	if gasByTxSize > gasByMsgType {
		ctx.GasMeter().ConsumeGas(gasByTxSize, "tx bytes length")
	} else {
		ctx.GasMeter().ConsumeGas(gasByMsgType, "msg type")
	}

	return next(ctx, tx, simulate)
}

func (cmfg ConsumeMsgGasDecorator) getMsgGas(params types.Params, tx sdk.Tx) (uint64, error) {
	msgs := tx.GetMsgs()
	totalGas := uint64(0)
	for _, msg := range msgs {
		feeCalcGen := types.GetGasCalculatorGen(sdk.MsgTypeURL(msg))
		if feeCalcGen == nil {
			return 0, fmt.Errorf("failed to find fee calculator")
		}
		feeCalc := feeCalcGen(params)
		gas, err := feeCalc(msg)
		if err != nil {
			return 0, err
		}
		totalGas += gas
	}
	return totalGas, nil
}

func (cmfg ConsumeMsgGasDecorator) getTxSizeGas(params types.Params, ctx sdk.Context) uint64 {
	return params.GetMinGasPerByte() * ctx.TxSize()
}
