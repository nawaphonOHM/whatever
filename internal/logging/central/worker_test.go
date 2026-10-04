package central

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/callstack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testBufferSize = 100
	testFrameLine  = 88
)

func newTestWorker(t *testing.T, logger *slog.Logger, exitFunc func(int)) *Worker {
	t.Helper()
	w := New(logger, testBufferSize, exitFunc)
	w.Start()
	t.Cleanup(func() { require.NoError(t, w.Close()) })
	return w
}

func newJSONLogger(buf *bytes.Buffer) *slog.Logger {
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(handler)
}

func TestWorker_AsyncDelivery(t *testing.T) {
	buf := &bytes.Buffer{}
	w := newTestWorker(t, newJSONLogger(buf), nil)
	ok := w.Enqueue(&LogEvent{Message: "async event test", Level: slog.LevelInfo,
		Attrs: []slog.Attr{slog.String("component", "test")}, Timestamp: time.Now()})
	assert.True(t, ok)
	w.Flush()
	assertAsyncOutput(t, buf.String())
}

func assertAsyncOutput(t *testing.T, out string) {
	t.Helper()
	assert.Contains(t, out, "async event test")
	assert.Contains(t, out, `"component":"test"`)
	assert.Contains(t, out, `"level":"INFO"`)
}

func TestWorker_ExitGraceful(t *testing.T) {
	buf := &bytes.Buffer{}
	var exitCode atomic.Int32
	w := newTestWorker(t, newJSONLogger(buf), func(code int) { exitCode.Store(int32(code)) })
	w.Exit(context.Background(), "graceful shutdown", ExitGraceful, nil, nil)
	assert.Equal(t, int32(0), exitCode.Load())
	assert.Contains(t, buf.String(), "graceful shutdown")
	assert.Contains(t, buf.String(), `"exit_flag":"ExitGraceful"`)
}

func TestWorker_ExitAbnormal_WithCallStack(t *testing.T) {
	buf := &bytes.Buffer{}
	var exitCode atomic.Int32
	w := newTestWorker(t, newJSONLogger(buf), func(code int) { exitCode.Store(int32(code)) })
	stack := callstack.New()
	stack.Push(&callstack.Frame{
		Function: "failingOperation",
		Package:  "testpkg",
		File:     "test.go",
		Line:     testFrameLine,
	})

	testErr := errors.New("connection failed")
	w.Exit(context.Background(), "fatal connection error", ExitAbnormal, testErr, stack)

	assert.Equal(t, int32(1), exitCode.Load())
	assertAbnormalOutput(t, buf.String())
}

func assertAbnormalOutput(t *testing.T, out string) {
	t.Helper()
	assert.Contains(t, out, "fatal connection error")
	assert.Contains(t, out, `"exit_flag":"ExitAbnormal"`)
	assert.Contains(t, out, "connection failed")
	assert.Contains(t, out, "Call stack snapshot")
	assert.Contains(t, out, "testpkg.failingOperation")
}
