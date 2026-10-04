package central

// Enqueue adds a LogEvent to the worker channel.
// If in synchronous mode, buffer is full, or worker not started/closed, falls back to direct processing.
func (w *Worker) Enqueue(event *LogEvent) bool {
	if w == nil || event == nil {
		return false
	}
	return w.enqueueValid(event)
}

func (w *Worker) enqueueValid(event *LogEvent) bool {
	if w.shouldProcessDirectly() {
		w.processEvent(event)
		return w.syncMode
	}
	return w.tryEnqueue(event)
}

// Flush blocks until all currently enqueued log events have been processed.
func (w *Worker) Flush() {
	if w == nil {
		return
	}
	if !w.canFlush() {
		return
	}

	done := make(chan struct{})
	w.ch <- &LogEvent{flushCh: done}
	<-done
}

func (w *Worker) canFlush() bool {
	return !w.syncMode && w.isRunning()
}

// Stop closes the worker and waits for queued events to finish.
func (w *Worker) Stop() error {
	return w.Close()
}

// Close gracefully closes the worker and flushes pending events.
func (w *Worker) Close() error {
	if w == nil {
		return nil
	}
	if w.syncMode {
		return nil
	}
	return w.closeAsync()
}
