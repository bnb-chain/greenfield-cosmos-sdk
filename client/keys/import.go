package keys

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/eth/ethsecp256k1"
	"github.com/cosmos/cosmos-sdk/version"
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

	_, err = clientCtx.Keyring.WriteLocalKey(keyName, privKey)
	if err != nil {
		return err
	}
	return nil
}

func ImportKeyHexCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-hex <name> <hex>",
		Short: "Import private keys into the local keybase",
		Long:  fmt.Sprintf("Import hex encoded private key into the local keybase.\nSupported key-types can be obtained with:\n%s list-key-types", version.AppName),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			keyType, _ := cmd.Flags().GetString(flags.FlagKeyType)
			return clientCtx.Keyring.ImportPrivKeyHex(args[0], args[1], keyType)
		},
	}
	cmd.Flags().String(flags.FlagKeyType, string(hd.Secp256k1Type), "private key signing algorithm kind")
	return cmd
}
