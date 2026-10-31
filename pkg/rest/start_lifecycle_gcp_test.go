package rest

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPublicGCPLogger(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()
	buf := new(bytes.Buffer)
	handler := core.NewGCPHandler(buf, nil)
	reset := core.SetTestLogger(core.NewWithHandler(handler))
	return buf, reset
}

func findPublicGCPLog(entries []map[string]any, msg string) map[string]any {
	for _, entry := range entries {
		if entry["message"] == msg {
			return entry
		}
	}
	return nil
}

func assertPublicGCPStartLog(t *testing.T, port int, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected GCP 'starting REST server' log entry")
	assert.Equal(t, "INFO", entry["severity"])
	assert.Equal(t, "application", entry["service"])
	assert.Equal(t, "1.0.0", entry["version"])
	assert.Equal(t, testHost, entry["host"])
	assert.Equal(t, float64(port), entry["port"])
	assert.Equal(t, gin.TestMode, entry["mode"])
	assert.Equal(t, "10s", entry["shutdown_timeout"])
}

func assertPublicGCPShutdownLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	sigEntry := findPublicGCPLog(entries, "shutting down REST server gracefully")
	require.NotNil(t, sigEntry, "expected GCP 'shutting down REST server gracefully' log entry")
	assert.Equal(t, "INFO", sigEntry["severity"])
	assert.Equal(t, "10s", sigEntry["shutdown_timeout"])

	doneEntry := findPublicGCPLog(entries, "REST server shut down successfully")
	require.NotNil(t, doneEntry, "expected GCP 'REST server shut down successfully' log entry")
	assert.Equal(t, "INFO", doneEntry["severity"])
}

func TestStartREST_LifecycleLogging_GCPFormat(t *testing.T) {
	buf, reset := setupPublicGCPLogger(t)
	defer reset()

	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(newTestE2EBlueprint())

	assertStatusOK(t, fmt.Sprintf("http://%s:%d/api/v1/ping", testHost, port))
	signalAndWait(t, errCh)

	entries := parseJSONLogEntries(t, buf.Bytes())
	assertPublicGCPStartLog(t, port, findPublicGCPLog(entries, "starting REST server"))
	assertPublicGCPShutdownLogs(t, entries)
}
