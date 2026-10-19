package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const envStorageBucket = "OHM9996_STORAGE_BUCKET"

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
	assert.Contains(t, err.Error(), envStorageBucket)
}

// TestConfig_Load_MissingBucket_PeacefulExit verifies missing required bucket triggers graceful exit.
func TestConfig_Load_MissingBucket_PeacefulExit(t *testing.T) {
	t.Setenv(envStorageBucket, "")
	state, restore := interceptExitHook()
	defer restore()

	cfg, err := LoadConfig()
	assert.Nil(t, cfg)
	assertMissingOutcome(t, err, state)
}

// TestConfig_Load_ParseError verifies loader error when an environment variable fails type parsing.
func TestConfig_Load_ParseError(t *testing.T) {
	t.Setenv(envStorageBucket, "valid-bucket")
	t.Setenv("OHM9996_STORAGE_PORT", "not-a-number")

	state, restore := interceptExitHook()
	defer restore()

	cfg, err := LoadConfig()
	assert.Nil(t, cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load storage config")
	assert.False(t, state.exitCalled)
}

// TestConfig_Load_ValidationError verifies error wrapping when parsed settings fail semantic validation.
func TestConfig_Load_ValidationError(t *testing.T) {
	t.Setenv(envStorageBucket, "valid-bucket")
	t.Setenv("OHM9996_STORAGE_PORT", "99999")

	state, restore := interceptExitHook()
	defer restore()

	cfg, err := LoadConfig()
	assert.Nil(t, cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid storage config")
	assert.False(t, state.exitCalled)
}
