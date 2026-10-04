package core

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/config"
	"github.com/stretchr/testify/assert"
)

const (
	testWorkerCount       = 10
	testIterationCount    = 20
	testExpectedLineCount = 200
)

type syncedWriter struct {
	mu  *sync.Mutex
	buf *bytes.Buffer
}

func (s *syncedWriter) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func spawnWorker(l *Logger, id, iterations int, done chan<- struct{}) {
	for j := range iterations {
		l.Debug("concurrent log", "worker", id, "iteration", j)
	}
	done <- struct{}{}
}

func runConcurrentWorkers(l *Logger, count, iterations int) {
	done := make(chan struct{}, count)
	for i := range count {
		go spawnWorker(l, i, iterations, done)
	}
	for range count {
		<-done
	}
}

func TestLogger_ConcurrentLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	var mu sync.Mutex
	syncedBuf := &syncedWriter{mu: &mu, buf: buf}

	l := NewJSON(syncedBuf, config.LevelDebug)
	runConcurrentWorkers(l, testWorkerCount, testIterationCount)
	assert.Equal(t, testExpectedLineCount, countSyncedLines(&mu, buf))
}

func TestLogger_ConcurrentDefaultAccess(t *testing.T) {
	const count = 50
	var wg sync.WaitGroup
	for range count {
		wg.Go(func() {
			l := Default()
			assert.NotNil(t, l)
			assert.NotNil(t, l.Slog())
		})
	}
	wg.Wait()
}

func runConcurrentSetDefaultWorkers(wg *sync.WaitGroup, count int) {
	for i := range count {
		id := i
		wg.Go(func() {
			Info("concurrent global msg", "id", id)
			_ = Default()
		})
	}
}

func countSyncedLines(mu *sync.Mutex, buf *bytes.Buffer) int {
	mu.Lock()
	defer mu.Unlock()
	return len(strings.Split(strings.TrimSpace(buf.String()), "\n"))
}

func TestLogger_ConcurrentSetDefaultAndLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	var mu sync.Mutex
	syncedBuf := &syncedWriter{mu: &mu, buf: buf}
	cleanup := SetTestLogger(NewJSON(syncedBuf, config.LevelDebug))
	defer cleanup()

	const count = 20
	var wg sync.WaitGroup
	runConcurrentSetDefaultWorkers(&wg, count)
	wg.Wait()
	assert.Equal(t, count, countSyncedLines(&mu, buf))
}
