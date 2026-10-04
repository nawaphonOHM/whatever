package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type missingTestCase struct {
	name        string
	setup       func(*testing.T)
	missingKeys []string
}

func buildMissingTestCases() []*missingTestCase {
	return []*missingTestCase{
		{
			name:  "all required keys unset",
			setup: setupEmptyEnv,
			missingKeys: []string{
				testEnvHost,
			},
		},
		{
			name: "missing Host only",
			setup: func(t *testing.T) {
				setupCustomEnv(t)
				t.Setenv(testEnvHost, "")
			},
			missingKeys: []string{testEnvHost},
		},
	}
}

type missingExitState struct {
	exitCalled bool
	exitCode   int
}

func interceptMissingExitHook() (*missingExitState, func()) {
	state := &missingExitState{}
	restore := SetExitFunc(func(code int) {
		state.exitCalled = true
		state.exitCode = code
	})
	return state, func() { SetExitFunc(restore) }
}

func assertMissingResult(t *testing.T, err error, missingKeys []string) {
	for _, key := range missingKeys {
		assert.Contains(t, err.Error(), key)
	}
}

func assertMissingOutcome(t *testing.T, err error, state *missingExitState, tc *missingTestCase) {
	require.Error(t, err)
	assert.True(t, state.exitCalled)
	assert.Equal(t, 0, state.exitCode)
	assertMissingResult(t, err, tc.missingKeys)
}

func runMissingTestCase(t *testing.T, tc *missingTestCase) {
	if tc == nil {
		return
	}
	tc.setup(t)
	state, restore := interceptMissingExitHook()
	defer restore()

	cfg, err := LoadConfig()
	assert.Nil(t, cfg)
	assertMissingOutcome(t, err, state, tc)
}

// TestConfig_Load_MissingKeys_PeacefulExit tests missing required keys trigger peaceful exit.
func TestConfig_Load_MissingKeys_PeacefulExit(t *testing.T) {
	for _, tc := range buildMissingTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			runMissingTestCase(t, tc)
		})
	}
}
