package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func TestConnect_InfoLevel_SuppressesDebugLogs(t *testing.T) {
	for _, mode := range storageLogModes() {
		t.Run(mode.name, func(t *testing.T) {
			assertInfoLevelSuppressesDebug(t, mode)
		})
	}
}

func assertLifecycleInfoLogsPresent(t *testing.T, entries []map[string]any) {
	t.Helper()
	assert.NotNil(t, findStorageLog(entries, "initializing storage client"))
	assert.NotNil(t, findStorageLog(entries, "storage client initialization complete"))
}

func assertDebugLogsSuppressed(t *testing.T, entries []map[string]any, levelKey string) {
	t.Helper()
	assert.Nil(t, findStorageLog(entries, "evaluating storage provider detection"))
	assert.Nil(t, findStorageLog(entries, "storage client compatibility configuration"))
	for _, entry := range entries {
		assert.NotEqual(t, "DEBUG", entry[levelKey])
	}
}

func captureInfoFilterLogs(t *testing.T, mode storageLogMode) []map[string]any {
	t.Helper()
	isolateAWSEnv(t)
	var buf bytes.Buffer
	setStorageTestLogger(t, mode.handler(&buf))
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	cfg.Endpoint = "https://storage.googleapis.com"
	_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
	require.NoError(t, err)
	return parseStorageLogs(t, buf.Bytes())
}

func assertInfoLevelSuppressesDebug(t *testing.T, mode storageLogMode) {
	t.Helper()
	entries := captureInfoFilterLogs(t, mode)
	assertLifecycleInfoLogsPresent(t, entries)
	assertDebugLogsSuppressed(t, entries, mode.levelKey)
}
