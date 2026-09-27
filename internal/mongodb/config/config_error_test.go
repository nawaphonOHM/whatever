package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupValidBaseEnv sets up valid connection and credential environment variables.
func setupValidBaseEnv(t *testing.T) {
	t.Setenv("OHM9996_MONGODB_HOST", "localhost")
	t.Setenv("OHM9996_MONGODB_PORT", "27017")
	t.Setenv("OHM9996_MONGODB_USERNAME", "testuser")
	t.Setenv("OHM9996_MONGODB_PASSWORD", "testpass")
}

// TestConfig_Load_InvalidTypeConversion tests type conversion failure on load.
func TestConfig_Load_InvalidTypeConversion(t *testing.T) {
	setupValidBaseEnv(t)
	t.Setenv("OHM9996_MONGODB_PORT", "invalid-port")

	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to load mongodb config")
}

// TestConfig_Load_InvalidValidation tests validation failure on load.
func TestConfig_Load_InvalidValidation(t *testing.T) {
	setupValidBaseEnv(t)
	t.Setenv("OHM9996_MONGODB_PROTOCOL", "invalid-proto")

	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "invalid mongodb config")
}
