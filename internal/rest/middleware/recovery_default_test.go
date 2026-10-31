package middleware

import (
	"bytes"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/stretchr/testify/assert"
)

// TestRecovery_NilLogger_UsesActiveJSONLogger verifies nil logger resolves active JSON default.
func TestRecovery_NilLogger_UsesActiveJSONLogger(t *testing.T) {
	var buf bytes.Buffer
	cleanup := core.CaptureLogs(&buf)
	defer cleanup()

	r := setupDefaultRecoveryRouter()
	w := executePanicRequest(t, r)

	verifyPanicResponse(t, w)
	entry := parseJSONLogEntry(t, buf.Bytes())
	verifyStandardJSONPanicLog(t, entry)
}

func assertTextPanicLogs(t *testing.T, logText string) {
	t.Helper()
	assert.Contains(t, logText, "ERROR")
	assert.Contains(t, logText, expectedLogMsg)
	assert.Contains(t, logText, panicReqID)
	assert.Contains(t, logText, panicMsg)
}

// TestRecovery_NilLogger_UsesActiveTextLogger verifies nil logger resolves active Text default.
func TestRecovery_NilLogger_UsesActiveTextLogger(t *testing.T) {
	var buf bytes.Buffer
	cleanup := core.CaptureTextLogs(&buf)
	defer cleanup()

	r := setupDefaultRecoveryRouter()
	w := executePanicRequest(t, r)

	verifyPanicResponse(t, w)
	assertTextPanicLogs(t, buf.String())
}

// TestRecoveryWithLogger_Nil_Fallback verifies RecoveryWithLogger(nil) uses active default logger.
func TestRecoveryWithLogger_Nil_Fallback(t *testing.T) {
	var buf bytes.Buffer
	cleanup := core.CaptureLogs(&buf)
	defer cleanup()

	r := setupPanicRouter(nil)
	w := executePanicRequest(t, r)

	verifyPanicResponse(t, w)
	entry := parseJSONLogEntry(t, buf.Bytes())
	verifyStandardJSONPanicLog(t, entry)
}
