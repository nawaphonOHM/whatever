package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runIsGCPTableTests(t *testing.T) {
	t.Helper()
	tests := []struct {
		cfg      *Config
		name     string
		expected bool
	}{
		{name: "nil", cfg: nil, expected: false},
		{name: "empty", cfg: &Config{CloudProvider: ""}, expected: false},
		{name: "lowercase", cfg: &Config{CloudProvider: testValGCP}, expected: true},
		{name: "uppercase", cfg: &Config{CloudProvider: "GCP"}, expected: true},
		{name: "mixed", cfg: &Config{CloudProvider: "Gcp"}, expected: true},
		{name: "whitespace", cfg: &Config{CloudProvider: "  gcp \t\n"}, expected: true},
		{name: "aws", cfg: &Config{CloudProvider: testValAWS}, expected: false},
		{name: "azure", cfg: &Config{CloudProvider: "azure"}, expected: false},
		{name: "spaces", cfg: &Config{CloudProvider: "   "}, expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.cfg.IsGCP())
		})
	}
}

func TestConfig_IsGCP(t *testing.T) {
	runIsGCPTableTests(t)
}

func TestLoad_CloudProviderEnv(t *testing.T) {
	clearEnv(t)
	t.Setenv(testEnvCloudProvider, testValGCP)

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, testValGCP, cfg.CloudProvider)
	assert.True(t, cfg.IsGCP())
}

func TestLoad_DotEnvFile_CloudProvider(t *testing.T) {
	clearEnv(t)
	envPath := filepath.Join(t.TempDir(), ".env.gcp")
	content := []byte("OHM9996_CLOUD_PROVIDER=gcp\n" +
		"OHM9996_LOGGING_OUTPUT=stdout\n" +
		"OHM9996_LOGGING_LEVEL=info\n" +
		"OHM9996_LOGGING_FORMAT=json\n")
	require.NoError(t, os.WriteFile(envPath, content, testDotEnvPerm))

	cfg, err := Load(envPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, testValGCP, cfg.CloudProvider)
	assert.True(t, cfg.IsGCP())
}
