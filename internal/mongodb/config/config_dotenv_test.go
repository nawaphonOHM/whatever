package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testFilePerm      = 0o600
	testDotEnvPort    = 27017
	testDotEnvMaxPool = 150
)

// createDotEnvFile writes a temporary dotenv configuration file.
func createDotEnvFile(t *testing.T) string {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env.mongodb")
	content := []byte(`
OHM9996_MONGODB_HOST=dotenv-host
OHM9996_MONGODB_PORT=27017
OHM9996_MONGODB_DATABASE=dotenv_db
OHM9996_MONGODB_USERNAME=dotenv_user
OHM9996_MONGODB_PASSWORD=dotenv_pass
OHM9996_MONGODB_AUTH_SOURCE=admin
OHM9996_MONGODB_APP_NAME=dotenv-app
OHM9996_MONGODB_PROTOCOL=mongodb
OHM9996_MONGODB_UUID_REPRESENTATION=standard
OHM9996_MONGODB_MAX_POOL_SIZE=150
`)
	err := os.WriteFile(envPath, content, testFilePerm)
	require.NoError(t, err)
	return envPath
}

// verifyDotEnvLoadedConfig asserts loaded configuration fields match dotenv contents.
func verifyDotEnvLoadedConfig(t *testing.T, cfg *Config) {
	assert.Equal(t, "dotenv-host", cfg.Host)
	assert.Equal(t, testDotEnvPort, cfg.Port)
	assert.Equal(t, "dotenv_db", cfg.Database)
	assert.Equal(t, "dotenv_user", cfg.Username)
	assert.Equal(t, "dotenv_pass", cfg.Password)
	assert.Equal(t, "admin", cfg.AuthSource)
	assert.Equal(t, "dotenv-app", cfg.AppName)
	assert.Equal(t, "mongodb", cfg.Protocol)
	assert.Equal(t, "standard", cfg.UUIDRepresentation)
	assert.Equal(t, uint64(testDotEnvMaxPool), cfg.MaxPoolSize)
}

// TestConfig_Load_DotEnvFile tests loading configuration from a dotenv file.
func TestConfig_Load_DotEnvFile(t *testing.T) {
	envPath := createDotEnvFile(t)
	cfg, err := LoadConfig(envPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyDotEnvLoadedConfig(t, cfg)
	assert.NoError(t, cfg.Validate())
}
