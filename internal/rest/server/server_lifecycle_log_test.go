package server

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/logging/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseJSONLogs(data []byte) []map[string]any {
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

func setupLifecycleLogger(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()
	buf := new(bytes.Buffer)
	reset := core.CaptureLogs(buf)
	return buf, reset
}

func findLogEntry(entries []map[string]any, msg string) map[string]any {
	for _, entry := range entries {
		if entry["msg"] == msg {
			return entry
		}
	}
	return nil
}

func assertStartLog(t *testing.T, srv *Server, port int, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected 'starting REST server' log entry")
	assert.Equal(t, "INFO", entry["level"])
	assert.Equal(t, srv.Config.DisplayName, entry["service"])
	assert.Equal(t, testHost, entry["host"])
	assert.Equal(t, float64(port), entry["port"])
	assert.Equal(t, gin.TestMode, entry["mode"])
	assert.Equal(t, srv.Config.ShutdownTimeout.String(), entry["shutdown_timeout"])
	assert.Equal(t, srv.Config.EnableAccessLog, entry["access_log"])
	assert.Equal(t, srv.Config.EnableMetrics, entry["metrics"])
}

func assertShutdownLogs(t *testing.T, srv *Server, entries []map[string]any) {
	t.Helper()
	sigEntry := findLogEntry(entries, "shutting down REST server gracefully")
	require.NotNil(t, sigEntry, "expected 'shutting down REST server gracefully' log entry")
	assert.Equal(t, "INFO", sigEntry["level"])
	assert.Equal(t, srv.Config.ShutdownTimeout.String(), sigEntry["shutdown_timeout"])

	doneEntry := findLogEntry(entries, "REST server shut down successfully")
	require.NotNil(t, doneEntry, "expected 'REST server shut down successfully' log entry")
	assert.Equal(t, "INFO", doneEntry["level"])
}

func TestServer_LifecycleLogging_StartupAndShutdown(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	srv, port := newBoundServer(t)
	cancel, errCh := startWithCancel(srv)

	assertPingOK(t, port)
	cancel()
	require.NoError(t, waitErr(t, errCh))

	entries := parseJSONLogs(buf.Bytes())
	assertStartLog(t, srv, port, findLogEntry(entries, "starting REST server"))
	assertShutdownLogs(t, srv, entries)
}
