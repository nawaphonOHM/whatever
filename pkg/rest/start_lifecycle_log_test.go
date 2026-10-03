package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/logging/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseJSONLogEntries(data []byte) []map[string]any {
	var entries []map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	for decoder.More() {
		var entry map[string]any
		if err := decoder.Decode(&entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	return entries
}

func setupTestLogger(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()
	buf := new(bytes.Buffer)
	reset := core.CaptureLogs(buf)
	return buf, reset
}

func findLogByMsg(entries []map[string]any, msg string) map[string]any {
	for _, entry := range entries {
		if entry["msg"] == msg {
			return entry
		}
	}
	return nil
}

func assertPublicStartLog(t *testing.T, port int, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected 'starting REST server' log entry")
	assert.Equal(t, "INFO", entry["level"])
	assert.Equal(t, "application", entry["service"])
	assert.Equal(t, "1.0.0", entry["version"])
	assert.Equal(t, testHost, entry["host"])
	assert.Equal(t, float64(port), entry["port"])
	assert.Equal(t, gin.TestMode, entry["mode"])
	assert.Equal(t, "10s", entry["shutdown_timeout"])
}

func assertPublicShutdownLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	sigEntry := findLogByMsg(entries, "shutting down REST server gracefully")
	require.NotNil(t, sigEntry, "expected 'shutting down REST server gracefully' log entry")
	assert.Equal(t, "INFO", sigEntry["level"])
	assert.Equal(t, "10s", sigEntry["shutdown_timeout"])

	doneEntry := findLogByMsg(entries, "REST server shut down successfully")
	require.NotNil(t, doneEntry, "expected 'REST server shut down successfully' log entry")
	assert.Equal(t, "INFO", doneEntry["level"])
}

func TestStartREST_LifecycleLogging_Success(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(newTestE2EBlueprint())

	assertStatusOK(t, fmt.Sprintf("http://%s:%d/api/v1/ping", testHost, port))
	signalAndWait(t, errCh)

	entries := parseJSONLogEntries(buf.Bytes())
	assertPublicStartLog(t, port, findLogByMsg(entries, "starting REST server"))
	assertPublicShutdownLogs(t, entries)
}
