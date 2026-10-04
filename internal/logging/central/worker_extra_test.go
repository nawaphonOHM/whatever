package central

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorker_FullBufferFallback(t *testing.T) {
	buf := &bytes.Buffer{}
	w := New(newJSONLogger(buf), 1, nil)
	ok := w.Enqueue(&LogEvent{Message: "fallback event", Level: slog.LevelWarn})
	assert.False(t, ok)
	assert.Contains(t, buf.String(), "fallback event")
}

func TestWorker_SettersAndGetters(t *testing.T) {
	w := New(nil, 0, nil)
	require.NotNil(t, w.Logger())
	require.NotNil(t, w.ExitFunc())
	called := false
	w.SetExitFunc(func(code int) { called = code == 0 })
	w.ExitFunc()(0)
	assert.True(t, called)
	newLogger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	w.SetLogger(newLogger)
	assert.Equal(t, newLogger, w.Logger())
}

func TestPackageLevel_CentralWorker(t *testing.T) {
	buf := &bytes.Buffer{}
	exitCode := setupPackageWorker(t, newJSONLogger(buf))
	Dispatch(NewLogEvent(context.Background(), slog.LevelInfo, "pkg level central log"))
	FlushDefault()
	assert.Contains(t, buf.String(), "pkg level central log")
	assertPackageExits(t, exitCode)
}

func setupPackageWorker(t *testing.T, logger *slog.Logger) *atomic.Int32 {
	t.Helper()
	var exitCode atomic.Int32
	exitCode.Store(-1)
	w := newTestWorker(t, logger, func(code int) { exitCode.Store(int32(code)) })
	previous := DefaultWorker()
	SetDefaultWorker(w)
	t.Cleanup(func() { SetDefaultWorker(previous) })
	return &exitCode
}

func assertPackageExits(t *testing.T, exitCode *atomic.Int32) {
	t.Helper()
	ExitWithGraceful(context.Background(), "pkg graceful")
	assert.Equal(t, int32(0), exitCode.Load())
	ExitWithAbnormal(context.Background(), "pkg abnormal", errors.New("err"), nil)
	assert.Equal(t, int32(1), exitCode.Load())
}

func assertBufferContainsAll(t *testing.T, buf *bytes.Buffer, messages ...string) {
	t.Helper()
	out := buf.String()
	for _, msg := range messages {
		assert.Contains(t, out, msg)
	}
}

func TestPackageLevel_Exit(t *testing.T) {
	buf := &bytes.Buffer{}
	exitCode := setupPackageWorker(t, newJSONLogger(buf))

	Exit(context.Background(), "pkg exit none", ExitNone)
	assert.Equal(t, int32(-1), exitCode.Load())

	Exit(context.Background(), "pkg exit graceful", ExitGraceful)
	assert.Equal(t, int32(0), exitCode.Load())

	Exit(context.Background(), "pkg exit abnormal", ExitAbnormal, errors.New("err"))
	assert.Equal(t, int32(1), exitCode.Load())

	assertBufferContainsAll(t, buf, "pkg exit none", "pkg exit graceful", "pkg exit abnormal")
}

func TestWorker_EventSpecificLogger(t *testing.T) {
	defaultBuf := &bytes.Buffer{}
	customBuf := &bytes.Buffer{}
	defaultLogger := newJSONLogger(defaultBuf)
	customLogger := newJSONLogger(customBuf).With("scoped", "custom_val")
	w := newTestWorker(t, defaultLogger, nil)
	ok := w.Enqueue(&LogEvent{Logger: customLogger, Message: "event with specific logger",
		Level: slog.LevelInfo, Attrs: []slog.Attr{slog.String("extra", "field")}, Timestamp: time.Now()})
	assert.True(t, ok)
	w.Flush()
	assert.Empty(t, defaultBuf.String())
	assertCustomLoggerOutput(t, customBuf.String())
}

func assertCustomLoggerOutput(t *testing.T, out string) {
	t.Helper()
	assert.Contains(t, out, "event with specific logger")
	assert.Contains(t, out, `"scoped":"custom_val"`)
	assert.Contains(t, out, `"extra":"field"`)
}
