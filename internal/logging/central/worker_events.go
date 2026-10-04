package central

import (
	"context"
	"log/slog"

	"github.com/nawaphonOHM/whatever/internal/logging/callstack"
)

func (w *Worker) handleExitEvent(ctx context.Context, event *LogEvent, logger *slog.Logger) {
	if event == nil {
		return
	}
	w.dispatchExit(ctx, event, logger)
}

func (w *Worker) dispatchExit(ctx context.Context, event *LogEvent, logger *slog.Logger) {
	switch event.ExitFlag {
	case ExitAbnormal:
		w.handleAbnormalExit(ctx, event, logger)
	case ExitGraceful:
		w.handleGracefulExit()
	}
}

func (w *Worker) handleAbnormalExit(ctx context.Context, event *LogEvent, logger *slog.Logger) {
	logAbnormalStack(ctx, event, logger)
	flushes := w.drainChannel()
	w.invokeExit(1)
	closeFlushChannels(flushes)
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
		logger.LogAttrs(ctx, slog.LevelError, "Call stack snapshot:\n"+stackText,
			slog.String("callstack", stackText))
	}
}

func (w *Worker) handleGracefulExit() {
	flushes := w.drainChannel()
	w.invokeExit(0)
	closeFlushChannels(flushes)
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
