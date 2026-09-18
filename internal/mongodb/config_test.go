package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDefaultMaxPool = uint64(100)
	testDefaultMinPool = uint64(5)
)

// setupEmptyEnv clears all mongodb environment variable overrides.
func setupEmptyEnv(t *testing.T) {
	// Clear all connection string and database env vars.
	t.Setenv("OHM9969_MONGODB_URI", "")
	t.Setenv("OHM9969_MONGODB_DATABASE", "")
	t.Setenv("OHM9969_MONGODB_CONNECT_TIMEOUT", "")
	t.Setenv("OHM9969_MONGODB_SERVER_SELECTION_TIMEOUT", "")
	t.Setenv("OHM9969_MONGODB_SOCKET_TIMEOUT", "")
	t.Setenv("OHM9969_MONGODB_MAX_POOL_SIZE", "")
	t.Setenv("OHM9969_MONGODB_MIN_POOL_SIZE", "")
	t.Setenv("OHM9969_MONGODB_MAX_CONN_IDLE_TIME", "")
	t.Setenv("OHM9969_MONGODB_APP_NAME", "")
}

// setupCustomEnv sets custom values for all mongodb environment variables.
func setupCustomEnv(t *testing.T) {
	// Set custom connection parameters.
	t.Setenv("OHM9969_MONGODB_URI", "mongodb://user:pass@remote:27018")
	t.Setenv("OHM9969_MONGODB_DATABASE", "analytics")
	t.Setenv("OHM9969_MONGODB_CONNECT_TIMEOUT", "20s")
	t.Setenv("OHM9969_MONGODB_SERVER_SELECTION_TIMEOUT", "15s")
	t.Setenv("OHM9969_MONGODB_SOCKET_TIMEOUT", "30s")
	t.Setenv("OHM9969_MONGODB_MAX_POOL_SIZE", "200")
	t.Setenv("OHM9969_MONGODB_MIN_POOL_SIZE", "10")
	t.Setenv("OHM9969_MONGODB_MAX_CONN_IDLE_TIME", "15m")
	t.Setenv("OHM9969_MONGODB_APP_NAME", "my-service")
}

// verifyDefaultTimeouts checks default timeout durations.
func verifyDefaultTimeouts(t *testing.T, cfg *Config) {
	assert.Equal(t, 10*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, 5*time.Second, cfg.ServerSelectionTimeout)
	assert.Equal(t, 10*time.Second, cfg.SocketTimeout)
	assert.Equal(t, 10*time.Minute, cfg.MaxConnIdleTime)
}

// TestDefaultConfig verifies default values created by DefaultConfig.
func TestDefaultConfig(t *testing.T) {
	// Load default configuration object.
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	assert.Equal(t, DefaultURI, cfg.URI)
	assert.Empty(t, cfg.Database)
	assert.Empty(t, cfg.AppName)

	// Verify timeouts and connection pool bounds.
	verifyDefaultTimeouts(t, cfg)
	assert.Equal(t, testDefaultMaxPool, cfg.MaxPoolSize)
	assert.Equal(t, testDefaultMinPool, cfg.MinPoolSize)
	assert.NoError(t, cfg.Validate())
}

// TestConfig_Load_Defaults tests loading default config without env overrides.
func TestConfig_Load_Defaults(t *testing.T) {
	// Ensure environment is clean before loading.
	setupEmptyEnv(t)

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, DefaultURI, cfg.URI)
	assert.Equal(t, 10*time.Second, cfg.ConnectTimeout)
	assert.NoError(t, cfg.Validate())
}

// TestConfig_Load_EnvOverrides tests loading config with environment variables.
func TestConfig_Load_EnvOverrides(t *testing.T) {
	// Populate custom environment variables.
	setupCustomEnv(t)

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "mongodb://user:pass@remote:27018", cfg.URI)
	assert.Equal(t, "analytics", cfg.Database)
	assert.Equal(t, 20*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, "my-service", cfg.AppName)
	assert.NoError(t, cfg.Validate())
}
