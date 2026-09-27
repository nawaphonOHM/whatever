package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupCustomEnv sets custom values for OpenTelemetry environment variables.
func setupCustomEnv(t *testing.T) {
	t.Setenv(testEnvEnabled, "true")
	t.Setenv(testEnvServiceName, testCustomServiceName)
	t.Setenv(testEnvEndpoint, testCustomEndpoint)
	t.Setenv(testEnvProtocol, ProtocolHTTPProtobuf)
	t.Setenv(testEnvInsecure, "false")
	t.Setenv(testEnvSampleRate, "0.75")
	t.Setenv(testEnvSkipPaths, "/health,/ready,/metrics")
}

// verifyLoadedEnvConfig verifies fields of custom env loaded configuration.
func verifyLoadedEnvConfig(t *testing.T, cfg *Config) {
	assert.True(t, cfg.Enabled)
	assert.Equal(t, testCustomServiceName, cfg.ServiceName)
	assert.Equal(t, testCustomEndpoint, cfg.Endpoint)
	assert.Equal(t, ProtocolHTTPProtobuf, cfg.Protocol)
	assert.False(t, cfg.Insecure)
	assert.Equal(t, testSampleRateHigh, cfg.SampleRate)
	assert.Equal(t, []string{testHealthPath, testReadyPath, testMetricsPath}, cfg.SkipPaths)
	assert.False(t, cfg.IsGRPC())
	assert.True(t, cfg.IsHTTP())
}

// TestLoad_Defaults verifies Load parses defaults when env vars are unset or empty.
func TestLoad_Defaults(t *testing.T) {
	clearEnv(t)

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyDefaultConfig(t, cfg)
}

// TestLoad_EnvOverrides verifies overriding default configuration with environment variables.
func TestLoad_EnvOverrides(t *testing.T) {
	setupCustomEnv(t)

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyLoadedEnvConfig(t, cfg)
}

// TestLoad_Disabled verifies loading configuration with OpenTelemetry disabled.
func TestLoad_Disabled(t *testing.T) {
	t.Setenv(testEnvEnabled, "false")
	t.Setenv(testEnvServiceName, "")
	t.Setenv(testEnvEndpoint, "")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.False(t, cfg.Enabled)
}

// TestLoad_InvalidTypeConversion verifies parsing fails on invalid types.
func TestLoad_InvalidTypeConversion(t *testing.T) {
	clearEnv(t)
	t.Setenv(testEnvSampleRate, "not-a-number")
	cfg, err := Load()
	require.Error(t, err)
	assert.Nil(t, cfg)
}

// TestLoad_InvalidValidation verifies Load returns error on invalid config values.
func TestLoad_InvalidValidation(t *testing.T) {
	clearEnv(t)
	t.Setenv(testEnvProtocol, "invalid-proto")
	cfg, err := Load()
	require.Error(t, err)
	assert.Nil(t, cfg)
}
