package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testEnvHost               = "OHM9996_MONGODB_HOST"
	testEnvUsername           = "OHM9996_MONGODB_USERNAME"
	testEnvPassword           = "OHM9996_MONGODB_PASSWORD"
	testEnvPort               = "OHM9996_MONGODB_PORT"
	testEnvProtocol           = "OHM9996_MONGODB_PROTOCOL"
	testEnvUUIDRepresentation = "OHM9996_MONGODB_UUID_REPRESENTATION"
	testCustomPort            = 27018
)

// setupEmptyEnv clears all mongodb environment variable overrides.
func setupEmptyEnv(t *testing.T) {
	for _, envVar := range []string{
		testEnvHost, testEnvUsername, testEnvPassword,
		testEnvPort, testEnvProtocol, testEnvUUIDRepresentation,
	} {
		t.Setenv(envVar, "")
	}
}

// setupCustomEnv sets custom values for all mongodb environment variables.
func setupCustomEnv(t *testing.T) {
	envVars := map[string]string{
		testEnvHost:               "db.internal",
		testEnvUsername:           "myuser",
		testEnvPassword:           "mypassword",
		testEnvPort:               "27018",
		testEnvProtocol:           "mongodb",
		testEnvUUIDRepresentation: "standard",
	}
	for k, v := range envVars {
		t.Setenv(k, v)
	}
}

// verifyDefaultFields asserts default connection fields.
func verifyDefaultFields(t *testing.T, cfg *Config) {
	assert.Empty(t, cfg.Host)
	assert.Empty(t, cfg.Username)
	assert.Empty(t, cfg.Password)
	assert.Zero(t, cfg.Port)
	assert.Equal(t, ProtocolMongoDB, cfg.Protocol)
	assert.Equal(t, UUIDRepresentationUnspecified, cfg.UUIDRepresentation)
}

// TestDefaultConfig verifies default values created by DefaultConfig.
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	verifyDefaultFields(t, cfg)
	assert.Error(t, cfg.Validate())
}

// verifyLoadedEnvFields asserts loaded configuration fields match custom env.
func verifyLoadedEnvFields(t *testing.T, cfg *Config) {
	assert.Equal(t, "db.internal", cfg.Host)
	assert.Equal(t, "myuser", cfg.Username)
	assert.Equal(t, "mypassword", cfg.Password)
	assert.Equal(t, testCustomPort, cfg.Port)
	assert.Equal(t, "mongodb", cfg.Protocol)
	assert.Equal(t, "standard", cfg.UUIDRepresentation)
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
