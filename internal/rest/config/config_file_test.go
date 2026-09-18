// Package config provides unit tests for file-based configuration loading.
package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeTempEnvFile writes temporary env content and returns path.
func writeTempEnvFile(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(p, []byte(content), filePerm0600))
	return p
}

// verifyDotEnvLoadedConfig validates config loaded from test .env file.
func verifyDotEnvLoadedConfig(t *testing.T, cfg *TestConfig) {
	assert.Equal(t, dotEnvPort, cfg.Port)
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 15*time.Second, cfg.Timeout)
	assert.Equal(t, sqliteURL, cfg.DatabaseURL)
}

// TestLoad_DotEnvFile tests loading configuration from custom .env file.
func TestLoad_DotEnvFile(t *testing.T) {
	content := "PORT=3000\nHOST=127.0.0.1\nDEBUG=true\n" +
		"TIMEOUT=15s\nDATABASE_URL=sqlite://test.db\n"
	envFile := writeTempEnvFile(t, ".env.test", content)
	unsetEnv(t, envPort, envHost, envDebug, envTimeout, envDBURL)

	cfg, err := Load[TestConfig](envFile)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyDotEnvLoadedConfig(t, cfg)
}

// TestLoad_PrecedenceOrder tests system env overriding dotEnv file values.
func TestLoad_PrecedenceOrder(t *testing.T) {
	envContent := "PORT=3000\nHOST=from-env-file\n"
	envFile := writeTempEnvFile(t, ".env.precedence", envContent)
	t.Setenv(envPort, "4000")
	unsetEnv(t, envHost, envOptVal)

	cfg, err := Load[TestConfig](envFile)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, precedencePort, cfg.Port)
	assert.Equal(t, "from-env-file", cfg.Host)
	assert.Equal(t, defaultOptVal, cfg.OptionalVal)
}

// TestLoad_MissingDotEnvFile tests fallback when specified file does not exist.
func TestLoad_MissingDotEnvFile(t *testing.T) {
	t.Setenv(envPort, "5000")
	cfg, err := Load[TestConfig]("non_existent_file.env")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, missingPort, cfg.Port)
}
