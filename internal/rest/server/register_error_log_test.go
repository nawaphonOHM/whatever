package server

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testExecutionErrText = "downstream failure"

type executionFailureCase struct {
	name string
	msg  string
}

func executionFailureCases() []executionFailureCase {
	return []executionFailureCase{
		{name: "middleware", msg: msgMiddlewareFailed},
		{name: "handler", msg: msgHandlerFailed},
	}
}

func captureJSONExecutionFailure(t *testing.T, msg string) map[string]any {
	t.Helper()
	var buf bytes.Buffer
	reset := core.CaptureLogs(&buf)
	logExecutionFailure(context.Background(), msg, errors.New(testExecutionErrText))
	reset()
	entry := findLogEntry(parseJSONLogs(t, buf.Bytes()), msg)
	require.NotNil(t, entry, "expected %q log entry", msg)
	return entry
}

func captureGCPExecutionFailure(t *testing.T, msg string) map[string]any {
	t.Helper()
	var buf bytes.Buffer
	reset := core.SetTestLogger(core.NewWithHandler(core.NewGCPHandler(&buf, nil)))
	logExecutionFailure(context.Background(), msg, errors.New(testExecutionErrText))
	reset()
	entry := findGCPLogEntry(parseJSONLogs(t, buf.Bytes()), msg)
	require.NotNil(t, entry, "expected %q GCP log entry", msg)
	return entry
}

func findGCPLogEntry(entries []map[string]any, msg string) map[string]any {
	for _, entry := range entries {
		if entry["message"] == msg {
			return entry
		}
	}
	return nil
}

func TestExecutionFailureLog_JSONHandler(t *testing.T) {
	for _, tc := range executionFailureCases() {
		t.Run(tc.name, func(t *testing.T) {
			entry := captureJSONExecutionFailure(t, tc.msg)
			assert.Equal(t, "ERROR", entry["level"])
			assert.Equal(t, testExecutionErrText, entry["error"])
			assert.NotContains(t, entry, "err")
		})
	}
}

func TestExecutionFailureLog_GCPHandler(t *testing.T) {
	for _, tc := range executionFailureCases() {
		t.Run(tc.name, func(t *testing.T) {
			entry := captureGCPExecutionFailure(t, tc.msg)
			assert.Equal(t, "ERROR", entry["severity"])
			assert.Equal(t, tc.msg, entry["message"])
			assert.Equal(t, testExecutionErrText, entry["error"])
			assert.NotContains(t, entry, "err")
		})
	}
}
