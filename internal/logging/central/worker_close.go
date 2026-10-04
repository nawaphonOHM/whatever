package central

func (w *Worker) shouldProcessDirectly() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.syncMode || !w.started || w.closed
}

func (w *Worker) tryEnqueue(event *LogEvent) bool {
	select {
	case w.ch <- event:
		return true
	default:
		w.processEvent(event)
		return false
	}
}

func (w *Worker) isRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.started && !w.closed
}

func (w *Worker) closeAsync() error {
	if w.alreadyClosed() {
		return nil
	}
	started := w.markClosed()
	if started {
		w.wg.Wait()
	}
	return nil
}

func (w *Worker) alreadyClosed() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.closed
}

func (w *Worker) markClosed() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.closed = true
	if !w.started {
		return false
	}
	close(w.stopCh)
	return true
}
