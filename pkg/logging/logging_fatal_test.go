package logging_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/logging/core"
	"github.com/nawaphonOHM/whatever/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func createFatalTestLogger(buf *bytes.Buffer, exitCode *int, called *bool) func() {
	testLogger := core.NewJSON(buf, logging.LevelDebug)
	testLogger.SetExitFunc(func(code int) {
		*exitCode = code
		*called = true
	})
	return core.SetTestLogger(testLogger)
}

func verifyFatalEntry(t *testing.T, data []byte) {
	t.Helper()
	entry := parseJSONLog(t, data)
	assert.Equal(t, testLevelFatal, entry["level"])
	assert.Equal(t, "fatal error occurred", entry["msg"])
	assert.Equal(t, "disk full", entry["reason"])
}

func TestLogging_Fatal(t *testing.T) {
	buf := new(bytes.Buffer)
	var exitCode int
	var called bool

	cleanup := createFatalTestLogger(buf, &exitCode, &called)
	defer cleanup()

	logging.Fatal("fatal error occurred", "reason", "disk full")
	assert.True(t, called)
	assert.Equal(t, 1, exitCode)
	verifyFatalEntry(t, buf.Bytes())
}

func createFatalSpanContext(t *testing.T) context.Context {
	t.Helper()
	traceID, err := trace.TraceIDFromHex(testFatalTraceIDHex)
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex(testFatalSpanIDHex)
	require.NoError(t, err)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(context.Background(), spanCtx)
}

func verifyFatalContextEntry(t *testing.T, data []byte) {
	t.Helper()
	entry := parseJSONLog(t, data)
	assert.Equal(t, testLevelFatal, entry["level"])
	assert.Equal(t, "fatal with context", entry["msg"])
	assert.Equal(t, "fatal err", entry["err"])
	assert.Equal(t, testFatalTraceIDHex, entry["trace_id"])
	assert.Equal(t, testFatalSpanIDHex, entry["span_id"])
}

func TestLogging_FatalContext(t *testing.T) {
	buf := new(bytes.Buffer)
	var exitCode int
	var called bool

	cleanup := createFatalTestLogger(buf, &exitCode, &called)
	defer cleanup()

	ctx := createFatalSpanContext(t)
	logging.FatalContext(ctx, "fatal with context", "err", errors.New("fatal err"))
	assert.True(t, called)
	assert.Equal(t, 1, exitCode)
	verifyFatalContextEntry(t, buf.Bytes())
}
