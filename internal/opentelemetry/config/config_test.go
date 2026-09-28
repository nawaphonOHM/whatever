package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testEnvEnabled     = "OHM9996_OTEL_ENABLED"
	testEnvServiceName = "OHM9996_OTEL_SERVICE_NAME"
	testEnvEndpoint    = "OHM9996_OTEL_EXPORTER_OTLP_ENDPOINT"
	testEnvProtocol    = "OHM9996_OTEL_EXPORTER_OTLP_PROTOCOL"
	testEnvInsecure    = "OHM9996_OTEL_INSECURE"
	testEnvSampleRate  = "OHM9996_OTEL_SAMPLE_RATE"
	testEnvSkipPaths   = "OHM9996_OTEL_SKIP_PATHS"

	testServiceName       = "whatever-service"
	testDefaultEndpoint   = "localhost:4317"
	testCustomServiceName = "billing-service"
	testCustomEndpoint    = "collector.internal:4318"
	testHealthPath        = "/health"
	testReadyPath         = "/ready"
	testMetricsPath       = "/metrics"
	testSampleRateQuarter = 0.25
	testSampleRateHalf    = 0.5
	testSampleRateHigh    = 0.75
)

// unsetEnv unsets an environment variable and restores it during test cleanup.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	if val, ok := os.LookupEnv(key); ok {
		require.NoError(t, os.Unsetenv(key))
		t.Cleanup(func() {
			require.NoError(t, os.Setenv(key, val))
		})
	}
}

// clearEnv clears all OpenTelemetry environment variables.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		testEnvEnabled,
		testEnvServiceName,
		testEnvEndpoint,
		testEnvProtocol,
		testEnvInsecure,
		testEnvSampleRate,
		testEnvSkipPaths,
	} {
		unsetEnv(t, k)
	}
}

// verifyDefaultConfig asserts default configuration values.
func verifyDefaultConfig(t *testing.T, cfg *Config) {
	t.Helper()
	assert.True(t, cfg.Enabled)
	assert.Equal(t, testServiceName, cfg.ServiceName)
	assert.Equal(t, testDefaultEndpoint, cfg.Endpoint)
	assert.Equal(t, ProtocolGRPC, cfg.Protocol)
	assert.True(t, cfg.Insecure)
	assert.Equal(t, DefaultSampleRate, cfg.SampleRate)
	assert.Nil(t, cfg.SkipPaths)
	assert.True(t, cfg.IsGRPC())
	assert.False(t, cfg.IsHTTP())
}

// TestDefaultConfig verifies DefaultConfig constructor outputs expected values.
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	verifyDefaultConfig(t, cfg)
	assert.NoError(t, cfg.Validate())
}

// TestSetDefaults verifies SetDefaults resets config fields to defaults.
func TestSetDefaults(t *testing.T) {
	cfg := &Config{
		Enabled:     false,
		ServiceName: "custom",
		Endpoint:    "custom:4317",
		Protocol:    ProtocolHTTP,
		Insecure:    false,
		SampleRate:  testSampleRateQuarter,
		SkipPaths:   []string{"/a"},
	}
	cfg.SetDefaults()
	verifyDefaultConfig(t, cfg)
}
