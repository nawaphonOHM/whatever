package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger_NewWithHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	h := slog.NewJSONHandler(buf, nil)
	l := NewWithHandler(h)
	require.NotNil(t, l)

	l.Info("message via custom handler")
	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "message via custom handler", logMap["msg"])
}

func TestLogger_Slog_And_NilSafety(t *testing.T) {
	var nilLogger *Logger
	assert.NotNil(t, nilLogger.Slog())
	assert.NotNil(t, nilLogger.ExitFunc())
	nilLogger.SetExitFunc(nil)
	assert.Nil(t, nilLogger.With("k", "v"))
	assert.Nil(t, nilLogger.WithGroup("grp"))
	nilLogger.TraceContext(context.Background(), "msg")
	nilLogger.FatalContext(context.Background(), "msg")

	l := New()
	assert.Equal(t, l.Logger, l.Slog())
}

func TestLogger_DisableTraceCorrelation(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(Config{
		Output:                  buf,
		Level:                   LevelInfo,
		DisableTraceCorrelation: true,
	})

	ctx := createTestSpanContext(t)
	l.InfoContext(ctx, "hello without trace correlation")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "hello without trace correlation", logMap["msg"])
	assert.Nil(t, logMap["trace_id"])
	assert.Nil(t, logMap["span_id"])
}

func emitGlobalContextLogs(ctx context.Context) {
	DebugContext(ctx, "debug ctx msg")
	InfoContext(ctx, "info ctx msg")
	WarnContext(ctx, "warn ctx msg")
	ErrorContext(ctx, "error ctx msg")
}

func verifyGlobalContextLogs(t *testing.T, out string) {
	assert.Contains(t, out, "debug ctx msg")
	assert.Contains(t, out, "info ctx msg")
	assert.Contains(t, out, "warn ctx msg")
	assert.Contains(t, out, "error ctx msg")
}

func TestGlobal_ContextFunctions(t *testing.T) {
	buf := &bytes.Buffer{}
	SetDefault(New(Config{Output: buf, Level: LevelDebug}))
	emitGlobalContextLogs(context.Background())
	verifyGlobalContextLogs(t, buf.String())
}

func TestLogger_InvalidLevelFallback(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(Config{Output: buf, Level: "INVALID_UNKNOWN_LEVEL"})
	l.Info("fallback info message")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "fallback info message", logMap["msg"])
}
