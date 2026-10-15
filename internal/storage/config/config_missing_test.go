package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type exitHookState struct {
	exitCalled bool
	exitCode   int
}

func interceptExitHook() (*exitHookState, func()) {
	state := &exitHookState{}
	restore := SetExitFunc(func(code int) {
		state.exitCalled = true
		state.exitCode = code
	})
	return state, func() { SetExitFunc(restore) }
}

func assertMissingOutcome(t *testing.T, err error, state *exitHookState) {
	require.Error(t, err)
	assert.True(t, state.exitCalled)
	assert.Equal(t, 0, state.exitCode)
	assert.Contains(t, err.Error(), "OHM9996_STORAGE_BUCKET")
}

// TestConfig_Load_MissingBucket_PeacefulExit verifies missing required bucket triggers graceful exit.
func TestConfig_Load_MissingBucket_PeacefulExit(t *testing.T) {
	t.Setenv("OHM9996_STORAGE_BUCKET", "")
	state, restore := interceptExitHook()
	defer restore()

	cfg, err := LoadConfig()
	assert.Nil(t, cfg)
	assertMissingOutcome(t, err, state)
}
