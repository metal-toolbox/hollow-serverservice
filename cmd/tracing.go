package cmd

import (
	"github.com/XSAM/otelsql"
	_ "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx" // crdb retries and postgres interface
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracingAndDB() *sqlx.DB {
	dbDriverName := "postgres"

	if viper.GetBool("tracing.enabled") {
		var err error

		initTracer(viper.GetString("tracing.endpoint"))

		// Register an OTel SQL driver
		dbDriverName, err = otelsql.Register(dbDriverName,
			otelsql.WithAttributes(semconv.DBSystemCockroachdb))
		if err != nil {
			logger.Fatalw("failed initializing sql tracer", "error", err)
		}
	}

	db, err := sqlx.Open(dbDriverName, viper.GetString("db.uri"))
	if err != nil {
		logger.Fatalw("failed to initialize database connection", "error", err)
	}

	if err := db.Ping(); err != nil {
		logger.Fatalw("failed verifying database connection", "error", err)
	}

	db.SetMaxOpenConns(viper.GetInt("db.connections.max_open"))
	db.SetMaxIdleConns(viper.GetInt("db.connections.max_idle"))
	db.SetConnMaxIdleTime(viper.GetDuration("db.connections.max_lifetime"))

	return db
}

// initTracer returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func initTracer(url string) *tracesdk.TracerProvider {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		logger.Fatalw("failed to initialize tracing exporter", "error", err)
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("serverservice"),
			attribute.String("environment", viper.GetString("tracing.environment")),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp
}
