package client

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
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
	setStorageTestLogger(t, core.NewGCPHandler(&buf, nil))
	previousExit := SetExitFunc(func(int) {})
	t.Cleanup(func() { SetExitFunc(previousExit) })
	testErr := errors.New("storage network timeout")
	require.ErrorIs(t, handleConnectionError(context.Background(), testErr), testErr)
	entries := parseStorageLogs(t, buf.Bytes())
	require.NotEmpty(t, entries)
	return entries[0]
}

func assertStorageConnectionFailureEntry(t *testing.T, entry map[string]any) {
	t.Helper()
	assert.Equal(t, "ERROR", entry["severity"])
	assert.Equal(t, "failed to connect to storage; exiting abnormally", entry["message"])
	assert.Equal(t, "storage network timeout", entry["error"])
	assert.Equal(t, "ExitAbnormal", entry["exit_flag"])
}
