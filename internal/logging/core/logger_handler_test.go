package core

import (
	"bytes"
	"log/slog"
	"sync/atomic"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v2/internal/logging/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWithHandler(t *testing.T) {
	t.Run("valid handler", func(t *testing.T) {
		buf := &bytes.Buffer{}
		h := slog.NewJSONHandler(buf, nil)
		l := NewWithHandler(h)
		require.NotNil(t, l)

		l.Info("message via custom handler")
		logMap := parseJSONLog(t, buf.Bytes())
		assert.Equal(t, "message via custom handler", logMap["msg"])
	})

	t.Run("nil handler defaults safely", func(t *testing.T) {
		l := NewWithHandler(nil)
		require.NotNil(t, l)
		assert.NotNil(t, l.Slog())
	})
}

func TestLogger_WithAndWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	l := NewJSON(buf, config.LevelDebug)

	child := l.With("service", "billing").WithGroup("req")
	assert.Same(t, l.Worker(), child.Worker(), "derived logger must share parent's worker instance")

	child.Debug("payment processed", "amount", testPaymentAmount)

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "billing", logMap["service"])
	reqMap, ok := logMap["req"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(testPaymentAmount), reqMap["amount"])
}

func TestLogger_DerivedLogger_SharedWorker_AsyncFlush(t *testing.T) {
	buf, parent := newSharedWorkerLogger(t)
	child1 := parent.With("component", "auth")
	child2 := child1.WithGroup("session").With("user_id", "u-123")
	assert.Same(t, parent.Worker(), child1.Worker())
	assert.Same(t, parent.Worker(), child2.Worker())
	child2.Info("user logged in", "ip", "127.0.0.1")
	parent.Flush()
	assertSharedWorkerOutput(t, buf.String())
}

func TestLogger_DerivedLogger_Fatal(t *testing.T) {
	buf, parent, exitCode := newFatalTestLogger(t)
	child := parent.With("service", "orders")
	child.Fatal("database unreachable", "db", "primary")
	assert.Equal(t, int32(1), exitCode.Load())
	assertFatalOutput(t, buf.String())
}

const testCoreWorkerBuffer = 100

func newSharedWorkerLogger(t *testing.T) (*bytes.Buffer, *Logger) {
	t.Helper()
	buf := &bytes.Buffer{}
	h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	worker := central.New(slog.New(h), testCoreWorkerBuffer, nil)
	worker.Start()
	t.Cleanup(func() { require.NoError(t, worker.Close()) })
	return buf, &Logger{Logger: slog.New(h), handler: h, worker: worker}
}

func newFatalTestLogger(t *testing.T) (*bytes.Buffer, *Logger, *atomic.Int32) {
	t.Helper()
	buf := &bytes.Buffer{}
	h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	var exitCode atomic.Int32
	worker := central.New(slog.New(h), testCoreWorkerBuffer, func(code int) { exitCode.Store(int32(code)) })
	worker.Start()
	t.Cleanup(func() { require.NoError(t, worker.Close()) })
	logger := &Logger{Logger: slog.New(h), handler: h, worker: worker,
		exitFunc: func(code int) { exitCode.Store(int32(code)) }}
	return buf, logger, &exitCode
}

func assertSharedWorkerOutput(t *testing.T, out string) {
	t.Helper()
	assert.Contains(t, out, "user logged in")
	assert.Contains(t, out, `"component":"auth"`)
	assert.Contains(t, out, `"session":{"user_id":"u-123","ip":"127.0.0.1"}`)
}

func assertFatalOutput(t *testing.T, out string) {
	t.Helper()
	assert.Contains(t, out, "database unreachable")
	assert.Contains(t, out, `"service":"orders"`)
	assert.Contains(t, out, `"db":"primary"`)
	assert.Contains(t, out, `"exit_flag":"ExitAbnormal"`)
}
