package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.metalkube.net/hollow/internal/db"
)

// migrateCmd represents the serve command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate the datastore schema",
	Run: func(cmd *cobra.Command, args []string) {
		migrate()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().String("db-uri", "postgresql://root@db:26257/hollow_dev?sslmode=disable", "URI for database connection")
	viperBindFlag("db.uri", migrateCmd.Flags().Lookup("db-uri"))
}

func migrate() {
	store, err := db.NewPostgresStore(viper.GetString("db.uri"), logger.Desugar())
	if err != nil {
		logger.Fatalw("failed to init data store", "error", err)
	}

	if err := store.Migrate(); err != nil {
		logger.Fatalw("failed to migrate data store", "error", err)
	}
}
