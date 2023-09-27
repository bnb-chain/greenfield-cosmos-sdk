package server

import (
	"errors"
	"fmt"
	"os"

	"github.com/cometbft/cometbft/node"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

// tmpMigratingDir is a temporary directory to facilitate the migration.
const tmpMigratingDir = "data-migrating"

// NewMigrateStoreCmd creates a command to migrate multistore from IAVL stores to plain DB stores to enable fast node.
func NewMigrateStoreCmd(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "migrate-store",
		Short: "migrate application db to use plain db stores instead of IAVL stores",
		Long: `
To run a fast node, plain DB store types is needed. To convert a normal full node to a fast full node.
We need to migrate the underlying stores. With this command, the old application db will be backed up, 
the new application db will use plain DB store types.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := GetServerContextFromCmd(cmd)
			cfg := ctx.Config
			home := cfg.RootDir
			db, err := openDB(home, GetAppDBBackend(ctx.Viper))
			if err != nil {
				return err
			}
			newDb, err := openDBWithDataDir(home, tmpMigratingDir, GetAppDBBackend(ctx.Viper))
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

			if err = app.CommitMultiStore().LoadLatestVersion(); err != nil {
				return err
			}
			rs, ok := app.CommitMultiStore().(*rootmulti.Store)
			if !ok {
				return errors.New("cannot convert store to root multi store")
			}

			if err = rs.MigrateStores(storetypes.StoreTypeDB, newDb); err != nil {
				return err
			}
			version, err := rootmulti.MigrateCommitInfos(db, newDb)
			if err != nil {
				return err
			}
			fmt.Printf("Multi root store is dumped at version %d \n", version)

			_ = db.Close()
			_ = newDb.Close()

			applicationPath := fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, "data", os.PathSeparator, "application.db")
			applicationBackupPath := fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, "data", os.PathSeparator, "application.db-backup")
			applicationMigratePath := fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, tmpMigratingDir, os.PathSeparator, "application.db")
			if err = os.Rename(applicationPath, applicationBackupPath); err != nil {
				return err
			}
			if err = os.Rename(applicationMigratePath, applicationPath); err != nil {
				return err
			}
			fmt.Printf("Database is replaced and the old one is backup %s\n", applicationBackupPath)

			_ = os.Remove(applicationMigratePath)
			fmt.Printf("Migrate database done, please update app.toml and config.toml to use fastnode feature")

			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	return cmd
}
