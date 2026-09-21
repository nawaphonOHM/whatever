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

func buildMissingTestCases() []missingTestCase {
	return []missingTestCase{
		{
			name:  "all required keys unset",
			setup: setupEmptyEnv,
			missingKeys: []string{
				testEnvURI, testEnvDatabase, testEnvUsername, testEnvPassword,
			},
		},
		{
			name: "missing URI only",
			setup: func(t *testing.T) {
				setupCustomEnv(t)
				t.Setenv(testEnvURI, "")
			},
			missingKeys: []string{testEnvURI},
		},
		{
			name: "missing Database only",
			setup: func(t *testing.T) {
				setupCustomEnv(t)
				t.Setenv(testEnvDatabase, "")
			},
			missingKeys: []string{testEnvDatabase},
		},
		{
			name: "missing Username only",
			setup: func(t *testing.T) {
				setupCustomEnv(t)
				t.Setenv(testEnvUsername, "")
			},
			missingKeys: []string{testEnvUsername},
		},
		{
			name: "missing Password only",
			setup: func(t *testing.T) {
				setupCustomEnv(t)
				t.Setenv(testEnvPassword, "")
			},
			missingKeys: []string{testEnvPassword},
		},
	}
}

// TestConfig_Load_MissingKeys_PeacefulExit tests missing required keys trigger peaceful exit.
func TestConfig_Load_MissingKeys_PeacefulExit(t *testing.T) {
	origExit := exitFunc
	defer func() { exitFunc = origExit }()

	for _, tc := range buildMissingTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t)
			var exitCalled bool
			var exitCode int
			exitFunc = func(code int) {
				exitCalled = true
				exitCode = code
			}

			cfg, err := LoadConfig()
			require.Error(t, err)
			assert.Nil(t, cfg)
			assert.True(t, exitCalled)
			assert.Equal(t, 0, exitCode)
			for _, key := range tc.missingKeys {
				assert.Contains(t, err.Error(), key)
			}
		})
	}
}
