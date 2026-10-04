package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupCustomEnv sets custom values for all mongodb environment variables.
func setupCustomEnv(t *testing.T) {
	envVars := map[string]string{
		testEnvHost:                   "db.internal",
		testEnvPort:                   "27018",
		testEnvProtocol:               "mongodb",
		testEnvDatabase:               "orders_db",
		testEnvUsername:               "myuser",
		testEnvPassword:               "mypassword",
		testEnvAuthSource:             "admin",
		testEnvAppName:                "order-service",
		testEnvUUIDRepresentation:     "standard",
		testEnvConnectTimeout:         "15s",
		testEnvServerSelectionTimeout: "3s",
		testEnvSocketTimeout:          "20s",
		testEnvMaxConnIdleTime:        "5m",
		testEnvMaxPoolSize:            "200",
		testEnvMinPoolSize:            "10",
	}
	for k, v := range envVars {
		t.Setenv(k, v)
	}
}

func verifyLoadedBaseFields(t *testing.T, cfg *Config) {
	assert.Equal(t, "db.internal", cfg.Host)
	assert.Equal(t, testCustomPort, cfg.Port)
	assert.Equal(t, "mongodb", cfg.Protocol)
	assert.Equal(t, "orders_db", cfg.Database)
	assert.Equal(t, "myuser", cfg.Username)
	assert.Equal(t, "mypassword", cfg.Password)
	assert.Equal(t, "admin", cfg.AuthSource)
	assert.Equal(t, "order-service", cfg.AppName)
	assert.Equal(t, "standard", cfg.UUIDRepresentation)
}

func verifyLoadedRuntimeFields(t *testing.T, cfg *Config) {
	assert.Equal(t, 15*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, 3*time.Second, cfg.ServerSelectionTimeout)
	assert.Equal(t, 20*time.Second, cfg.SocketTimeout)
	assert.Equal(t, 5*time.Minute, cfg.MaxConnIdleTime)
	assert.Equal(t, uint64(testCustomMaxPoolVal), cfg.MaxPoolSize)
	assert.Equal(t, uint64(testCustomMinPoolVal), cfg.MinPoolSize)
}

// TestConfig_Load_EnvOverrides tests loading config with environment variables.
func TestConfig_Load_EnvOverrides(t *testing.T) {
	setupCustomEnv(t)

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyLoadedBaseFields(t, cfg)
	verifyLoadedRuntimeFields(t, cfg)
	assert.NoError(t, cfg.Validate())
}

func verifyMinimalConfig(t *testing.T, cfg *Config) {
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, defaultPort, cfg.Port)
	assert.Equal(t, ProtocolMongoDB, cfg.Protocol)
	assert.Empty(t, cfg.Username)
	assert.Empty(t, cfg.Password)
	assert.Equal(t, testDefaultTimeoutSec*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, uint64(testDefaultMaxPoolVal), cfg.MaxPoolSize)
}

// TestConfig_Load_MinimalEnv tests loading config with only mandatory host set.
func TestConfig_Load_MinimalEnv(t *testing.T) {
	setupEmptyEnv(t)
	t.Setenv(testEnvHost, "localhost")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyMinimalConfig(t, cfg)
	assert.NoError(t, cfg.Validate())
}
