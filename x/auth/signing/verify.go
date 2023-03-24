package signing

import (
	"fmt"

	"cosmossdk.io/errors"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

// VerifySignature verifies a transaction signature contained in SignatureData abstracting over different signing modes
// and single vs multi-signatures.
func VerifySignature(pubKey cryptotypes.PubKey, signerData SignerData, sigData signing.SignatureData, handler SignModeHandler, tx sdk.Tx) error {
	switch data := sigData.(type) {
	case *signing.SingleSignatureData:
		if data.SignMode == signing.SignMode_SIGN_MODE_EIP_712 {
			sig := data.Signature
			sigHash, err := handler.GetSignBytes(data.SignMode, signerData, tx)
			if err != nil {
				return err
			}

			// check signature length
			if len(sig) != ethcrypto.SignatureLength {
				return errors.Wrap(sdkerrors.ErrorInvalidSigner, "signature length doesn't match typical [R||S||V] signature 65 bytes")
			}

			// remove the recovery offset if needed (ie. Metamask eip712 signature)
			if sig[ethcrypto.RecoveryIDOffset] == 27 || sig[ethcrypto.RecoveryIDOffset] == 28 {
				sig[ethcrypto.RecoveryIDOffset] -= 27
			}

			// recover the pubkey from the signature
			feePayerPubkey, err := secp256k1.RecoverPubkey(sigHash, sig)
			if err != nil {
				return errors.Wrap(err, "failed to recover fee payer from sig")
			}
			ecPubKey, err := ethcrypto.UnmarshalPubkey(feePayerPubkey)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal recovered fee payer pubkey")
			}

			// check that the recovered pubkey matches the one in the signerData data
			pk := &ethsecp256k1.PubKey{
				Key: ethcrypto.CompressPubkey(ecPubKey),
			}
			if !pubKey.Equals(pk) {
				return errors.Wrapf(sdkerrors.ErrorInvalidSigner, "feePayer's pubkey %s is different from signature's pubkey %s", pubKey, pk)
			}
			return nil
		} else {
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
