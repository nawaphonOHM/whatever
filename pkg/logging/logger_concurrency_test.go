package logging

import (
	"bytes"
	"strings"
	"sync"
	"testing"

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

	l := New(Config{Output: syncedBuf, Level: LevelDebug})
	runConcurrentWorkers(l, testWorkerCount, testIterationCount)

	mu.Lock()
	lineCount := len(strings.Split(strings.TrimSpace(buf.String()), "\n"))
	mu.Unlock()

	assert.Equal(t, testExpectedLineCount, lineCount)
}
