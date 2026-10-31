package central

import "log/slog"

// SetLogger updates the underlying slog.Logger.
// If l is nil, the worker reverts to dynamic slog.Default() delegation.
func (w *Worker) SetLogger(l *slog.Logger) {
	if w == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger = l
}

// ResetLogger clears any explicit logger so the worker dynamically delegates to slog.Default().
func (w *Worker) ResetLogger() {
	if w == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger = nil
}

// Logger returns the underlying slog.Logger, or slog.Default() if no custom logger is configured.
func (w *Worker) Logger() *slog.Logger {
	if w == nil {
		return slog.Default()
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.logger != nil {
		return w.logger
	}
	return slog.Default()
}
