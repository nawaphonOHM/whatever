package central

import (
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testConcurrentWorkers = 8
	testWorkerIterations  = 20
	testModInterval       = 5
)

func runDynamicWorkerRoutine(w *Worker, workerID int, counter *atomic.Int64) {
	for j := range testWorkerIterations {
		if j%testModInterval == 0 {
			slog.SetDefault(newCountingLogger(fmt.Sprintf("w-%d-%d", workerID, j), counter))
		}
		_ = w.Enqueue(&LogEvent{Message: "concurrent log", Level: slog.LevelDebug, Timestamp: time.Now()})
	}
}

func runConcurrentDynamicWorkers(w *Worker, counter *atomic.Int64) {
	var wg sync.WaitGroup
	for i := range testConcurrentWorkers {
		workerID := i
		wg.Go(func() { runDynamicWorkerRoutine(w, workerID, counter) })
	}
	wg.Wait()
}

func TestWorker_ConcurrentDynamicAccess(t *testing.T) {
	restoreLoggingState(t)
	var counter atomic.Int64
	slog.SetDefault(newCountingLogger("init", &counter))

	w := newTestWorker(t, nil, nil)
	runConcurrentDynamicWorkers(w, &counter)
	w.Flush()
	assert.Equal(t, int64(testConcurrentWorkers*testWorkerIterations), counter.Load())
}
