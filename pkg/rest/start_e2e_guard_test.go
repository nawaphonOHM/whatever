package rest

import (
	"testing"
	"time"
)

const guardSettleDelay = 50 * time.Millisecond

// TestGuardedSIGTERM_SurvivesWithoutListener ensures signaling without a server listener is safe.
func TestGuardedSIGTERM_SurvivesWithoutListener(t *testing.T) {
	stop := guardedSIGTERM(t)
	time.Sleep(guardSettleDelay)
	stop()
}
