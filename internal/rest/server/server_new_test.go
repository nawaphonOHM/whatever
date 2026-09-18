package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew_WithNilConfig uses package defaults.
func TestNew_WithNilConfig(t *testing.T) {
	// Act
	srv := New(nil)

	// Assert
	require.NotNil(t, srv)
	require.NotNil(t, srv.Engine)
	require.NotNil(t, srv.Config)
	assert.Equal(t, defaultServerPort, srv.Config.Port)
	assert.Equal(t, gin.ReleaseMode, srv.Config.Mode)
	assert.Equal(t, fmt.Sprintf(":%d", defaultServerPort),
		srv.httpServer.Addr)
}

// customTestConfig returns a non-default server config.
func customTestConfig() *Config {
	return &Config{
		Host:         testHost,
		Port:         customPort,
		Mode:         gin.TestMode,
		ReadTimeout:  customReadSec * time.Second,
		WriteTimeout: customReadSec * time.Second,
		IdleTimeout:  customIdleSec * time.Second,
	}
}

// TestNew_WithCustomConfig preserves caller values.
func TestNew_WithCustomConfig(t *testing.T) {
	// Arrange
	cfg := customTestConfig()

	// Act
	srv := New(cfg)

	// Assert
	require.NotNil(t, srv)
	assert.Equal(t, testHost, srv.Config.Host)
	assert.Equal(t, customPort, srv.Config.Port)
	assert.Equal(t, gin.TestMode, srv.Config.Mode)
	assert.Equal(t,
		fmt.Sprintf("%s:%d", testHost, customPort),
		srv.httpServer.Addr,
	)
}

// TestDefaultConfig_Values checks recommended defaults.
func TestDefaultConfig_Values(t *testing.T) {
	// Act
	cfg := DefaultConfig()

	// Assert
	require.NotNil(t, cfg)
	assert.Equal(t, defaultServerPort, cfg.Port)
	assert.Equal(t, gin.ReleaseMode, cfg.Mode)
	assert.Equal(t, defaultTimeoutSec*time.Second, cfg.ReadTimeout)
	assert.Equal(t, defaultTimeoutSec*time.Second, cfg.WriteTimeout)
	assert.Equal(t, defaultIdleTimeoutSec*time.Second, cfg.IdleTimeout)
	assert.Equal(t, defaultTimeoutSec*time.Second, cfg.ShutdownTimeout)
}
