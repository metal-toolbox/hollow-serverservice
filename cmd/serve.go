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

	serveCmd.Flags().String("jwt-aud", "", "expected audience on JWT tokens")
	viperBindFlag("jwt.audience", serveCmd.Flags().Lookup("jwt-aud"))
	serveCmd.Flags().String("jwt-issuer", "https://equinixmetal.us.auth0.com/", "expected issuer of JWT tokens")
	viperBindFlag("jwt.issuer", serveCmd.Flags().Lookup("jwt-issuer"))
	serveCmd.Flags().String("jwt-jwksuri", "https://equinixmetal.us.auth0.com/.well-known/jwks.json", "URI for JWKS listing for JWTs")
	viperBindFlag("jwt.jwksuri", serveCmd.Flags().Lookup("jwt-jwksuri"))
}

func serve() {
	store, err := db.NewPostgresStore(viper.GetString("db.uri"), logger.Desugar())
	if err != nil {
		logger.Fatalw("failed to init data store", "error", err)
	}

	logger.Infow("starting server",
		"address", viper.GetString("listen"),
	)

	hs := &hollowserver.Server{
		Logger: logger.Desugar(),
		Listen: viper.GetString("listen"),
		Debug:  viper.GetBool("logging.debug"),
		Store:  store,
		AuthConfig: hollowserver.AuthConfig{
			Audience: viper.GetString("jwt.audience"),
			Issuer:   viper.GetString("jwt.issuer"),
			JWKSURI:  viper.GetString("jwt.jwksuri"),
		},
	}

	if err := hs.Run(); err != nil {
		logger.Fatalw("failed starting server", "error", err)
	}
}
