package keys

import (
	"encoding/hex"
	"errors"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

// SignMsgKeysCmd returns the Cobra Command for signing messages with the private key of a given name.
func SignMsgKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign [message]",
		Short: "Sign message",
		Long: `Return a signature from their associated name and address private key.
!!!NOTE!!! 
This is not a secure way to sign messages.
This command is allowed to sign any message from your private key.
Please *DO NOT* use this command unless you know what you are doing.

Example For Signing BLS PoP Message:
	$ gnfd keys add bls --keyring-backend test --algo eth_bls
	$ BLS=$(./build/bin/gnfd keys show bls --keyring-backend test --output json | jq -r .pubkey_hex)
	$ gnfd keys sign $BLS --from bls 
	
`,
		RunE: runSignMsgCmd,
	}

	cmd.Flags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	return cmd
}

func runSignMsgCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid number of arguments")
	}

	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	_, name, _, err := client.GetFromFields(clientCtx, clientCtx.Keyring, clientCtx.From)
	if err != nil {
		return err
	}

	msg, err := hex.DecodeString(args[0])
	if err != nil {
		return err
	}
	sig, _, err := clientCtx.Keyring.Sign(name, tmhash.Sum(msg))
	if err != nil {
		return err
	}

	cmd.Println(hex.EncodeToString(sig))
	return nil
}

// VerifySignatureCmd returns the Cobra Command for verifying signatures with a given public key and message.
func VerifySignatureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify [message] [signature]",
		Short: "Verify signature",
		Long:  "Verify signature with public key and message",
		RunE:  runVerifySignatureCmd,
	}

	cmd.Flags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	return cmd
}

func runVerifySignatureCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid number of arguments")
	}

	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	_, name, _, err := client.GetFromFields(clientCtx, clientCtx.Keyring, clientCtx.From)
	if err != nil {
		return err
	}
	record, err := clientCtx.Keyring.Key(name)
	if err != nil {
		return err
	}

	priv, err := record.ExtractPrivKey()
	if err != nil {
		return nil
	}

	msg, err := hex.DecodeString(args[0])
	if err != nil {
		return nil
	}
	signature, err := hex.DecodeString(args[1])
	if err != nil {
		return nil
	}
	if priv.PubKey().VerifySignature(tmhash.Sum(msg), signature) {
		cmd.Println("Signature verify successfully")
	} else {
		cmd.Println("Signature verify failed")
	}
	return nil
}
