package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDotEnvPerm        = 0o600
	testDotEnvServiceName = "file-loaded-service"
	testDotEnvEndpoint    = "http://otel-collector:4318"
)

// createDotEnvFile writes a temporary dotenv configuration file.
func createDotEnvFile(t *testing.T) string {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env.otel")
	content := []byte(`
OHM9996_OTEL_ENABLED=true
OHM9996_OTEL_SERVICE_NAME=file-loaded-service
OHM9996_OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4318
OHM9996_OTEL_EXPORTER_OTLP_PROTOCOL=http
OHM9996_OTEL_INSECURE=true
OHM9996_OTEL_SAMPLE_RATE=0.25
OHM9996_OTEL_SKIP_PATHS=/health,/ready
`)
	require.NoError(t, os.WriteFile(envPath, content, testDotEnvPerm))
	return envPath
}

// verifyDotEnvLoadedConfig asserts loaded configuration fields match dotenv contents.
func verifyDotEnvLoadedConfig(t *testing.T, cfg *Config) {
	assert.True(t, cfg.Enabled)
	assert.Equal(t, testDotEnvServiceName, cfg.ServiceName)
	assert.Equal(t, testDotEnvEndpoint, cfg.Endpoint)
	assert.Equal(t, ProtocolHTTP, cfg.Protocol)
	assert.True(t, cfg.Insecure)
	assert.Equal(t, testSampleRateQuarter, cfg.SampleRate)
	assert.Equal(t, []string{testHealthPath, testReadyPath}, cfg.SkipPaths)
	assert.True(t, cfg.IsHTTP())
}

// TestLoad_DotEnvFile tests loading configuration from a dotenv file.
func TestLoad_DotEnvFile(t *testing.T) {
	clearEnv(t)
	envPath := createDotEnvFile(t)
	cfg, err := Load(envPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	verifyDotEnvLoadedConfig(t, cfg)
	assert.NoError(t, cfg.Validate())
}
