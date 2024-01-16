package cmd

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"go.hollow.sh/toolbox/events"
	"go.hollow.sh/toolbox/ginjwt"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/viperx"
	"gocloud.dev/secrets"

	// import gocdk secret drivers
	_ "gocloud.dev/secrets/localsecrets"

	"go.hollow.sh/serverservice/internal/config"
	"go.hollow.sh/serverservice/internal/dbtools"
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
	ginjwt.BindFlagFromViperInst(viper.GetViper(), "oidc.enabled", serveCmd.Flags().Lookup("oidc"))

	// DB Flags
	serveCmd.Flags().String("db-encryption-driver", "", "encryption driver uri; 32 byte base64 encoded string, (example: base64key://your-encoded-secret-key)")
	viperx.MustBindFlag(viper.GetViper(), "db.encryption_driver", serveCmd.Flags().Lookup("db-encryption-driver"))

	// NATs Flags
	rootCmd.PersistentFlags().String("nats-url", "", "NATS server connection url")
	viperx.MustBindFlag(viper.GetViper(), "nats.url", rootCmd.PersistentFlags().Lookup("nats-url"))

	rootCmd.PersistentFlags().String("nats-stream-user", "", "NATS basic auth account user name")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream.user", rootCmd.PersistentFlags().Lookup("nats-stream-user"))

	rootCmd.PersistentFlags().String("nats-stream-pass", "", "NATS basic auth account password")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream.pass", rootCmd.PersistentFlags().Lookup("nats-stream-pass"))

	rootCmd.PersistentFlags().String("nats-creds-file", "", "Path to the file containing the NATS nkey keypair")
	viperx.MustBindFlag(viper.GetViper(), "nats.creds.file", rootCmd.PersistentFlags().Lookup("nats-creds-file"))

	rootCmd.PersistentFlags().String("nats-stream-name", appName, "prefix for NATS subjects")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream.name", rootCmd.PersistentFlags().Lookup("nats-stream-name"))

	rootCmd.PersistentFlags().String("nats-stream-prefix", "com.hollow.sh.serverservice.events", "NATS stream prefix")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream.prefix", rootCmd.PersistentFlags().Lookup("nats-stream-prefix"))

	rootCmd.PersistentFlags().StringSlice("nats-stream-subjects", []string{"com.hollow.sh.serverservice.events.>"}, "NATS stream subject(s)")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream.subjects", rootCmd.PersistentFlags().Lookup("nats-stream-subjects"))

	rootCmd.PersistentFlags().String("nats-stream-urn-ns", "hollow", "NATS stream URN namespace value")
	viperx.MustBindFlag(viper.GetViper(), "nats.stream.urn.ns", rootCmd.PersistentFlags().Lookup("nats-stream-urn-ns"))

	rootCmd.PersistentFlags().Duration("nats-connect-timeout", natsConnectTimeout, "Timeout when connecting to NATs")
	viperx.MustBindFlag(viper.GetViper(), "nats.connect.timeout", rootCmd.PersistentFlags().Lookup("nats-connect-timeout"))
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

	logger.Infow("oidc enabled", "oidc", viper.GetString("oidc"))

	var authCfgs []ginjwt.AuthConfig
	if viper.GetViper().GetBool("oidc.enabled") {
		authCfgs, err = ginjwt.GetAuthConfigsFromFlags(viper.GetViper())
		if err != nil {
			logger.Fatal(err)
		}
	}

	hs := &httpsrv.Server{
		Logger:      logger.Desugar(),
		Listen:      viper.GetString("listen"),
		Debug:       config.AppConfig.Logging.Debug,
		DB:          db,
		AuthConfigs: authCfgs,
	}

	// init event stream - for now, only when nats.url is specified
	eventstream := initStream()
	if eventstream != nil {
		hs.EventStream = eventstream
		defer hs.EventStream.Close()
	}

	if err := hs.Run(); err != nil {
		logger.Fatalw("failed starting server", "error", err)
	}
}

func initStream() events.Stream {
	streamURL := viper.GetString("nats.url")
	if streamURL == "" {
		return nil
	}

	stream, err := events.NewStream(natsOptions(appName, streamURL))
	if err != nil {
		logger.Warnw("error in event stream configuration", "error", err.Error())

		return nil
	}

	if err := stream.Open(); err != nil {
		logger.Warnw("error in event stream configuration", "error", err.Error())

		return nil
	}

	return stream
}

func natsOptions(appName, serverURL string) events.NatsOptions {
	return events.NatsOptions{
		AppName:                appName,
		URL:                    serverURL,
		StreamUser:             viper.GetString("nats.stream.user"),
		StreamPass:             viper.GetString("nats.stream.pass"),
		CredsFile:              viper.GetString("nats.creds.file"),
		PublisherSubjectPrefix: viper.GetString("nats.stream.prefix"),
		StreamURNNamespace:     viper.GetString("nats.stream.urn.ns"),
		ConnectTimeout:         viper.GetDuration("nats.connect.timeout"),
		Stream: &events.NatsStreamOptions{
			Name:     viper.GetString("nats.stream.name"),
			Subjects: viper.GetStringSlice("nats.stream.subjects"),
		},
	}
}

func initDB() *sqlx.DB {
	dbDriverName := "postgres"

	sqldb, err := crdbx.NewDB(config.AppConfig.CRDB, config.AppConfig.Tracing.Enabled)
	if err != nil {
		logger.Fatalw("failed to initialize database connection", "error", err)
	}

	boil.SetDB(sqldb)

	db := sqlx.NewDb(sqldb, dbDriverName)

	return db
}
