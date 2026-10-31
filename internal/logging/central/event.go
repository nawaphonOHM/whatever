// Package central provides the asynchronous central logging worker.
package central

import (
	"context"
	"log/slog"
	"time"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/callstack"
)

// ExitFlag determines post-logging termination behavior.
type ExitFlag int

const (
	ExitNone ExitFlag = iota
	ExitGraceful
	ExitAbnormal
)

var exitFlagNames = map[ExitFlag]string{
	ExitNone:     "ExitNone",
	ExitGraceful: "ExitGraceful",
	ExitAbnormal: "ExitAbnormal",
}

func (f ExitFlag) String() string {
	if name, ok := exitFlagNames[f]; ok {
		return name
	}
	return "ExitUnknown"
}

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

func hasStackFrames(stack *callstack.CallStack) bool {
	return stack != nil && stack.Len() > 0
}

func logAbnormalStack(ctx context.Context, event *LogEvent, logger *slog.Logger) {
	if logger == nil {
		return
	}
	stack := resolveEventStack(event)
	if hasStackFrames(stack) {
		stackText := stack.String()
		logger.LogAttrs(ctx, slog.LevelError, "Call stack snapshot",
			slog.String("callstack", stackText))
	}
}

func resolveEventStack(event *LogEvent) *callstack.CallStack {
	if event == nil {
		return callstack.DefaultStack()
	}
	if event.Stack != nil {
		return event.Stack
	}
	return resolveContextStack(event.Ctx)
}

func resolveContextStack(ctx context.Context) *callstack.CallStack {
	if ctx == nil {
		return callstack.DefaultStack()
	}
	if stack := callstack.FromContext(ctx); stack != nil {
		return stack
	}
	return callstack.DefaultStack()
}
