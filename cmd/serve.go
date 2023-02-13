package cmd

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.hollow.sh/toolbox/ginjwt"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/viperx"
	"gocloud.dev/secrets"

	// import gocdk secret drivers
	_ "gocloud.dev/secrets/localsecrets"

	"go.hollow.sh/serverservice/internal/config"
	"go.hollow.sh/serverservice/internal/dbtools"
	"go.hollow.sh/serverservice/internal/events"
	"go.hollow.sh/serverservice/internal/httpsrv"
)

var (
	apiDefaultListen   = "0.0.0.0:8000"
	natsConnectTimeout = 100 * time.Millisecond
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the hollow server",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("listen", apiDefaultListen, "address to listen on")
	viperx.MustBindFlag(viper.GetViper(), "listen", serveCmd.Flags().Lookup("listen"))

	otelx.MustViperFlags(viper.GetViper(), serveCmd.Flags())
	crdbx.MustViperFlags(viper.GetViper(), serveCmd.Flags())

	// OIDC Flags
	serveCmd.Flags().Bool("oidc", true, "use oidc auth")
	viperx.MustBindFlag(viper.GetViper(), "oidc.enabled", serveCmd.Flags().Lookup("oidc"))
	serveCmd.Flags().String("oidc-aud", "", "expected audience on OIDC JWT")
	viperx.MustBindFlag(viper.GetViper(), "oidc.audience", serveCmd.Flags().Lookup("oidc-aud"))
	serveCmd.Flags().String("oidc-issuer", "", "expected issuer of OIDC JWT")
	viperx.MustBindFlag(viper.GetViper(), "oidc.issuer", serveCmd.Flags().Lookup("oidc-issuer"))
	serveCmd.Flags().String("oidc-jwksuri", "", "URI for JWKS listing for JWTs")
	viperx.MustBindFlag(viper.GetViper(), "oidc.jwksuri", serveCmd.Flags().Lookup("oidc-jwksuri"))
	serveCmd.Flags().String("oidc-roles-claim", "claim", "field containing the permissions of an OIDC JWT")
	viperx.MustBindFlag(viper.GetViper(), "oidc.claims.roles", serveCmd.Flags().Lookup("oidc-roles-claim"))
	serveCmd.Flags().String("oidc-username-claim", "", "additional fields to output in logs from the JWT token, ex (email)")
	viperx.MustBindFlag(viper.GetViper(), "oidc.claims.username", serveCmd.Flags().Lookup("oidc-username-claim"))
	// DB Flags
	serveCmd.Flags().String("db-encryption-driver", "", "encryption driver uri; 32 byte base64 encoded string, (example: base64key://your-encoded-secret-key)")
	viperx.MustBindFlag(viper.GetViper(), "db.encryption_driver", serveCmd.Flags().Lookup("db-encryption-driver"))

	// NATs Flags
	rootCmd.PersistentFlags().String("nats-url", "", "NATS server connection url")
	viperx.MustBindFlag(viper.GetViper(), "nats.url", rootCmd.PersistentFlags().Lookup("nats-url"))

	rootCmd.PersistentFlags().String("nats-stream-user", "serverservice", "NATS basic auth account user name")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream-user", rootCmd.PersistentFlags().Lookup("nats-stream-user"))

	rootCmd.PersistentFlags().String("nats-stream-pass", "password", "NATS basic auth account password")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream-pass", rootCmd.PersistentFlags().Lookup("nats-stream-pass"))

	rootCmd.PersistentFlags().String("nats-creds-file", "", "Path to the file containing the NATS nkey keypair")
	viperx.MustBindFlag(viper.GetViper(), "nats.creds-file", rootCmd.PersistentFlags().Lookup("nats-creds-file"))

	rootCmd.PersistentFlags().String("nats-stream-name", appName, "prefix for NATS subjects")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream-name", rootCmd.PersistentFlags().Lookup("nats-stream-name"))

	rootCmd.PersistentFlags().String("nats-stream-prefix", "com.hollow.sh.events", "NATS stream prefix")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream-prefix", rootCmd.PersistentFlags().Lookup("nats-stream-prefix"))

	rootCmd.PersistentFlags().StringSlice("nats-stream-subjects", []string{"com.hollow.sh.events.>"}, "NATS stream subject(s)")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream-subjects", rootCmd.PersistentFlags().Lookup("nats-stream-subjects"))

	rootCmd.PersistentFlags().Duration("nats-connect-timeout", natsConnectTimeout, "Timeout when connecting to NATs")
	viperx.MustBindFlag(viper.GetViper(), "nats.connect-timeout", rootCmd.PersistentFlags().Lookup("nats-connect-timeout"))
}

func serve(ctx context.Context) {
	err := otelx.InitTracer(config.AppConfig.Tracing, appName, logger)
	if err != nil {
		logger.Fatalw("unable to initialize tracing system", "error", err)
	}

	db := initDB()

	dbtools.RegisterHooks()

	keeper, err := secrets.OpenKeeper(ctx, viper.GetString("db.encryption_driver"))
	if err != nil {
		logger.Fatalw("failed to open secrets keeper", "error", err)
	}
	defer keeper.Close()

	logger.Infow("starting server",
		"address", viper.GetString("listen"),
	)

	hs := &httpsrv.Server{
		Logger:        logger.Desugar(),
		Listen:        viper.GetString("listen"),
		Debug:         config.AppConfig.Logging.Debug,
		DB:            db,
		SecretsKeeper: keeper,
		AuthConfig: ginjwt.AuthConfig{
			Enabled:       viper.GetBool("oidc.enabled"),
			Audience:      viper.GetString("oidc.audience"),
			Issuer:        viper.GetString("oidc.issuer"),
			JWKSURI:       viper.GetString("oidc.jwksuri"),
			LogFields:     viper.GetStringSlice("oidc.log"),
			RolesClaim:    viper.GetString("oidc.claims.roles"),
			UsernameClaim: viper.GetString("oidc.claims.username"),
		},
	}

	// init event stream - for now, only when nats.url is specified
	streamURL := viper.GetString("nats.url")
	if streamURL != "" {
		eventStream := events.NewStreamBroker(
			appName,
			viper.GetString("nats.creds-file"),
			streamURL,
			viper.GetString("nats.stream-name"),
			viper.GetString("nats.stream-prefix"),
			viper.GetStringSlice("nats.stream-subjects"),
			// TODO(joel): fix these params are not being populated
			viper.GetString("nats.stream-user"),
			viper.GetString("nats.stream-pass"),
			viper.GetDuration("nats.connect-timeout"),
		)

		if err := eventStream.Open(); err != nil {
			logger.Warnw("failed to open event stream", "error", err.Error())
		} else {
			hs.EventStream = eventStream
			defer hs.EventStream.Close()
		}
	}

	if err := hs.Run(); err != nil {
		logger.Fatalw("failed starting server", "error", err)
	}
}

func initDB() *sqlx.DB {
	dbDriverName := "postgres"

	sqldb, err := crdbx.NewDB(config.AppConfig.CRDB, config.AppConfig.Tracing.Enabled)
	if err != nil {
		logger.Fatalw("failed to initialize database connection", "error", err)
	}

	db := sqlx.NewDb(sqldb, dbDriverName)

	return db
}
