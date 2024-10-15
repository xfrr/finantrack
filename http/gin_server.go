package ftkhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/propagation"
)

type Handler interface {
	Method() string
	Path() string
	Handle(c *gin.Context)
}

type Option func(*GinServer)

func WithHealthCheck() Option {
	return func(s *GinServer) {
		s.srv.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}
}

func WithOpenTracing(serviceName string) Option {
	return func(s *GinServer) {
		s.Use(otelgin.Middleware(serviceName,
			otelgin.WithFilter(func(r *http.Request) bool {
				return r.URL.Path != "/health"
			}),
			otelgin.WithPropagators(
				propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{}, propagation.Baggage{},
				),
			),
		))
	}
}

func WithZeroLogger(logger *zerolog.Logger) Option {
	return func(s *GinServer) {
		s.Use(GinRequestZeroLogger(logger))
	}
}

func WithRecovery() Option {
	return func(s *GinServer) {
		s.Use(gin.Recovery())
	}
}

func WithHandlers(handlers ...Handler) Option {
	return func(s *GinServer) {
		for _, h := range handlers {
			s.Handle(h.Method(), h.Path(), h.Handle)
		}
	}
}

type GinServer struct {
	srv *gin.Engine
}

// Handle adds a new route to the GinServer.
func (s *GinServer) Handle(method, path string, handler gin.HandlerFunc) {
	s.srv.Handle(method, path, handler)
}

// AddMiddleware adds a new middleware to the GinServer.
func (s *GinServer) Use(middleware gin.HandlerFunc) {
	s.srv.Use(middleware)
}

// Run starts the GinServer on the specified port.
func (s *GinServer) Run(port string) error {
	return s.srv.Run(makeAddr(port))
}

// NewGinServer creates a new GinServer with default healthcheck route and middlewares.
func NewGinServer(opts ...Option) *GinServer {
	engine := gin.Default()
	for _, opt := range opts {
		opt(&GinServer{engine})
	}
	return &GinServer{engine}
}

func makeAddr(port string) string {
	return ":" + port
}
