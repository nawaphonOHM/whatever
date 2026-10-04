package central

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
