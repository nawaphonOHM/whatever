package client

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

type storageLogMode struct {
	name     string
	levelKey string
	format   string
}

func storageLogModes() []storageLogMode {
	return []storageLogMode{
		{name: "json", levelKey: "level", format: "json"},
		{name: "gcp", levelKey: "severity", format: "gcp"},
	}
}

func (m storageLogMode) handler(output io.Writer) slog.Handler {
	if m.format == "gcp" {
		return core.NewGCPHandler(output, nil)
	}
	return slog.NewJSONHandler(output, nil)
}

func storageLoggingConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	cfg.Region = "ap-southeast-2"
	cfg.Endpoint = "https://storage.example.test"
	cfg.UsePathStyle = true
	return cfg
}

func TestStorageInitializationAndCompletionLogs_JSONAndGCP(t *testing.T) {
	for _, mode := range storageLogModes() {
		t.Run(mode.name, func(t *testing.T) {
			assertStorageLifecycleLogs(t, mode)
		})
	}
}

func assertStorageLifecycleLogs(t *testing.T, mode storageLogMode) {
	t.Helper()
	var output bytes.Buffer
	setStorageTestLogger(t, mode.handler(&output))
	cfg := storageLoggingConfig()
	require.NoError(t, validateAndLogConfig(context.Background(), cfg, contracts.ProviderAWS))
	logStartupComplete(context.Background(), contracts.ProviderAWS, cfg.Bucket)
	entries := parseStorageLogs(t, output.Bytes())
	assertStorageInitializationEntry(t, findStorageLog(entries, "initializing storage client"), mode)
	assertStorageCompletionEntry(t, findStorageLog(entries, "storage client initialization complete"), mode)
}

func assertStorageInitializationEntry(t *testing.T, entry map[string]any, mode storageLogMode) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, "INFO", entry[mode.levelKey])
	assert.Equal(t, "aws", entry["provider"])
	assert.Equal(t, "test-bucket", entry["bucket"])
	assert.Equal(t, "ap-southeast-2", entry["region"])
	assert.Equal(t, "https://storage.example.test", entry["endpoint"])
	assert.Equal(t, true, entry["use_path_style"])
}

func assertStorageCompletionEntry(t *testing.T, entry map[string]any, mode storageLogMode) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, "INFO", entry[mode.levelKey])
	assert.Equal(t, "aws", entry["provider"])
	assert.Equal(t, "test-bucket", entry["bucket"])
}
