package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/xfrr/finantrack/internal/shared/xos"
	"github.com/xfrr/finantrack/services"

	assets "github.com/xfrr/finantrack/services/assets/wire"
)

func main() {
	var (
		httpServerPort   = xos.GetEnvWithDefault("FINANCES_MANAGER_HTTP_SERVER_PORT", "6000")
		environment      = xos.GetEnvWithDefault("FINANCES_MANAGER_ENVIRONMENT", "development")
		otelCollectorURL = xos.GetEnvWithDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")

		dbURI    = xos.GetEnvWithDefault("FINANCES_MANAGER_DB_URI", "mongodb://localhost:27017")
		dbName   = xos.GetEnvWithDefault("FINANCES_MANAGER_DB_NAME", "finantrack")
		dbEngine = services.DatabaseEngineType(xos.GetEnvWithDefault("FINANCES_MANAGER_DB_ENGINE", string(services.MongoDatabaseEngine)))
	)

	ctx, stopNotification := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopNotification()

	service, err := assets.NewService(
		services.Environment(environment),
		services.Traces(otelCollectorURL),
		services.HTTPServer(
			services.Port(httpServerPort),
		),
		services.Database(
			services.DatabaseEngine(dbEngine),
			services.DatabaseURI(dbURI),
			services.DatabaseName(dbName),
		),
	)
	if err != nil {
		panic(err)
	}

	if err = service.Start(ctx); err != nil {
		panic(err)
	}
}
