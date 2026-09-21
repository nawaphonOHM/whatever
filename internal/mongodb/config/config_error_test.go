package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupValidBaseEnv sets up valid connection and credential environment variables.
func setupValidBaseEnv(t *testing.T) {
	t.Setenv("OHM9996_MONGODB_URI", "mongodb://test:27017")
	t.Setenv("OHM9996_MONGODB_DATABASE", "testdb")
	t.Setenv("OHM9996_MONGODB_USERNAME", "testuser")
	t.Setenv("OHM9996_MONGODB_PASSWORD", "testpass")
}

// TestConfig_Load_InvalidTypeConversion tests type conversion failure on load.
func TestConfig_Load_InvalidTypeConversion(t *testing.T) {
	setupValidBaseEnv(t)
	t.Setenv("OHM9996_MONGODB_MAX_POOL_SIZE", "invalid-uint")

	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to load mongodb config")
}

// TestConfig_Load_InvalidValidation tests validation failure on load.
func TestConfig_Load_InvalidValidation(t *testing.T) {
	setupValidBaseEnv(t)
	t.Setenv("OHM9996_MONGODB_CONNECT_TIMEOUT", "-10s")

	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "invalid mongodb config")
}
