// Package central provides the asynchronous central logging worker.
package central

import (
	"context"
	"log/slog"
	"time"

	"github.com/nawaphonOHM/whatever/internal/logging/callstack"
)

// LogEvent represents an envelope dispatched to the Central Log worker.
type LogEvent struct {
	Timestamp time.Time
	Ctx       context.Context
	Logger    *slog.Logger
	Stack     *callstack.CallStack
	Err       error
	flushCh   chan struct{}
	Message   string
	Attrs     []slog.Attr
	Level     slog.Level
	ExitFlag  ExitFlag
}

// NewLogEvent creates a basic LogEvent.
func NewLogEvent(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) *LogEvent {
	return &LogEvent{
		Ctx:       ctx,
		Message:   msg,
		Level:     level,
		Attrs:     attrs,
		Timestamp: time.Now(),
	}
}
