package xhttp

import (
	"context"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/propagation"
)

const (
	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 10 * time.Second
)

type Handler interface {
	Method() string
	Path() string
	Handle(c *gin.Context)
}

type Option func(*Server)

func WithHealthCheck() Option {
	return func(s *Server) {
		s.router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}
}

func WithOpenTracing(serviceName string) Option {
	return func(s *Server) {
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
	return func(s *Server) {
		s.Use(GinRequestZeroLogger(logger))
	}
}

func WithHandlers(handlers ...Handler) Option {
	return func(s *Server) {
		for _, h := range handlers {
			s.Handle(
				h.Method(),
				path.Join(s.basePath, h.Path()),
				h.Handle,
			)
		}
	}
}

type Server struct {
	basePath string

	router *gin.Engine
	srv    *http.Server
}

// Handle adds a new route to the GinServer.
func (s *Server) Handle(method, path string, handler gin.HandlerFunc) {
	s.router.Handle(method, path, handler)
}

// AddMiddleware adds a new middleware to the GinServer.
func (s *Server) Use(middleware gin.HandlerFunc) {
	s.router.Use(middleware)
}

// Run starts the GinServer on the specified port.
func (s *Server) Run(port string) error {
	s.srv = &http.Server{
		Addr:              makeAddr(port),
		Handler:           s.router,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadTimeout,
		WriteTimeout:      DefaultWriteTimeout,
	}
	return s.srv.ListenAndServe()
}

// Shutdown stops the GinServer gracefully.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// NewGinServer creates a new GinServer with default healthcheck route and middlewares.
func NewGinServer(basePath string, opts ...Option) Server {
	router := gin.New()
	router.Use(gin.Recovery())

	server := &Server{
		basePath: basePath,
		router:   router,
	}

	for _, opt := range opts {
		opt(server)
	}

	return *server
}

func makeAddr(port string) string {
	return ":" + port
}
