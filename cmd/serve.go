package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.metalkube.net/hollow/internal/db"
	"go.metalkube.net/hollow/internal/hollowserver"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the hollow server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("listen", "0.0.0.0:8000", "address to listen on")
	viperBindFlag("listen", serveCmd.Flags().Lookup("listen"))

	serveCmd.Flags().String("db-uri", "postgresql://root@db:26257/hollow_dev?sslmode=disable", "URI for database connection")
	viperBindFlag("db.uri", serveCmd.Flags().Lookup("db-uri"))
}

func serve() {
	if err := db.NewPostgresStore(viper.GetString("db.uri"), logger.Desugar()); err != nil {
		logger.Fatalw("failed to init data store", "error", err)
	}

	logger.Infow("starting server",
		"address", viper.GetString("listen"),
	)

	// gin.SetMode(gin.ReleaseMode)

	hs := &hollowserver.Server{
		Logger: logger.Desugar(),
		Listen: viper.GetString("listen"),
	}
	s := hs.NewServer()

	if err := s.ListenAndServe(); err != nil {
		logger.Fatalw("failed starting server", "error", err)
	}
}
