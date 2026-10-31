package server

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const gcpInfoLevel = "INFO"

func setupGCPLifecycleLogger(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()
	buf := new(bytes.Buffer)
	handler := core.NewGCPHandler(buf, nil)
	reset := core.SetTestLogger(core.NewWithHandler(handler))
	return buf, reset
}

func assertGCPStartBase(t *testing.T, srv *Server, port int, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected GCP 'starting REST server' log entry")
	assert.Equal(t, gcpInfoLevel, entry["severity"])
	assert.Equal(t, srv.Config.DisplayName, entry["service"])
	assert.Equal(t, srv.Config.AppVersion, entry["version"])
	assert.Equal(t, fmt.Sprintf("%s:%d", testHost, port), entry["addr"])
	assert.Equal(t, testHost, entry["host"])
	assert.Equal(t, float64(port), entry["port"])
}

func assertGCPStartFlags(t *testing.T, srv *Server, entry map[string]any) {
	t.Helper()
	assert.Equal(t, gin.TestMode, entry["mode"])
	assert.Equal(t, srv.Config.ShutdownTimeout.String(), entry["shutdown_timeout"])
	assert.Equal(t, srv.Config.EnableAccessLog, entry["access_log"])
	assert.Equal(t, srv.Config.EnableMetrics, entry["metrics"])
}

func assertGCPStartLog(t *testing.T, srv *Server, port int, entry map[string]any) {
	t.Helper()
	assertGCPStartBase(t, srv, port, entry)
	assertGCPStartFlags(t, srv, entry)
}

func assertGCPShutdownLogs(t *testing.T, srv *Server, entries []map[string]any) {
	t.Helper()
	sigEntry := findGCPLogEntry(entries, "shutting down REST server gracefully")
	require.NotNil(t, sigEntry, "expected GCP 'shutting down REST server gracefully' log entry")
	assert.Equal(t, gcpInfoLevel, sigEntry["severity"])
	assert.Equal(t, srv.Config.ShutdownTimeout.String(), sigEntry["shutdown_timeout"])

	doneEntry := findGCPLogEntry(entries, "REST server shut down successfully")
	require.NotNil(t, doneEntry, "expected GCP 'REST server shut down successfully' log entry")
	assert.Equal(t, gcpInfoLevel, doneEntry["severity"])
}

func TestServer_LifecycleLogging_GCPFormat(t *testing.T) {
	buf, reset := setupGCPLifecycleLogger(t)
	defer reset()

	srv, port := newBoundServer(t)
	cancel, errCh := startWithCancel(srv)

	assertPingOK(t, port)
	cancel()
	require.NoError(t, waitErr(t, errCh))

	entries := parseJSONLogs(t, buf.Bytes())
	assertGCPStartLog(t, srv, port, findGCPLogEntry(entries, "starting REST server"))
	assertGCPShutdownLogs(t, srv, entries)
}
