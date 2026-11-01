package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

const (
	logFieldSeverity   = "severity"
	msgStorageFail     = "failed to connect to storage; exiting abnormally"
	msgStorageInit     = "initializing storage client"
	msgStorageComplete = "storage client initialization complete"
)

func assertStorageInitEntryGCP(t *testing.T, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected initializing storage client entry")
	assertStorageGCPTimestamp(t, entry)
	assertStorageGCPSeverity(t, entry, "INFO")
	assert.Equal(t, "s3-compatible", entry["provider"])
	assert.Equal(t, "test-bucket", entry["bucket"])
	assert.Equal(t, "ap-southeast-2", entry["region"])
	assert.Equal(t, "https://storage.example.test", entry["endpoint"])
	assert.Equal(t, true, entry["use_path_style"])
	assert.Equal(t, "when_required", entry["checksum_mode"])
}

func assertStorageCompleteEntryGCP(t *testing.T, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected storage client initialization complete entry")
	assertStorageGCPTimestamp(t, entry)
	assertStorageGCPSeverity(t, entry, "INFO")
	assert.Equal(t, "s3-compatible", entry["provider"])
	assert.Equal(t, "test-bucket", entry["bucket"])
	assert.Equal(t, "ap-southeast-2", entry["region"])
	assert.Equal(t, "https://storage.example.test", entry["endpoint"])
	assert.Equal(t, true, entry["use_path_style"])
	assert.Equal(t, "when_required", entry["checksum_mode"])
}

func gcpStorageTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	cfg.Region = "ap-southeast-2"
	cfg.Endpoint = "https://storage.example.test"
	cfg.UsePathStyle = true
	cfg.AccessKeyID = "test-key"
	cfg.SecretAccessKey = "test-secret"
	return cfg
}
