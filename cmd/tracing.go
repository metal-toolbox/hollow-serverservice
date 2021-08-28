package cmd

import (
	"database/sql"

	"github.com/XSAM/otelsql"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracingAndDB() *sql.DB {
	dbDriverName := "postgres"

	if viper.GetBool("tracing.enabled") {
		var err error

		initTracer(viper.GetString("tracing.endpoint"))

		// Register an OTel SQL driver
		dbDriverName, err = otelsql.Register(dbDriverName, semconv.DBSystemCockroachdb.Value.AsString())
		if err != nil {
			logger.Fatalw("failed initializing sql tracer", "error", err)
		}
	}

	db, err := sql.Open(dbDriverName, viper.GetString("db.uri"))
	if err != nil {
		logger.Fatalw("failed to initialize database connection", "error", err)
	}

	if _, err := db.Exec("select 1;"); err != nil {
		logger.Fatalw("failed verifying database connection", "error", err)
	}

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
