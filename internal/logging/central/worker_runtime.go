package central

import (
	"context"
	"log/slog"
	"os"
)

func resolveLogger(logger *slog.Logger) *slog.Logger {
	if logger != nil {
		return logger
	}
	return slog.Default()
}

func resolveExitFunc(exitFunc func(int)) func(int) {
	if exitFunc != nil {
		return exitFunc
	}
	return os.Exit
}

func resolveBufferSize(bufferSize int) int {
	if bufferSize < 0 {
		return DefaultBufferSize
	}
	return bufferSize
}

func (w *Worker) canStart() bool {
	return !w.syncMode && !w.started && !w.closed
}

func (w *Worker) runLoop() {
	for {
		event, ok := w.waitForEvent()
		if !ok {
			return
		}
		w.processEvent(event)
	}
}

func (w *Worker) waitForEvent() (*LogEvent, bool) {
	select {
	case event, ok := <-w.ch:
		return event, ok
	case <-w.stopCh:
		w.drainPendingEvents()
		return nil, false
	}
}

func (w *Worker) drainPendingEvents() {
	for len(w.ch) > 0 {
		w.processEvent(<-w.ch)
	}
}

func eventContext(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return context.Background()
}

func extraAttrs(event *LogEvent) []slog.Attr {
	var extras []slog.Attr
	if event.Err != nil {
		extras = append(extras, slog.Any("error", event.Err))
	}
	if event.ExitFlag != ExitNone {
		extras = append(extras, slog.String("exit_flag", event.ExitFlag.String()))
	}
	return extras
}

func eventAttrs(event *LogEvent) []slog.Attr {
	if event == nil {
		return nil
	}
	extras := extraAttrs(event)
	attrs := make([]slog.Attr, 0, len(event.Attrs)+len(extras))
	attrs = append(attrs, event.Attrs...)
	attrs = append(attrs, extras...)
	return attrs
}

func (w *Worker) eventLogger(logger *slog.Logger) *slog.Logger {
	if logger != nil {
		return logger
	}
	return w.logger
}

func (w *Worker) nextQueuedEvent() (*LogEvent, bool) {
	select {
	case event := <-w.ch:
		return event, true
	default:
		return nil, false
	}
}
