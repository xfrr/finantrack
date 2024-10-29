package assetshttp

import (
	"github.com/rs/zerolog"
	"github.com/xfrr/finantrack/internal/shared/xhttp"
	"github.com/xfrr/go-cqrsify/cqrs"
)

//	@title			Asset Management APIs
//	@version		1.0
//	@description	This is the API for managing assets in the Finantrack application.
//	@openapi		3.0.0

//	@contact.name	API Support
//	@contact.url
//	@contact.email

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:6000
//	@BasePath	/api/v1

//	@schemes	http

const BasePath = "/api/v1"

func NewServer(
	serviceName string,
	commandBus cqrs.Bus,
	logger zerolog.Logger,
) xhttp.Server {
	return xhttp.NewGinServer(
		BasePath,
		xhttp.WithHealthCheck(),
		xhttp.WithOpenTracing(serviceName),
		xhttp.WithZeroLogger(&logger),
		xhttp.WithHandlers(
			NewCreateAssetHandler(commandBus),
			NewModifyAssetHandler(commandBus),
			NewDeleteAssetHandler(commandBus),
		),
	)
}
