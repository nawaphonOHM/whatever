package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/logging/callstack"
	"github.com/nawaphonOHM/whatever/pkg/logging"
)

// Server encapsulates the Gin engine and HTTP server lifecycle.
type Server struct {
	Engine     *gin.Engine
	Config     *Config
	httpServer *http.Server
}

// resolveConfig returns cfg or defaults when cfg is nil.
func resolveConfig(cfg *Config) *Config {
	if cfg != nil {
		return cfg
	}
	logging.Info("nil REST server configuration provided; falling back to default configuration")
	return DefaultConfig()
}

// ensureShutdownTimeout sets a positive shutdown timeout.
func ensureShutdownTimeout(c *Config) {
	if c.ShutdownTimeout > 0 {
		return
	}
	c.ShutdownTimeout = defaultTimeoutSec * time.Second
}

// applyGinMode sets gin mode when configured.
func applyGinMode(c *Config) {
	if c.Mode == "" {
		return
	}
	gin.SetMode(c.Mode)
}

// applyServerDefaults normalizes config before server construction.
func applyServerDefaults(c *Config) {
	applyGinMode(c)
	ensureShutdownTimeout(c)
}

// buildHTTPServer constructs the net/http server for cfg and engine.
func buildHTTPServer(c *Config, engine *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", c.Host, c.Port),
		Handler:           engine,
		ReadTimeout:       c.ReadTimeout,
		ReadHeaderTimeout: c.ReadHeaderTimeout,
		WriteTimeout:      c.WriteTimeout,
		IdleTimeout:       c.IdleTimeout,
		MaxHeaderBytes:    c.MaxHeaderBytes,
	}
}

// New creates and initializes a new Server instance.
func New(cfg *Config) (*Server, error) {
	decorated := callstack.DecorateFuncErr("server.New", func() (*Server, error) {
		logging.Info("entering REST server initialization state")
		c := resolveConfig(cfg)
		applyServerDefaults(c)
		logging.Info("REST server configuration choices",
			"host", c.Host,
			"port", c.Port,
			"mode", c.Mode,
			"shutdown_timeout", c.ShutdownTimeout.String(),
			"access_log", c.EnableAccessLog,
			"metrics", c.EnableMetrics,
			"profiling", c.EnableProfiling,
		)
		engine := gin.New()
		if err := applyServerOptions(engine, c); err != nil {
			return nil, err
		}
		logging.Info("REST server initialization complete", "addr", fmt.Sprintf("%s:%d", c.Host, c.Port))
		return &Server{
			Engine:     engine,
			Config:     c,
			httpServer: buildHTTPServer(c, engine),
		}, nil
	})
	return decorated()
}
