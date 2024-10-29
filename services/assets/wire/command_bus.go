package assets

import (
	"context"

	"github.com/xfrr/go-cqrsify/cqrs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	assetscommands "github.com/xfrr/finantrack/internal/contexts/assets/commands"
	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
)

// newAssetCommandBus creates a new command bus for the assets context
// and registers all command handlers.
func newAssetCommandBus(
	ctx context.Context,
	repository assetdomain.Repository,
	tracer trace.Tracer,
) (cqrs.Bus, error) {
	bus := cqrs.NewBus()
	bus.Use(func(f func(context.Context, interface{}) (interface{}, error)) func(context.Context, interface{}) (interface{}, error) {
		return func(ctx context.Context, command interface{}) (interface{}, error) {
			var attrs []attribute.KeyValue

			name := "command"
			if cmd, isCommand := command.(cqrs.Command); isCommand {
				name = cmd.CommandName()
				attrs = append(attrs, attribute.String("command.name", cmd.CommandName()))
				attrs = append(attrs, attribute.String("command.bus.type", "inmemory"))
				attrs = append(attrs, attribute.String("command.bus.name", "cqrsify"))
				attrs = append(attrs, attribute.String("command.bus.service", "assets"))
			}

			ctx, span := tracer.Start(ctx, name)
			defer span.End()
			span.SetAttributes(attrs...)

			return f(ctx, command)
		}
	})

	// Register command handlers
	err := cqrs.Handle(ctx, bus, assetscommands.NewCreateAssetCommandHandler(repository).Handle)
	if err != nil {
		return nil, err
	}

	err = cqrs.Handle(ctx, bus, assetscommands.NewDeleteAssetCommandHandler(repository).Handle)
	if err != nil {
		return nil, err
	}

	return bus, nil
}
