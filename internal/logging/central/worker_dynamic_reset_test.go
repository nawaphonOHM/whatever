package central

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorker_SetLoggerNil(t *testing.T) {
	restoreLoggingState(t)
	w := newTestWorker(t, nil, nil)
	custPair := newLogPair("override")
	w.SetLogger(custPair.logger)
	assert.Equal(t, custPair.logger, w.Logger())

	p := newLogPair("dynamic-after-nil")
	slog.SetDefault(p.logger)
	w.SetLogger(nil)
	assert.Equal(t, p.logger, w.Logger())
}

func TestWorker_ResetLogger(t *testing.T) {
	restoreLoggingState(t)
	w := newTestWorker(t, nil, nil)
	custPair := newLogPair("override")
	w.SetLogger(custPair.logger)
	assert.Equal(t, custPair.logger, w.Logger())

	p := newLogPair("dynamic-after-reset")
	slog.SetDefault(p.logger)
	w.ResetLogger()
	assert.Equal(t, p.logger, w.Logger())
}

func TestDefaultWorker_ResetDefaultLogger(t *testing.T) {
	restoreLoggingState(t)
	custPair := newLogPair("cust-default")
	DefaultWorker().SetLogger(custPair.logger)
	assert.Equal(t, custPair.logger, DefaultWorker().Logger())

	dynPair := newLogPair("dyn-default")
	slog.SetDefault(dynPair.logger)
	ResetDefaultLogger()
	assert.Equal(t, dynPair.logger, DefaultWorker().Logger())
}

func TestWorker_EventLoggerPrecedence(t *testing.T) {
	restoreLoggingState(t)
	defPair, evPair := newLogPair("default"), newLogPair("scoped")
	slog.SetDefault(defPair.logger)

	w := newTestWorker(t, nil, nil)
	ok := w.Enqueue(&LogEvent{
		Logger:    evPair.logger,
		Message:   "scoped event",
		Level:     slog.LevelInfo,
		Timestamp: time.Now(),
	})
	assert.True(t, ok)
	w.Flush()
	assert.Contains(t, evPair.buf.String(), "scoped event")
	assert.Empty(t, defPair.buf.String())
}
