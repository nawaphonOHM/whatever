// Package middleware provides unit tests for panic recovery middleware.
package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func verifyStandardJSONPanicLog(t *testing.T, entry map[string]any) {
	t.Helper()
	assert.Equal(t, "ERROR", entry["level"])
	assert.Equal(t, expectedLogMsg, entry["msg"])
	assert.Equal(t, panicReqID, entry["request_id"])
	assert.Equal(t, panicPath, entry["path"])
	assert.Equal(t, http.MethodGet, entry["method"])
	assert.Equal(t, panicMsg, entry["error"])
	assert.NotEmpty(t, entry["stack_trace"])
}

// TestRecovery_HandlesPanic verifies that panics are caught properly.
func TestRecovery_HandlesPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	r := setupPanicRouter(logger)

	w := executePanicRequest(t, r)
	verifyPanicResponse(t, w)
	entry := parseJSONLogEntry(t, buf.Bytes())
	verifyStandardJSONPanicLog(t, entry)
}
