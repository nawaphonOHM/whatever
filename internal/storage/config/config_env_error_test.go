package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig_Load_ParseError verifies loader error when an environment variable fails type parsing.
func TestConfig_Load_ParseError(t *testing.T) {
	t.Setenv("OHM9996_STORAGE_BUCKET", "valid-bucket")
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
	t.Setenv("OHM9996_STORAGE_BUCKET", "valid-bucket")
	t.Setenv("OHM9996_STORAGE_PORT", "99999")

	state, restore := interceptExitHook()
	defer restore()

	cfg, err := LoadConfig()
	assert.Nil(t, cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid storage config")
	assert.False(t, state.exitCalled)
}
