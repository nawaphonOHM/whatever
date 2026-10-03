package core

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
	"github.com/stretchr/testify/assert"
)

func exerciseNilContextLogs(ctx context.Context, nilLogger *Logger) {
	nilLogger.TraceContext(ctx, testMsg)
	nilLogger.DebugContext(ctx, testMsg)
	nilLogger.InfoContext(ctx, testMsg)
	nilLogger.WarnContext(ctx, testMsg)
	nilLogger.ErrorContext(ctx, testMsg)
	nilLogger.FatalContext(ctx, testMsg)
	nilLogger.Log(ctx, config.SlogLevelInfo, testMsg)
	nilLogger.LogAttrs(ctx, config.SlogLevelInfo, testMsg)
}

func exerciseNilSimpleLogs(nilLogger *Logger) {
	nilLogger.Trace(testMsg)
	nilLogger.Debug(testMsg)
	nilLogger.Info(testMsg)
	nilLogger.Warn(testMsg)
	nilLogger.Error(testMsg)
	nilLogger.Fatal(testMsg)
}

func verifyNilLoggerProperties(t *testing.T, nilLogger *Logger) {
	t.Helper()
	assert.NotNil(t, nilLogger.Slog())
	assert.Nil(t, nilLogger.Handler())
	assert.NotNil(t, nilLogger.ExitFunc())
	nilLogger.SetExitFunc(nil)
	assert.NoError(t, nilLogger.Close())
	assert.Nil(t, nilLogger.With("k", "v"))
	assert.Nil(t, nilLogger.WithGroup("grp"))
}

func TestLogger_Slog_And_NilSafety(t *testing.T) {
	var nilLogger *Logger
	verifyNilLoggerProperties(t, nilLogger)
	exerciseNilSimpleLogs(nilLogger)
	exerciseNilContextLogs(context.Background(), nilLogger)

	l := New()
	assert.Equal(t, l.Logger, l.Slog())
	assert.NotNil(t, l.Handler())
	assert.NoError(t, l.Close())
}

func TestLogger_NilContext_Safe(t *testing.T) {
	buf := &bytes.Buffer{}
	l := NewJSON(buf, config.LevelTrace)

	var nilCtx context.Context
	l.TraceContext(nilCtx, "trace msg")
	l.DebugContext(nilCtx, "debug msg")
	l.InfoContext(nilCtx, "info msg")
	l.WarnContext(nilCtx, "warn msg")
	l.ErrorContext(nilCtx, "error msg")

	lines := strings.Split(strings.TrimSpace(buf.String()), testNewline)
	assert.Len(t, lines, testFiveCount)
}
