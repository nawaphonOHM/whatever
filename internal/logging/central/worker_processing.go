package central

import (
	"context"
	"log/slog"
)

// processEvent formats and dispatches a single LogEvent to slog.
func (w *Worker) processEvent(event *LogEvent) {
	if event == nil {
		return
	}
	if event.flushCh != nil {
		close(event.flushCh)
		return
	}
	ctx := eventContext(event.Ctx)
	attrs := eventAttrs(event)
	targetLogger := w.eventLogger(event.Logger)
	targetLogger.LogAttrs(ctx, event.Level, event.Message, attrs...)
	w.handleExitEvent(ctx, event, targetLogger)
}

func (w *Worker) handleExitEvent(ctx context.Context, event *LogEvent, logger *slog.Logger) {
	if event == nil {
		return
	}
	w.dispatchExit(ctx, event, logger)
}

func (w *Worker) dispatchExit(ctx context.Context, event *LogEvent, logger *slog.Logger) {
	if event.ExitFlag == ExitAbnormal {
		w.handleAbnormalExit(ctx, event, logger)
		return
	}
	if event.ExitFlag == ExitGraceful {
		w.handleGracefulExit()
	}
}

func (w *Worker) handleAbnormalExit(ctx context.Context, event *LogEvent, logger *slog.Logger) {
	logAbnormalStack(ctx, event, logger)
	flushes := w.drainChannel()
	w.invokeExit(1)
	closeFlushChannels(flushes)
}

func (w *Worker) handleGracefulExit() {
	flushes := w.drainChannel()
	w.invokeExit(0)
	closeFlushChannels(flushes)
}

func (w *Worker) drainChannel() []chan struct{} {
	var flushes []chan struct{}
	for {
		ev, ok := w.nextQueuedEvent()
		if !ok {
			return flushes
		}
		flushes = w.appendFlush(flushes, ev)
	}
}

func (w *Worker) appendFlush(flushes []chan struct{}, event *LogEvent) []chan struct{} {
	if event == nil {
		return flushes
	}
	if event.flushCh != nil {
		return append(flushes, event.flushCh)
	}
	w.processEvent(event)
	return flushes
}

func closeFlushChannels(flushes []chan struct{}) {
	for _, flushCh := range flushes {
		close(flushCh)
	}
}

func (w *Worker) invokeExit(code int) {
	if exitFunc := w.ExitFunc(); exitFunc != nil {
		exitFunc(code)
	}
}
