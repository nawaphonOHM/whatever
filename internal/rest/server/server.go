package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
		Addr:         fmt.Sprintf("%s:%d", c.Host, c.Port),
		Handler:      engine,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	}
}

// New creates and initializes a new Server instance.
func New(cfg *Config) *Server {
	c := resolveConfig(cfg)
	applyServerDefaults(c)
	engine := gin.New()
	return &Server{
		Engine:     engine,
		Config:     c,
		httpServer: buildHTTPServer(c, engine),
	}
}
