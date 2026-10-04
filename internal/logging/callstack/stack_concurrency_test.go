package callstack

import (
	"sync"
	"testing"
)

const (
	testWorkers    = 20
	testIterations = 50
)

func runStackWorker(cs *CallStack, workerID int) {
	for j := range testIterations {
		cs.Push(&Frame{Function: "concur", Line: workerID*testIterations + j})
		_ = cs.Frames()
		_ = cs.Len()
		_ = cs.String()
		_, _ = cs.Pop()
	}
}

func TestCallStack_Concurrency(t *testing.T) {
	t.Helper()
	cs := New()
	var wg sync.WaitGroup
	for i := range testWorkers {
		workerID := i
		wg.Go(func() {
			runStackWorker(cs, workerID)
		})
	}
	wg.Wait()
}
