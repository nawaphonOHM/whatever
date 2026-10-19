package central

import (
	"log/slog"
	"os"
	"sync"
)

// DefaultBufferSize is the default capacity for the central log channel.
const DefaultBufferSize = 4096

// Worker is the central background worker goroutine that processes log events.
type Worker struct {
	ch       chan *LogEvent
	logger   *slog.Logger
	exitFunc func(int)
	stopCh   chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
	syncMode bool
	started  bool
	closed   bool
}

// New creates a new central log Worker.
// If bufferSize == 0, the worker operates in synchronous mode.
// If bufferSize > 0, the worker operates in asynchronous background mode.
func New(logger *slog.Logger, bufferSize int, exitFunc func(int)) *Worker {
	resolvedLogger := resolveLogger(logger)
	resolvedExit := resolveExitFunc(exitFunc)
	if bufferSize == 0 {
		return &Worker{
			logger:   resolvedLogger,
			exitFunc: resolvedExit,
			syncMode: true,
		}
	}
	resolvedBuffer := resolveBufferSize(bufferSize)
	return &Worker{
		ch:       make(chan *LogEvent, resolvedBuffer),
		logger:   resolvedLogger,
		exitFunc: resolvedExit,
		stopCh:   make(chan struct{}),
		syncMode: false,
	}
}

// Start launches the dedicated background goroutine named "central log".
func (w *Worker) Start() {
	if w == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.canStart() {
		return
	}
	w.started = true
	w.wg.Add(1)
	go w.centralLogLoop()
}

// centralLogLoop runs the central log worker background loop.
func (w *Worker) centralLogLoop() {
	defer w.wg.Done()
	w.runLoop()
}

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
