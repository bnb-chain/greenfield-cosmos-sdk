package signing

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/crypto/keys/eth/ethsecp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

// VerifySignature verifies a transaction signature contained in SignatureData abstracting over different signing modes
// and single vs multi-signatures.
func VerifySignature(ctx sdk.Context, pubKey cryptotypes.PubKey, signerData SignerData, sigData signing.SignatureData, handler SignModeHandler, tx sdk.Tx) error {
	switch data := sigData.(type) {
	case *signing.SingleSignatureData:
		// EIP712 signatures are verified in a different way
		// In greenfield, we adapt another antehandler to reject non-EIP712 signatures
		if data.SignMode == signing.SignMode_SIGN_MODE_EIP_712 {
			// check signature length
			if len(data.Signature) != ethcrypto.SignatureLength {
				return errorsmod.Wrap(sdkerrors.ErrorInvalidSigner, "signature length doesn't match typical [R||S||V] signature 65 bytes")
			}

			// skip signature verification if we have a cache and the tx is already in it
			if ctx.SigCache() != nil && ctx.TxBytes() != nil {
				if _, known := ctx.SigCache().Get(string(ctx.TxBytes())); known {
					return nil
				}
			}

			// verify signature
			err := verifyEip712SignatureWithFallback(ctx, pubKey, data.Signature, handler, signerData, tx)
			if err == nil && ctx.SigCache() != nil && ctx.TxBytes() != nil {
				ctx.SigCache().Add(string(ctx.TxBytes()), tx)
			}
			return err
		} else {
			// original cosmos-sdk signature verification
			signBytes, err := handler.GetSignBytes(data.SignMode, signerData, tx)
			if err != nil {
				return err
			}
			if !pubKey.VerifySignature(signBytes, data.Signature) {
				return fmt.Errorf("unable to verify single signer signature")
			}
			return nil
		}

	case *signing.MultiSignatureData:
		return fmt.Errorf("multi signature is not allowed")
	default:
		return fmt.Errorf("unexpected SignatureData %T", sigData)
	}
}

func verifyEip712SignatureWithFallback(ctx sdk.Context, pubKey cryptotypes.PubKey, sig []byte, handler SignModeHandler, signerData SignerData, tx sdk.Tx) error {
	// try with the old sign scheme first (for backward compatibility)
	sigHash, err := handler.GetSignBytes(signing.SignMode_SIGN_MODE_EIP_712, signerData, tx)
	if err == nil {
		if err := verifyEip712Signature(pubKey, sig, sigHash); err == nil {
			return nil
		}
	}

	// try with the new sign scheme
	sigHash, err = handler.GetSignBytesRuntime(ctx, signing.SignMode_SIGN_MODE_EIP_712, signerData, tx)
	if err != nil {
		return err
	}
	return verifyEip712Signature(pubKey, sig, sigHash)
}

func verifyEip712Signature(pubKey cryptotypes.PubKey, sig []byte, msg []byte) error {
	// remove the recovery offset if needed (ie. Metamask eip712 signature)
	if sig[ethcrypto.RecoveryIDOffset] == 27 || sig[ethcrypto.RecoveryIDOffset] == 28 {
		sig[ethcrypto.RecoveryIDOffset] -= 27
	}

	// recover the pubkey from the signature
	feePayerPubkey, err := secp256k1.RecoverPubkey(msg, sig)
	if err != nil {
		return errorsmod.Wrap(err, "failed to recover fee payer from sig")
	}

	ecPubKey, err := ethcrypto.UnmarshalPubkey(feePayerPubkey)
	if err != nil {
		return errorsmod.Wrap(err, "failed to unmarshal recovered fee payer pubkey")
	}

	// check that the recovered pubkey matches the one in the signerData data
	pk := &ethsecp256k1.PubKey{
		Key: ethcrypto.CompressPubkey(ecPubKey),
	}
	if !pubKey.Equals(pk) {
		return errorsmod.Wrapf(sdkerrors.ErrorInvalidSigner, "feePayer's pubkey %s is different from signature's pubkey %s", pubKey, pk)
	}

	return nil
}
