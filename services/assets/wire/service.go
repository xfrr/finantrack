package assets

import (
	"context"

	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
	"github.com/xfrr/finantrack/internal/shared/xlog"
	"github.com/xfrr/finantrack/internal/shared/xtracing"
	"github.com/xfrr/finantrack/services"
	assetshttp "github.com/xfrr/finantrack/services/assets/http"
)

type Service struct {
	services.Base

	repoFactory services.RepositoryFactory[assetdomain.Repository]
}

func (s Service) Start(ctx context.Context) error {
	logger := xlog.NewZerologger(s.Name(), s.Config().Environment)
	logger.Debug().
		Any("config", s.Config()).
		Msg("starting assets service...")

	// create new tracer provider
	tracer, stopTracer, err := xtracing.NewOtelTracerProvider(ctx, s.Name(), s.Config().OtelCollectorURL)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create tracer provider")
	}

	// create database based on the engine type
	repository, stopDatabase, err := s.repoFactory.CreateRepository(ctx, s.Config().DatabaseEngine)
	if err != nil {
		return err
	}

	// creates new command bus and register all commands
	cmdbus, err := newAssetCommandBus(ctx, repository, tracer)
	if err != nil {
		return err
	}

	// create new http server instance
	httpServer := assetshttp.NewServer(
		s.Name(),
		cmdbus,
		logger,
	)

	go func() {
		<-ctx.Done()
		logger.Info().Msg("shutting down assets service")

		// gracefully shutdown http server
		err = httpServer.Shutdown(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("failed to shutdown http server")
		}

		// stop database connection
		err = stopDatabase()
		if err != nil {
			logger.Error().Err(err).Msg("failed to close database connection")
		}

		// stop tracer
		stopTracer()
	}()

	return httpServer.Run(s.Config().HTTPServerPort)
}

func NewService(opts ...services.InitializeOption) (*Service, error) {
	var (
		err error
	)

	service := &Service{
		Base: *services.NewService("assets", opts...),
	}

	// register all events for the assets context
	eventsRegistry := newAssetEventsRegistry()

	// Register asset repository factory
	service.repoFactory, err = newRepositoryFactory(
		service.Config().DatabaseURI,
		service.Config().DatabaseName,
		eventsRegistry,
	)
	if err != nil {
		return nil, err
	}

	return service, nil
}
