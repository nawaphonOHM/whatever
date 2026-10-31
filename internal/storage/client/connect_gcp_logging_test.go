package client

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func TestHandleConnectionErrorLogsGCP(t *testing.T) {
	assertStorageConnectionFailureGCP(t)
}

func assertStorageConnectionFailureGCP(t *testing.T) {
	t.Helper()
	entry := captureStorageConnectionFailureGCP(t)
	assertStorageConnectionFailureEntry(t, entry)
}

func captureStorageConnectionFailureGCP(t *testing.T) map[string]any {
	t.Helper()
	var buf bytes.Buffer
	setupStorageGCPTestLogger(t, &buf)
	previousExit := SetExitFunc(func(int) {})
	t.Cleanup(func() { SetExitFunc(previousExit) })
	testErr := errors.New("storage network timeout")
	require.ErrorIs(t, handleConnectionError(context.Background(), testErr), testErr)
	return findStorageLog(parseStorageLogs(t, buf.Bytes()), msgStorageFail)
}

func assertStorageConnectionFailureEntry(t *testing.T, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry)
	assertStorageGCPTimestamp(t, entry)
	assertStorageGCPSeverity(t, entry, "ERROR")
	assert.Equal(t, msgStorageFail, entry["message"])
	assert.Equal(t, "storage network timeout", entry["error"])
	assert.Equal(t, "ExitAbnormal", entry["exit_flag"])
}

func captureStorageLifecycleLogsGCP(t *testing.T, cfg *config.Config) []map[string]any {
	t.Helper()
	var buf bytes.Buffer
	setupStorageGCPTestLogger(t, &buf)
	_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
	require.NoError(t, err)
	return parseStorageLogs(t, buf.Bytes())
}

func TestStorageLifecycleLogsGCP(t *testing.T) {
	entries := captureStorageLifecycleLogsGCP(t, gcpStorageTestConfig())
	assertStorageInitEntryGCP(t, findStorageLog(entries, msgStorageInit))
	assertStorageCompleteEntryGCP(t, findStorageLog(entries, msgStorageComplete))
}
