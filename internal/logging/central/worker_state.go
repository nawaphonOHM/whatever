package central

import (
	"log/slog"
	"os"
)

// IsSync reports whether the worker runs synchronously.
func (w *Worker) IsSync() bool {
	if w == nil {
		return false
	}
	return w.syncMode
}

// SetExitFunc updates the exit function.
func (w *Worker) SetExitFunc(fn func(int)) {
	if w == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.exitFunc = resolveExitFunc(fn)
}

// ExitFunc returns the current exit function.
func (w *Worker) ExitFunc() func(int) {
	if w == nil {
		return os.Exit
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.exitFunc
}

// SetLogger updates the underlying slog.Logger.
func (w *Worker) SetLogger(l *slog.Logger) {
	if w == nil || l == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger = l
}

// Logger returns the underlying slog.Logger.
func (w *Worker) Logger() *slog.Logger {
	if w == nil {
		return slog.Default()
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.logger
}
