package middleware

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/stretchr/testify/assert"
)

func verifyGCPPanicLogEntry(t *testing.T, entry map[string]any) {
	t.Helper()
	assert.Equal(t, "ERROR", entry["severity"])
	assert.Equal(t, expectedLogMsg, entry["message"])
	assert.Equal(t, panicReqID, entry["request_id"])
	assert.Equal(t, panicPath, entry["path"])
	assert.Equal(t, http.MethodGet, entry["method"])
	assert.Equal(t, panicMsg, entry["error"])
	assert.NotEmpty(t, entry["stack_trace"])
}

func setupGCPTestLogger(buf *bytes.Buffer) (*core.Logger, func()) {
	h := core.NewGCPHandler(buf, nil)
	l := core.NewWithHandler(h)
	cleanup := core.SetTestLogger(l)
	return l, cleanup
}

// TestRecovery_GCPCloudLogging verifies panic logging in GCP format with default logger.
func TestRecovery_GCPCloudLogging(t *testing.T) {
	var buf bytes.Buffer
	logger, cleanup := setupGCPTestLogger(&buf)
	defer cleanup()

	r := setupDefaultRecoveryRouter()
	w := executePanicRequest(t, r)

	logger.Flush()
	verifyPanicResponse(t, w)
	entry := parseJSONLogEntry(t, buf.Bytes())
	verifyGCPPanicLogEntry(t, entry)
}

// TestRecoveryWithLogger_GCPCloudLogging verifies panic logging in GCP format with explicit logger.
func TestRecoveryWithLogger_GCPCloudLogging(t *testing.T) {
	var buf bytes.Buffer
	h := core.NewGCPHandler(&buf, nil)
	logger := core.NewWithHandler(h)
	defer logger.Flush()

	r := setupPanicRouter(logger.Slog())
	w := executePanicRequest(t, r)

	logger.Flush()
	verifyPanicResponse(t, w)
	entry := parseJSONLogEntry(t, buf.Bytes())
	verifyGCPPanicLogEntry(t, entry)
}
