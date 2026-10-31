package central

import (
	"bytes"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testLogPair struct {
	buf    *bytes.Buffer
	logger *slog.Logger
}

func newLogPair(name string) *testLogPair {
	buf := &bytes.Buffer{}
	h := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return &testLogPair{buf: buf, logger: slog.New(h).With("name", name)}
}

func enqueueNamedEvent(w *Worker, msg string) bool {
	return w.Enqueue(&LogEvent{Message: msg, Level: slog.LevelInfo, Timestamp: time.Now()})
}

func restoreLoggingState(t *testing.T) {
	t.Helper()
	origDefault := slog.Default()
	origWriter := log.Writer()
	origFlags := log.Flags()
	t.Cleanup(func() {
		slog.SetDefault(origDefault)
		log.SetOutput(origWriter)
		log.SetFlags(origFlags)
	})
}

func verifyDynamicSwitch(t *testing.T, w *Worker, p1, p2 *testLogPair) {
	t.Helper()
	slog.SetDefault(p2.logger)
	assert.Equal(t, p2.logger, w.Logger())
	assert.True(t, enqueueNamedEvent(w, "event2"))
	w.Flush()
	assert.Contains(t, p2.buf.String(), "event2")
	assert.NotContains(t, p1.buf.String(), "event2")
}

func setupDynamicTest(t *testing.T, p1 *testLogPair) *Worker {
	t.Helper()
	restoreLoggingState(t)
	slog.SetDefault(p1.logger)
	return newTestWorker(t, nil, nil)
}

func TestWorker_DynamicDefaultLoggerDelegation(t *testing.T) {
	p1, p2 := newLogPair("l1"), newLogPair("l2")
	w := setupDynamicTest(t, p1)

	assert.Equal(t, p1.logger, w.Logger())
	assert.True(t, enqueueNamedEvent(w, "event1"))
	w.Flush()
	assert.Contains(t, p1.buf.String(), "event1")
	verifyDynamicSwitch(t, w, p1, p2)
}

func TestWorker_CustomLoggerPrecedence(t *testing.T) {
	restoreLoggingState(t)
	defPair, custPair := newLogPair("default"), newLogPair("custom")
	slog.SetDefault(defPair.logger)

	w := newTestWorker(t, custPair.logger, nil)
	assert.Equal(t, custPair.logger, w.Logger())
	assert.True(t, enqueueNamedEvent(w, "custom event"))
	w.Flush()

	assert.Contains(t, custPair.buf.String(), "custom event")
	assert.Empty(t, defPair.buf.String())
}

func TestWorker_SetLoggerOverridesDynamic(t *testing.T) {
	w := newTestWorker(t, nil, nil)
	custPair := newLogPair("override")
	w.SetLogger(custPair.logger)

	assert.Equal(t, custPair.logger, w.Logger())
	assert.True(t, enqueueNamedEvent(w, "override event"))
	w.Flush()
	assert.Contains(t, custPair.buf.String(), "override event")
}
