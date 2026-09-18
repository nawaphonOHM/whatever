// Package config provides unit tests for environment configuration parsing.
package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envPort        = "PORT"
	envHost        = "HOST"
	envDebug       = "DEBUG"
	envTimeout     = "TIMEOUT"
	envDBURL       = "DATABASE_URL"
	envOptVal      = "OPTIONAL_VAL"
	filePerm0600   = 0o600
	defaultOptVal  = "default_opt"
	customOptVal   = "custom_opt"
	postgresURL    = "postgres://user:pass@localhost:5432/db"
	sqliteURL      = "sqlite://test.db"
	localHostStr   = "localhost"
	defaultPortVal = 8080
	overridePort   = 9090
	dotEnvPort     = 3000
	precedencePort = 4000
	missingPort    = 5000
)

// TestConfig represents a test configuration structure with struct tags.
type TestConfig struct {
	Host        string        `env:"HOST" envDefault:"localhost"`
	DatabaseURL string        `env:"DATABASE_URL"`
	OptionalVal string        `env:"OPTIONAL_VAL" envDefault:"default_opt"`
	Timeout     time.Duration `env:"TIMEOUT" envDefault:"5s"`
	Port        int           `env:"PORT" envDefault:"8080"`
	Debug       bool          `env:"DEBUG" envDefault:"false"`
}

// unsetEnv safely unsets environment variables and registers cleanup.
func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		orig, had := os.LookupEnv(key)
		require.NoError(t, os.Unsetenv(key))
		t.Cleanup(func() {
			if had {
				assert.NoError(t, os.Setenv(key, orig))
			} else {
				assert.NoError(t, os.Unsetenv(key))
			}
		})
	}
}

// verifyDefaultConfig validates default values of parsed config.
func verifyDefaultConfig(t *testing.T, cfg *TestConfig) {
	assert.Equal(t, defaultPortVal, cfg.Port)
	assert.Equal(t, localHostStr, cfg.Host)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 5*time.Second, cfg.Timeout)
	assert.Equal(t, "", cfg.DatabaseURL)
	assert.Equal(t, defaultOptVal, cfg.OptionalVal)
}

// TestLoad_Defaults tests parsing config using only default tags.
func TestLoad_Defaults(t *testing.T) {
	for _, k := range []string{
		envPort, envHost, envDebug, envTimeout, envDBURL, envOptVal,
	} {
		t.Setenv(k, "")
	}

	cfg, err := Load[TestConfig]()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyDefaultConfig(t, cfg)
}

// verifyOverrideConfig validates overridden configuration values.
func verifyOverrideConfig(t *testing.T, cfg *TestConfig) {
	assert.Equal(t, overridePort, cfg.Port)
	assert.Equal(t, "0.0.0.0", cfg.Host)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 10*time.Second, cfg.Timeout)
	assert.Equal(t, postgresURL, cfg.DatabaseURL)
	assert.Equal(t, customOptVal, cfg.OptionalVal)
}

// TestLoad_EnvOverrides tests overriding struct defaults with system env.
func TestLoad_EnvOverrides(t *testing.T) {
	t.Setenv(envPort, "9090")
	t.Setenv(envHost, "0.0.0.0")
	t.Setenv(envDebug, "true")
	t.Setenv(envTimeout, "10s")
	t.Setenv(envDBURL, postgresURL)
	t.Setenv(envOptVal, customOptVal)

	cfg, err := Load[TestConfig]()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyOverrideConfig(t, cfg)
}

// TestLoad_InvalidTypeConversion tests error on invalid type conversions.
func TestLoad_InvalidTypeConversion(t *testing.T) {
	t.Setenv(envPort, "not-a-number")
	cfg, err := Load[TestConfig]()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
