// Package server implements the internal HTTP engine and lifecycle routines.
package server

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultServerPort     = 8080
	defaultTimeoutSec     = 10
	defaultIdleTimeoutSec = 60
)

// dur shortens struct field lines for tag length limits.
type dur = time.Duration

// timeoutFields groups duration settings for the HTTP server.
type timeoutFields struct {
	ReadTimeout     dur `env:"OHM9996_SERVER_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout    dur `env:"OHM9996_SERVER_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout     dur `env:"OHM9996_SERVER_IDLE_TIMEOUT" envDefault:"60s"`
	ShutdownTimeout dur `env:"OHM9996_SERVER_SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

// Config defines the configuration for the HTTP server.
// All environment variable keys use the OHM9996_ prefix.
type Config struct {
	Host       string `env:"OHM9996_SERVER_HOST" envDefault:""`
	Mode       string `env:"OHM9996_GIN_MODE" envDefault:"release"`
	AppVersion string `env:"OHM9996_APP_VERSION" envDefault:""`
	Port       int    `env:"OHM9996_SERVER_PORT" envDefault:"8080"`
	timeoutFields
}

// DefaultConfig returns server configuration with defaults.
func DefaultConfig() *Config {
	timeout := time.Duration(defaultTimeoutSec) * time.Second
	idle := time.Duration(defaultIdleTimeoutSec) * time.Second
	return &Config{
		Host:            "",
		Port:            defaultServerPort,
		Mode:            gin.ReleaseMode,
		AppVersion:      "",
		ReadTimeout:     timeout,
		WriteTimeout:    timeout,
		IdleTimeout:     idle,
		ShutdownTimeout: timeout,
	}
}
