package central

import (
	"bytes"
	"context"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const proofWorkerQueueSize = 100

func setupProofWorker(t *testing.T, buf *bytes.Buffer, exitCode *atomic.Int32) *Worker {
	t.Helper()
	testWorker := New(newJSONLogger(buf), proofWorkerQueueSize, func(code int) {
		exitCode.Store(int32(code))
	})
	testWorker.Start()
	t.Cleanup(func() {
		require.NoError(t, testWorker.Close())
	})

	prev := DefaultWorker()
	SetDefaultWorker(testWorker)
	t.Cleanup(func() {
		SetDefaultWorker(prev)
	})
	return testWorker
}

func verifyWorkerBorrowing(t *testing.T, expected *Worker) *Worker {
	t.Helper()
	w1 := DefaultWorker()
	w2 := DefaultWorker()
	assert.Same(t, expected, w1)
	assert.Same(t, w1, w2)
	assert.True(t, w1.isRunning())
	return w1
}

func verifyWorkerEventDispatch(t *testing.T, buf *bytes.Buffer, exitCode *atomic.Int32, w *Worker) {
	t.Helper()
	w.Logger().InfoContext(context.Background(), "borrowed worker info", "component", "proof")
	ok := w.Enqueue(&LogEvent{
		Message:   "borrowed worker enqueue",
		Level:     slog.LevelInfo,
		Timestamp: time.Now(),
	})
	assert.True(t, ok)
	w.Flush()

	assert.Contains(t, buf.String(), "borrowed worker info")
	assert.Contains(t, buf.String(), "borrowed worker enqueue")
	assert.Equal(t, int32(-1), exitCode.Load())
	assert.True(t, w.isRunning())
}

func verifyPackageDispatch(t *testing.T, buf *bytes.Buffer, exitCode *atomic.Int32) {
	t.Helper()
	Dispatch(&LogEvent{
		Message:   "pkg dispatch borrowed",
		Level:     slog.LevelInfo,
		Timestamp: time.Now(),
	})
	FlushDefault()

	assert.Contains(t, buf.String(), "pkg dispatch borrowed")
	assert.Equal(t, int32(-1), exitCode.Load())
	assert.True(t, DefaultWorker().isRunning())
}

func TestDefaultWorker_SingletonResourceProof(t *testing.T) {
	// Callers across packages (internal/mongodb, internal/rest, internal/logging) borrow
	// DefaultWorker() without owning its lifecycle or closing it per operation.
	buf := &bytes.Buffer{}
	var exitCode atomic.Int32
	exitCode.Store(-1)

	testWorker := setupProofWorker(t, buf, &exitCode)
	w1 := verifyWorkerBorrowing(t, testWorker)
	verifyWorkerEventDispatch(t, buf, &exitCode, w1)
	verifyPackageDispatch(t, buf, &exitCode)
}
