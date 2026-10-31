package logging_test

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envPkgResetHelper = "TEST_PKG_LOGGING_RESET_HELPER"
	msgResetLogging   = "pkg logging reset message"
	msgResetSlog      = "pkg slog reset message"
)

func restorePkgLoggingState(t *testing.T) {
	t.Helper()
	prevWriter := log.Writer()
	prevFlags := log.Flags()

	t.Cleanup(func() {
		// Restore stdlib log settings first: ResetDefault then re-binds log to the fresh handler
		// via slog.SetDefault instead of leaving it pointing at the closed previous logger.
		log.SetOutput(prevWriter)
		log.SetFlags(prevFlags)
		logging.ResetDefault()
	})
}

func TestRestorePkgLoggingState_RebindsStdLogToCurrentDefault(t *testing.T) {
	restorePkgLoggingState(t)
	t.Setenv("OHM9996_LOGGING_LEVEL", "error")

	buf := new(bytes.Buffer)
	stale := core.NewJSON(buf, config.LevelDebug)
	core.SetDefault(stale)

	t.Run("inner", func(t *testing.T) {
		restorePkgLoggingState(t)
	})

	require.NotEqual(t, stale.Handler(), slog.Default().Handler())
	log.Print("probe after cleanup")
	stale.Flush()
	assert.Empty(t, buf.String(), "stdlib log must not route to the replaced logger")
}

func TestResetDefault_SynchronizesSlogAndAppliesNewEnv(t *testing.T) {
	restorePkgLoggingState(t)

	t.Setenv("OHM9996_CLOUD_PROVIDER", "gcp")
	t.Setenv("OHM9996_LOGGING_LEVEL", "info")
	t.Setenv("OHM9996_LOGGING_OUTPUT", "stdout")

	logging.ResetDefault()

	assert.Equal(t, core.Default().Handler(), slog.Default().Handler())
}

func runPkgResetSubprocess(t *testing.T) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestLoggingResetHelperProcess$", "-test.v")
	cmd.Env = append(os.Environ(),
		envPkgResetHelper+"=1",
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	require.NoError(t, cmd.Run())
	return out.String()
}

func mustSetEnvHelper(key, val string) {
	if err := os.Setenv(key, val); err != nil {
		panic(err)
	}
}

func TestLoggingResetHelperProcess(t *testing.T) {
	if os.Getenv(envPkgResetHelper) != "1" {
		t.Skip("not running as pkg reset helper process")
	}

	mustSetEnvHelper("OHM9996_CLOUD_PROVIDER", "gcp")
	mustSetEnvHelper("OHM9996_LOGGING_LEVEL", "info")
	mustSetEnvHelper("OHM9996_LOGGING_OUTPUT", "stdout")

	logging.ResetDefault()

	logging.Info(msgResetLogging, keyModule, "pkg_reset")
	slog.Info(msgResetSlog, keyModule, "slog_reset")
	logging.Flush()
}

func verifyResetEntries(t *testing.T, jsonLines []string) {
	t.Helper()
	require.Len(t, jsonLines, 2)
	entries := make([]map[string]any, len(jsonLines))
	for i, line := range jsonLines {
		require.NoError(t, json.Unmarshal([]byte(line), &entries[i]))
		assert.Equal(t, "INFO", entries[i]["severity"])
		assert.NotEmpty(t, entries[i]["timestamp"])
	}

	assert.True(t, findMatchingEntry(entries, msgResetLogging, "pkg_reset"))
	assert.True(t, findMatchingEntry(entries, msgResetSlog, "slog_reset"))
}

func TestLogging_ResetDefault_Subprocess(t *testing.T) {
	output := runPkgResetSubprocess(t)
	lines := extractJSONLogLines(output)
	verifyResetEntries(t, lines)
}
