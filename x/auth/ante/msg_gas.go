package ante

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

const (
	EthSecp256k1SigSize = 65
	FeeSize             = 42
)

// ValidateTxSizeDecorator will validate tx bytes length given the parameters passed in
// If tx is too large decorator returns with error, otherwise call next AnteHandler
//
// CONTRACT: If simulate=true, then signatures must either be completely filled
// in or empty.
type ValidateTxSizeDecorator struct {
	ak  AccountKeeper
	ghk GashubKeeper
}

func NewValidateTxSizeDecorator(ak AccountKeeper, ghk GashubKeeper) ValidateTxSizeDecorator {
	return ValidateTxSizeDecorator{
		ak:  ak,
		ghk: ghk,
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

		txSize += FeeSize
		for i := range sigTx.GetSigners() {
			// if signature is already filled in, no need to simulate gas cost
			if i < n && !isIncompleteSignature(sigs[i].Data) {
				continue
			}

			txSize += EthSecp256k1SigSize
		}

		newCtx = ctx.WithTxSize(txSize)
	}

	params := vtsd.ghk.GetParams(ctx)
	if txSize > params.GetMaxTxSize() {
		return ctx, errors.Wrapf(sdkerrors.ErrTxTooLarge, "tx length: %d, limit: %d", txSize, params.GetMaxTxSize())
	}

	return next(newCtx, tx, simulate)
}

// ConsumeMsgGasDecorator will take in parameters and consume gas depending on
// the size of tx and msg type before calling next AnteHandler.
type ConsumeMsgGasDecorator struct {
	ak  AccountKeeper
	ghk GashubKeeper
}

func NewConsumeMsgGasDecorator(ak AccountKeeper, ghk GashubKeeper) ConsumeMsgGasDecorator {
	return ConsumeMsgGasDecorator{
		ak:  ak,
		ghk: ghk,
	}
}

func (cmfg ConsumeMsgGasDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, errors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}

	params := cmfg.ghk.GetParams(ctx)
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
