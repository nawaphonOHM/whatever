package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig_LoadFromEnv loads server config from OHM9996_ keys.
func TestConfig_LoadFromEnv(t *testing.T) {
	// Arrange
	t.Setenv(envServerPort, "9090")
	t.Setenv(envGinMode, gin.TestMode)

	// Act
	cfg, err := intcfg.Load[Config]()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, customPort, cfg.Port)
	assert.Equal(t, gin.TestMode, cfg.Mode)
}
