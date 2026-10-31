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
)

type countWriter struct {
	delivered *atomic.Int64
}

func (cw *countWriter) Write(p []byte) (int, error) {
	cw.delivered.Add(1)
	return len(p), nil
}

func newCountingLogger(name string, counter *atomic.Int64) *slog.Logger {
	cw := &countWriter{delivered: counter}
	h := slog.NewJSONHandler(cw, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(h).With("name", name)
}

func verifyDynDefault(t *testing.T, dw *Worker, p2 *testLogPair) {
	t.Helper()
	slog.SetDefault(p2.logger)
	assert.Equal(t, p2.logger, dw.Logger())
	Dispatch(&LogEvent{Message: "msg dyn2", Level: slog.LevelInfo, Timestamp: time.Now()})
	FlushDefault()
	assert.Contains(t, p2.buf.String(), "msg dyn2")
}

func TestDefaultWorker_DynamicSlogDefault(t *testing.T) {
	restoreLoggingState(t)
	p1, p2 := newLogPair("dyn1"), newLogPair("dyn2")
	slog.SetDefault(p1.logger)
	dw := DefaultWorker()
	assert.Equal(t, p1.logger, dw.Logger())
	verifyDynDefault(t, dw, p2)
}

func setupExitWorker(t *testing.T, code *atomic.Int32) {
	t.Helper()
	code.Store(-1)
	origExit := DefaultWorker().ExitFunc()
	DefaultWorker().SetExitFunc(func(c int) { code.Store(int32(c)) })
	t.Cleanup(func() { DefaultWorker().SetExitFunc(origExit) })
}

func TestDefaultWorker_ExitWithDynamicDefault(t *testing.T) {
	restoreLoggingState(t)
	exitPair := newLogPair("exit-logger")
	slog.SetDefault(exitPair.logger)
	var exitCode atomic.Int32
	setupExitWorker(t, &exitCode)

	ExitWithAbnormal(context.Background(), "fatal abort", errors.New("err"), nil)
	assert.Equal(t, int32(1), exitCode.Load())
	assert.Contains(t, exitPair.buf.String(), "fatal abort")
	assert.Contains(t, exitPair.buf.String(), `"exit_flag":"ExitAbnormal"`)
}

func TestWorker_NilWorkerDelegation(t *testing.T) {
	var nilWorker *Worker
	assert.Equal(t, slog.Default(), nilWorker.Logger())
	assert.Equal(t, slog.Default(), nilWorker.eventLogger(nil))

	custom := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	assert.Equal(t, custom, nilWorker.eventLogger(custom))
}
