package config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEnvOverrides(t *testing.T) {
	t.Helper()
	t.Setenv(testEnvOutput, testStderr)
	t.Setenv(testEnvLevel, testDebug)
	t.Setenv(testEnvFormat, testJSON)
	t.Setenv(testEnvAddSource, "true")
	t.Setenv(testEnvDisableTraceCorrelation, "true")
}

func verifyEnvOverridesFields(t *testing.T, cfg *Config) {
	t.Helper()
	assert.Equal(t, testStderr, cfg.Output)
	assert.Equal(t, testDebug, cfg.Level)
	assert.Equal(t, testJSON, cfg.Format)
	assert.True(t, cfg.AddSource)
	assert.True(t, cfg.DisableTraceCorrelation)
}

func verifyEnvOverridesLevel(t *testing.T, cfg *Config) {
	t.Helper()
	lvl, err := cfg.ParsedLevel()
	require.NoError(t, err)
	assert.Equal(t, LevelDebug, lvl)
	slvl, err := cfg.ParsedSlogLevel()
	require.NoError(t, err)
	assert.Equal(t, SlogLevelDebug, slvl)
}

func verifyEnvOverridesOutput(t *testing.T, cfg *Config) {
	t.Helper()
	fmtEnum, err := cfg.ParsedFormat()
	require.NoError(t, err)
	assert.Equal(t, FormatJSON, fmtEnum)
	out, err := cfg.ResolveOutput()
	require.NoError(t, err)
	assert.Equal(t, os.Stderr, out)
}

func verifyEnvOverridesParsed(t *testing.T, cfg *Config) {
	t.Helper()
	verifyEnvOverridesLevel(t, cfg)
	verifyEnvOverridesOutput(t, cfg)
}

func TestLoad_EnvOverrides(t *testing.T) {
	setEnvOverrides(t)
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyEnvOverridesFields(t, cfg)
	verifyEnvOverridesParsed(t, cfg)
}

func createDotEnvFile(t *testing.T) string {
	t.Helper()
	envPath := filepath.Join(t.TempDir(), ".env.logging")
	content := []byte("OHM9996_LOGGING_OUTPUT=discard\n" +
		"OHM9996_LOGGING_LEVEL=warn\n" +
		"OHM9996_LOGGING_FORMAT=json\n" +
		"OHM9996_LOGGING_ADD_SOURCE=true\n" +
		"OHM9996_LOGGING_DISABLE_TRACE_CORRELATION=false\n")
	require.NoError(t, os.WriteFile(envPath, content, testDotEnvPerm))
	return envPath
}

func verifyDotEnvConfig(t *testing.T, cfg *Config) {
	t.Helper()
	assert.Equal(t, testDiscard, cfg.Output)
	assert.Equal(t, "warn", cfg.Level)
	assert.Equal(t, testJSON, cfg.Format)
	assert.True(t, cfg.AddSource)
	assert.False(t, cfg.DisableTraceCorrelation)

	out, err := cfg.ResolveOutput()
	require.NoError(t, err)
	assert.Equal(t, io.Discard, out)
}

func TestLoad_DotEnvFile(t *testing.T) {
	clearEnv(t)
	cfg, err := Load(createDotEnvFile(t))
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyDotEnvConfig(t, cfg)
}
