package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	intcfg "github.com/example/go-boilerplate/internal/config"
	intmw "github.com/example/go-boilerplate/internal/middleware"
	"github.com/example/go-boilerplate/pkg/logger"
	"github.com/example/go-boilerplate/pkg/rest/response"
	"github.com/gin-gonic/gin"
)

// Config defines the configuration for the HTTP server.
// All environment variable keys use the OHM9969_ prefix.
type Config struct {
	Host            string        `env:"OHM9969_SERVER_HOST" envDefault:""`
	Port            int           `env:"OHM9969_SERVER_PORT" envDefault:"8080"`
	Mode            string        `env:"OHM9969_GIN_MODE" envDefault:"release"`
	AppVersion      string        `env:"OHM9969_APP_VERSION" envDefault:""`
	ReadTimeout     time.Duration `env:"OHM9969_SERVER_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout    time.Duration `env:"OHM9969_SERVER_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout     time.Duration `env:"OHM9969_SERVER_IDLE_TIMEOUT" envDefault:"60s"`
	ShutdownTimeout time.Duration `env:"OHM9969_SERVER_SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

// DefaultConfig returns server configuration with recommended defaults.
func DefaultConfig() *Config {
	return &Config{
		Host:            "",
		Port:            8080,
		Mode:            gin.ReleaseMode,
		AppVersion:      "",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}
}

// Server encapsulates the Gin engine and HTTP server lifecycle.
type Server struct {
	Engine     *gin.Engine
	Config     *Config
	httpServer *http.Server
}

// New creates and initializes a new Server instance.
func New(cfg *Config) *Server {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if cfg.Mode != "" {
		gin.SetMode(cfg.Mode)
	}

	if cfg.ShutdownTimeout <= 0 {
		cfg.ShutdownTimeout = 10 * time.Second
	}

	engine := gin.New()

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		Engine:     engine,
		Config:     cfg,
		httpServer: httpServer,
	}
}

// NewFromRegistrations loads environment configuration with OHM9969_ prefix,
// validates registrations, attaches default middlewares, mounts reserved health endpoints,
// and prepares a Server ready to start.
func NewFromRegistrations(registrations []*RestApiRegistration) (*Server, error) {
	cfg, err := intcfg.Load[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to load server config: %w", err)
	}

	srv := New(cfg)
	srv.SetupDefaultMiddlewares()

	if err := srv.RegisterRoutesWithVersion(registrations, cfg.AppVersion); err != nil {
		return nil, err
	}

	return srv, nil
}

// Start launches the HTTP server and blocks until the context is canceled or
// an OS interrupt signal (SIGINT, SIGTERM) is received, triggering graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("http server failed to start: %w", err)
		}
		return nil
	case <-ctx.Done():
		// Signal or context cancellation caught, begin graceful shutdown
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.Config.ShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	return nil
}

// Shutdown initiates an immediate graceful shutdown of the server using the provided context.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// SetupDefaultMiddlewares attaches recommended middlewares (RequestID, Logger,
// Recovery, CORS) and standard 404/405 RFC 9457 error handlers to the Gin engine.
func (s *Server) SetupDefaultMiddlewares(loggerSkipPaths ...string) {
	s.Engine.HandleMethodNotAllowed = true

	skipPaths := []string{ReservedHealthPath, ReservedReadyPath}
	if len(loggerSkipPaths) > 0 {
		skipPaths = append(skipPaths, loggerSkipPaths...)
	}

	s.Engine.Use(
		intmw.RequestID(),
		logger.LoggerWithConfig(logger.LoggerConfig{
			SkipPaths: skipPaths,
		}),
		intmw.Recovery(),
		intmw.CORS(intmw.DefaultCORSConfig()),
	)

	s.Engine.NoRoute(func(c *gin.Context) {
		resp := response.NotFound("NOT_FOUND", "Route not found")
		resp.Write(c)
	})

	s.Engine.NoMethod(func(c *gin.Context) {
		resp := response.Error(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		resp.Write(c)
	})
}

// RegisterRoutes validates and mounts caller-supplied API registrations onto the server's Gin engine,
// pre-registering the reserved /health and /ready health endpoints.
func (s *Server) RegisterRoutes(registrations []*RestApiRegistration) error {
	return RegisterRoutesWithVersion(s.Engine, registrations, "")
}

// RegisterRoutesWithVersion validates and mounts caller-supplied API registrations with version metadata,
// pre-registering the reserved /health and /ready health endpoints.
func (s *Server) RegisterRoutesWithVersion(registrations []*RestApiRegistration, version string) error {
	return RegisterRoutesWithVersion(s.Engine, registrations, version)
}
