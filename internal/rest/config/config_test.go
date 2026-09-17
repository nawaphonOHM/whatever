package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Port        int           `env:"PORT" envDefault:"8080"`
	Host        string        `env:"HOST" envDefault:"localhost"`
	Debug       bool          `env:"DEBUG" envDefault:"false"`
	Timeout     time.Duration `env:"TIMEOUT" envDefault:"5s"`
	DatabaseURL string        `env:"DATABASE_URL"`
	OptionalVal string        `env:"OPTIONAL_VAL" envDefault:"default_opt"`
}

func TestLoad_Defaults(t *testing.T) {
	// Clean environment
	t.Setenv("PORT", "")
	t.Setenv("HOST", "")
	t.Setenv("DEBUG", "")
	t.Setenv("TIMEOUT", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("OPTIONAL_VAL", "")

	cfg, err := Load[TestConfig]()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "localhost", cfg.Host)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 5*time.Second, cfg.Timeout)
	assert.Equal(t, "", cfg.DatabaseURL)
	assert.Equal(t, "default_opt", cfg.OptionalVal)
}

func TestLoad_EnvOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("HOST", "0.0.0.0")
	t.Setenv("DEBUG", "true")
	t.Setenv("TIMEOUT", "10s")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("OPTIONAL_VAL", "custom_opt")

	cfg, err := Load[TestConfig]()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "0.0.0.0", cfg.Host)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 10*time.Second, cfg.Timeout)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	assert.Equal(t, "custom_opt", cfg.OptionalVal)
}

func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		orig, had := os.LookupEnv(key)
		require.NoError(t, os.Unsetenv(key))
		t.Cleanup(func() {
			if had {
				_ = os.Setenv(key, orig)
			} else {
				_ = os.Unsetenv(key)
			}
		})
	}
}

func TestLoad_DotEnvFile(t *testing.T) {
	// Create temporary .env file
	tmpDir := t.TempDir()
	envFilePath := filepath.Join(tmpDir, ".env.test")
	envContent := `PORT=3000
HOST=127.0.0.1
DEBUG=true
TIMEOUT=15s
DATABASE_URL=sqlite://test.db
`
	err := os.WriteFile(envFilePath, []byte(envContent), 0600)
	require.NoError(t, err)

	// Fully unset so godotenv can populate from the file (empty string still blocks it).
	unsetEnv(t, "PORT", "HOST", "DEBUG", "TIMEOUT", "DATABASE_URL")

	cfg, err := Load[TestConfig](envFilePath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, 3000, cfg.Port)
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 15*time.Second, cfg.Timeout)
	assert.Equal(t, "sqlite://test.db", cfg.DatabaseURL)
}

func TestLoad_PrecedenceOrder(t *testing.T) {
	// Create temporary .env file
	tmpDir := t.TempDir()
	envFilePath := filepath.Join(tmpDir, ".env.precedence")
	envContent := `PORT=3000
HOST=from-env-file
`
	err := os.WriteFile(envFilePath, []byte(envContent), 0600)
	require.NoError(t, err)

	// System env overrides .env file
	t.Setenv("PORT", "4000")
	// HOST must be truly unset so the .env value is applied.
	unsetEnv(t, "HOST", "OPTIONAL_VAL")

	cfg, err := Load[TestConfig](envFilePath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// System env wins for PORT
	assert.Equal(t, 4000, cfg.Port)
	// .env file wins for HOST
	assert.Equal(t, "from-env-file", cfg.Host)
	// Default wins for OPTIONAL_VAL
	assert.Equal(t, "default_opt", cfg.OptionalVal)
}

func TestLoad_MissingDotEnvFile(t *testing.T) {
	t.Setenv("PORT", "5000")
	// Non-existent file should gracefully fall back to system env / defaults
	cfg, err := Load[TestConfig]("non_existent_file.env")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, 5000, cfg.Port)
}

func TestLoad_InvalidTypeConversion(t *testing.T) {
	t.Setenv("PORT", "not-a-number")
	cfg, err := Load[TestConfig]()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
