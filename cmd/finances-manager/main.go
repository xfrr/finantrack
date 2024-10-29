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

		dbHost   = xos.GetEnvWithDefault("FINANCES_MANAGER_DB_HOST", "localhost")
		dbPort   = xos.GetEnvWithDefault("FINANCES_MANAGER_DB_PORT", "27017")
		dbUser   = xos.GetEnvWithDefault("FINANCES_MANAGER_DB_USER", "root")
		dbPass   = xos.GetEnvWithDefault("FINANCES_MANAGER_DB_PASS", "root")
		dbName   = xos.GetEnvWithDefault("FINANCES_MANAGER_DB_NAME", "finantrack")
		dbEngine = services.DatabaseEngineType(xos.GetEnvWithDefault(
			"FINANCES_MANAGER_DB_ENGINE",
			string(services.MongoDatabaseEngine)),
		)
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
			services.DatabaseHost(dbHost),
			services.DatabasePort(dbPort),
			services.DatabaseUser(dbUser),
			services.DatabasePass(dbPass),
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
