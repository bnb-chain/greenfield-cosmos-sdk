package server

import (
	"fmt"

	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	"github.com/cometbft/cometbft/node"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/server/types"
)

// NewRollbackCmd creates a command to rollback tendermint and multistore state by one height.
func NewRollbackCmd(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {
	rollbackBlocks := uint(1)
	var removeBlock bool

	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "rollback cosmos-sdk and tendermint state by blocks",
		Long: `
A state rollback is performed to recover from an incorrect application state transition,
when Tendermint has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - blocks.
The application also rolls back to height n - blocks. No blocks are removed, so upon
restarting Tendermint the transactions in block n will be re-executed against the
application.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := GetServerContextFromCmd(cmd)
			cfg := ctx.Config
			home := cfg.RootDir
			db, err := openDB(home, GetAppDBBackend(ctx.Viper))
			if err != nil {
				return err
			}
			config, err := serverconfig.GetConfig(ctx.Viper)
			if err != nil {
				return err
			}
			genDocProvider := node.DefaultGenesisDocProviderFunc(ctx.Config)
			genDoc, err := genDocProvider()
			if err != nil {
				return err
			}
			app := appCreator(ctx.Logger, db, nil, genDoc.ChainID, &config, ctx.Viper)
			// rollback CometBFT state
			height, hash, err := cmtcmd.RollbackState(ctx.Config, removeBlock, int64(rollbackBlocks))
			if err != nil {
				return fmt.Errorf("failed to rollback tendermint state: %w", err)
			}
			// rollback the multistore

			if err := app.CommitMultiStore().RollbackToVersion(height); err != nil {
				return fmt.Errorf("failed to rollback to version: %w", err)
			}

			fmt.Printf("Rolled back state to height %d and hash %X", height, hash)
			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().UintVar(&rollbackBlocks, "blocks", 1, "number of blocks to rollback")
	cmd.Flags().BoolVar(&removeBlock, "hard", false, "remove last block as well as state")

	return cmd
}
