package central

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nawaphonOHM/whatever/internal/logging/callstack"
)

// Exit dispatches an exit event to the worker and flushes.
func (w *Worker) Exit(ctx context.Context, msg string, flag ExitFlag, args ...any) {
	w.Enqueue(newExitEvent(ctx, msg, flag, args...))
	w.Flush()
}

func newExitEvent(ctx context.Context, msg string, flag ExitFlag, args ...any) *LogEvent {
	stack, err := parseExitArgs(args)
	return &LogEvent{
		Ctx:       ctx,
		Message:   msg,
		Level:     exitLevel(flag),
		Err:       err,
		ExitFlag:  flag,
		Stack:     stack,
		Timestamp: time.Now(),
	}
}

func parseExitArgs(args []any) (*callstack.CallStack, error) {
	var err error
	var stack *callstack.CallStack
	for _, arg := range args {
		stack, err = parseExitArg(arg, stack, err)
	}
	return stack, err
}

func parseExitArg(arg any, stack *callstack.CallStack, err error) (*callstack.CallStack, error) {
	if value, ok := arg.(error); ok {
		return stack, value
	}
	if value, ok := arg.(*callstack.CallStack); ok {
		return value, err
	}
	return stack, err
}

func exitLevel(flag ExitFlag) slog.Level {
	if flag == ExitAbnormal {
		return slog.LevelError
	}
	return slog.LevelInfo
}

var (
	defaultWorker  atomic.Pointer[Worker]
	initWorkerOnce sync.Once
)

func initDefaultWorker() *Worker {
	w := New(slog.Default(), DefaultBufferSize, os.Exit)
	w.Start()
	return w
}

// DefaultWorker returns the package-level default worker.
func DefaultWorker() *Worker {
	if w := defaultWorker.Load(); w != nil {
		return w
	}
	initWorkerOnce.Do(func() {
		w := initDefaultWorker()
		defaultWorker.Store(w)
	})
	return defaultWorker.Load()
}

// SetDefaultWorker sets the package-level default worker.
func SetDefaultWorker(w *Worker) {
	if w != nil {
		defaultWorker.Store(w)
	}
}

// Dispatch enqueues an event to the default central worker.
func Dispatch(event *LogEvent) bool {
	// False positive: this wrapper borrows the singleton worker;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	return DefaultWorker().Enqueue(event)
}

// FlushDefault flushes all events in the default worker.
func FlushDefault() {
	// False positive: this wrapper borrows the singleton worker;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	DefaultWorker().Flush()
}

// Exit dispatches an exit event to the default central worker.
func Exit(ctx context.Context, msg string, flag ExitFlag, args ...any) {
	// False positive: this wrapper borrows the singleton worker;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	DefaultWorker().Exit(ctx, msg, flag, args...)
}

// ExitWithGraceful dispatches a graceful exit event to the default central worker.
func ExitWithGraceful(ctx context.Context, msg string) {
	// False positive: this wrapper borrows the singleton worker;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	DefaultWorker().Exit(ctx, msg, ExitGraceful, nil, nil)
}

// ExitWithAbnormal dispatches an abnormal exit event to the default central worker.
func ExitWithAbnormal(ctx context.Context, msg string, err error, stack *callstack.CallStack) {
	// False positive: this wrapper borrows the singleton worker;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	DefaultWorker().Exit(ctx, msg, ExitAbnormal, err, stack)
}
