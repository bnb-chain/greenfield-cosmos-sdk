package keys

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/eth/ethsecp256k1"
)

const (
	flagSecp256k1PrivateKey = "secp256k1-private-key"
)

// ImportKeyCommand imports private keys from a keyfile.
func ImportKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <name> <keyfile>/<privateKey>",
		Short: "Import private keys into the local keybase",
		Long:  "Import a ASCII armored/Secp256k1 private key into the local keybase.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			isSecp256k1, _ := cmd.Flags().GetBool(flagSecp256k1PrivateKey)

			if !isSecp256k1 {
				return importASCIIArmored(clientCtx, args)
			}

			return importSecp256k1(clientCtx, args)
		},
	}

	cmd.Flags().Bool(flagSecp256k1PrivateKey, false, "import Secp256k1 format private key")

	return cmd
}

func importASCIIArmored(clientCtx client.Context, args []string) error {
	buf := bufio.NewReader(clientCtx.Input)

	bz, err := os.ReadFile(args[1])
	if err != nil {
		return err
	}

	passphrase, err := input.GetPassword("Enter passphrase to decrypt your key:", buf)
	if err != nil {
		return err
	}

	return clientCtx.Keyring.ImportPrivKey(args[0], string(bz), passphrase)
}

func importSecp256k1(clientCtx client.Context, args []string) error {
	keyName := args[0]
	keyBytes, err := hex.DecodeString(args[1])
	if err != nil {
		return err
	}
	if len(keyBytes) != 32 {
		return fmt.Errorf("len of keybytes is not equal to 32")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], keyBytes[:32])
	privKey := hd.EthSecp256k1.Generate()(keyBytesArray[:]).(*ethsecp256k1.PrivKey)

	clientCtx.Keyring.WriteLocalKey(keyName, privKey)
	return nil
}
