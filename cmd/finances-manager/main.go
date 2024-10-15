package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/xfrr/finantrack/pkg/xlog"
	"github.com/xfrr/finantrack/pkg/xos"
	"github.com/xfrr/finantrack/pkg/xtracing"

	ftkhttp "github.com/xfrr/finantrack/http"
)

// env variables
var (
	httpServerPort   = xos.GetEnvWithDefault("FINANCES_MANAGER_HTTP_SERVER_PORT", "6000")
	environment      = xos.GetEnvWithDefault("FINANCES_MANAGER_ENVIRONMENT", "development")
	serviceName      = xos.GetEnvWithDefault("FINANCES_MANAGER_SERVICE_NAME", "finances-manager")
	otelCollectorURL = xos.GetEnvWithDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "otelcol:4317")
)

func main() {
	logger := xlog.NewZerologger(serviceName, environment)
	logger.Debug().Msg("starting finances manager service")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tracer, err := xtracing.InitOpenTelemetryTracer(ctx, serviceName, otelCollectorURL)
	if err != nil {
		panic(err)
	}
	defer tracer.Shutdown(ctx)

	httpServer := ftkhttp.NewGinServer(
		ftkhttp.WithRecovery(),
		ftkhttp.WithHealthCheck(),
		ftkhttp.WithOpenTracing(serviceName),
		ftkhttp.WithZeroLogger(&logger),
	)

	go httpServer.Run(httpServerPort)

	<-ctx.Done()
}
