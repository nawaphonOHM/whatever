package core_test

import (
	"bytes"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	concurrentGoroutines = 10
	logsPerGoroutine     = 20
)

func runWorkerLogs(logger *slog.Logger, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for j := range logsPerGoroutine {
		logger.Info("concurrent message", "worker_id", id, "iteration", j)
	}
}

func launchWorkers(logger *slog.Logger) {
	var wg sync.WaitGroup
	for i := range concurrentGoroutines {
		wg.Add(1)
		go runWorkerLogs(logger, &wg, i)
	}
	wg.Wait()
}

func verifyLogLines(t *testing.T, output string) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedCount := concurrentGoroutines * logsPerGoroutine
	require.Len(t, lines, expectedCount)
	for _, line := range lines {
		entry := parseLogMap(t, []byte(line))
		assert.Equal(t, core.GCPSeverityInfo, entry[testKeySeverity])
	}
}

func TestGCPHandlerConcurrentWrites(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, nil)
	logger := slog.New(handler)

	launchWorkers(logger)
	verifyLogLines(t, buf.String())
}
