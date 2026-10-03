package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func verifyDefaultLevel(t *testing.T, cfg *Config) {
	t.Helper()
	lvl, err := cfg.ParsedLevel()
	require.NoError(t, err)
	assert.Equal(t, LevelInfo, lvl)

	slvl, err := cfg.ParsedSlogLevel()
	require.NoError(t, err)
	assert.Equal(t, SlogLevelInfo, slvl)
}

func verifyDefaultOutput(t *testing.T, cfg *Config) {
	t.Helper()
	fmtEnum, err := cfg.ParsedFormat()
	require.NoError(t, err)
	assert.Equal(t, FormatText, fmtEnum)

	out, err := cfg.ResolveOutput()
	require.NoError(t, err)
	assert.Equal(t, os.Stdout, out)
}

// verifyDefaultConfig asserts default configuration values.
func verifyDefaultConfig(t *testing.T, cfg *Config) {
	t.Helper()
	assert.Equal(t, DefaultOutput, cfg.Output)
	assert.Equal(t, DefaultLevel, cfg.Level)
	assert.Equal(t, DefaultFormat, cfg.Format)
	assert.Equal(t, DefaultAddSource, cfg.AddSource)
	assert.Equal(t, DefaultDisableTraceCorrelation, cfg.DisableTraceCorrelation)
	verifyDefaultLevel(t, cfg)
	verifyDefaultOutput(t, cfg)
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	verifyDefaultConfig(t, cfg)
	assert.NoError(t, cfg.Validate())
}

func TestSetDefaults(t *testing.T) {
	cfg := &Config{
		Output:                  testStderr,
		Level:                   "debug",
		Format:                  testJSON,
		AddSource:               true,
		DisableTraceCorrelation: true,
	}
	cfg.SetDefaults()
	verifyDefaultConfig(t, cfg)
}

func TestLoad_Defaults(t *testing.T) {
	clearEnv(t)

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyDefaultConfig(t, cfg)
}
