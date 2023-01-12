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
			sigHash, err := handler.GetSignBytes(data.SignMode, signerData, tx)
			if err != nil {
				return err
			}

			senderSig := data.Signature
			if len(senderSig) != ethcrypto.SignatureLength {
				return errors.Wrap(sdkerrors.ErrorInvalidSigner, "signature length doesn't match typical [R||S||V] signature 65 bytes")
			}

			// Remove the recovery offset if needed (ie. Metamask eip712 signature)
			if senderSig[ethcrypto.RecoveryIDOffset] == 27 || senderSig[ethcrypto.RecoveryIDOffset] == 28 {
				senderSig[ethcrypto.RecoveryIDOffset] -= 27
			}

			feePayerPubkey, err := secp256k1.RecoverPubkey(sigHash, senderSig)
			if err != nil {
				return errors.Wrap(err, "failed to recover delegated fee payer from sig")
			}

			ecPubKey, err := ethcrypto.UnmarshalPubkey(feePayerPubkey)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal recovered fee payer pubkey")
			}

			pk := &ethsecp256k1.PubKey{
				Key: ethcrypto.CompressPubkey(ecPubKey),
			}

			if !pubKey.Equals(pk) {
				return errors.Wrapf(sdkerrors.ErrInvalidPubKey, "feePayer pubkey %s is different from transaction pubkey %s", pubKey, pk)
			}

			recoveredFeePayerAcc := sdk.AccAddress(pk.Address().Bytes())

			if !recoveredFeePayerAcc.Equals(sdk.MustAccAddressFromHex(signerData.Address)) {
				return errors.Wrapf(sdkerrors.ErrorInvalidSigner, "failed to verify delegated fee payer %s signature", recoveredFeePayerAcc)
			}

			// VerifySignature of ethsecp256k1 accepts 64 byte signature [R||S]
			// WARNING! Under NO CIRCUMSTANCES try to use pubKey.VerifySignature there
			if !secp256k1.VerifySignature(pubKey.Bytes(), sigHash, senderSig[:len(senderSig)-1]) {
				return errors.Wrap(sdkerrors.ErrorInvalidSigner, "unable to verify signer signature of EIP712 typed data")
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
