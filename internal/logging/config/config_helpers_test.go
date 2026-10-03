package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testEnvOutput                  = "OHM9996_LOGGING_OUTPUT"
	testEnvLevel                   = "OHM9996_LOGGING_LEVEL"
	testEnvFormat                  = "OHM9996_LOGGING_FORMAT"
	testEnvAddSource               = "OHM9996_LOGGING_ADD_SOURCE"
	testEnvDisableTraceCorrelation = "OHM9996_LOGGING_DISABLE_TRACE_CORRELATION"

	testDotEnvPerm = 0o600

	testStderr  = "stderr"
	testJSON    = "json"
	testDiscard = "discard"
	testTrace   = "TRACE"
	testDebug   = "DEBUG"
	testInfo    = "INFO"
	testWarn    = "WARN"
	testError   = "ERROR"
	testFatal   = "FATAL"
)

// unsetEnv unsets an environment variable and restores it during test cleanup.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	if val, ok := os.LookupEnv(key); ok {
		require.NoError(t, os.Unsetenv(key))
		t.Cleanup(func() {
			require.NoError(t, os.Setenv(key, val))
		})
	}
}

// clearEnv clears all logging environment variables.
func clearEnv(t *testing.T) {
	t.Helper()
	keys := []string{
		testEnvOutput,
		testEnvLevel,
		testEnvFormat,
		testEnvAddSource,
		testEnvDisableTraceCorrelation,
	}
	for _, k := range keys {
		unsetEnv(t, k)
	}
}
