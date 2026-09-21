package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDefaultMaxPool  = uint64(100)
	testDefaultMinPool  = uint64(5)
	testEnvURI          = "OHM9996_MONGODB_URI"
	testEnvDatabase     = "OHM9996_MONGODB_DATABASE"
	testEnvUsername     = "OHM9996_MONGODB_USERNAME"
	testEnvPassword     = "OHM9996_MONGODB_PASSWORD"
	testEnvConnect      = "OHM9996_MONGODB_CONNECT_TIMEOUT"
	testEnvServerSelect = "OHM9996_MONGODB_SERVER_SELECTION_TIMEOUT"
	testEnvSocket       = "OHM9996_MONGODB_SOCKET_TIMEOUT"
	testEnvMaxPool      = "OHM9996_MONGODB_MAX_POOL_SIZE"
	testEnvMinPool      = "OHM9996_MONGODB_MIN_POOL_SIZE"
	testEnvIdleTime     = "OHM9996_MONGODB_MAX_CONN_IDLE_TIME"
	testEnvAppName      = "OHM9996_MONGODB_APP_NAME"
)

// setupEmptyEnv clears all mongodb environment variable overrides.
func setupEmptyEnv(t *testing.T) {
	for _, envVar := range []string{
		testEnvURI, testEnvDatabase, testEnvUsername, testEnvPassword,
		testEnvConnect, testEnvServerSelect, testEnvSocket,
		testEnvMaxPool, testEnvMinPool, testEnvIdleTime, testEnvAppName,
	} {
		t.Setenv(envVar, "")
	}
}

// setupCustomEnv sets custom values for all mongodb environment variables.
func setupCustomEnv(t *testing.T) {
	envVars := map[string]string{
		testEnvURI:          "mongodb://remote:27018",
		testEnvDatabase:     "analytics",
		testEnvUsername:     "myuser",
		testEnvPassword:     "mypassword",
		testEnvConnect:      "20s",
		testEnvServerSelect: "15s",
		testEnvSocket:       "30s",
		testEnvMaxPool:      "200",
		testEnvMinPool:      "10",
		testEnvIdleTime:     "15m",
		testEnvAppName:      "my-service",
	}
	for k, v := range envVars {
		t.Setenv(k, v)
	}
}

// verifyDefaultBaseFields asserts default base connection fields are empty.
func verifyDefaultBaseFields(t *testing.T, cfg *Config) {
	assert.Empty(t, cfg.URI)
	assert.Empty(t, cfg.Database)
	assert.Empty(t, cfg.Username)
	assert.Empty(t, cfg.Password)
	assert.Empty(t, cfg.AppName)
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
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	verifyDefaultBaseFields(t, cfg)
	verifyDefaultTimeouts(t, cfg)
	assert.Equal(t, testDefaultMaxPool, cfg.MaxPoolSize)
	assert.Equal(t, testDefaultMinPool, cfg.MinPoolSize)
	assert.Error(t, cfg.Validate())
}

// verifyLoadedEnvFields asserts loaded configuration fields match custom env.
func verifyLoadedEnvFields(t *testing.T, cfg *Config) {
	assert.Equal(t, "mongodb://remote:27018", cfg.URI)
	assert.Equal(t, "analytics", cfg.Database)
	assert.Equal(t, "myuser", cfg.Username)
	assert.Equal(t, "mypassword", cfg.Password)
	assert.Equal(t, 20*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, "my-service", cfg.AppName)
}

// TestConfig_Load_EnvOverrides tests loading config with environment variables.
func TestConfig_Load_EnvOverrides(t *testing.T) {
	setupCustomEnv(t)

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyLoadedEnvFields(t, cfg)
	assert.NoError(t, cfg.Validate())
}
