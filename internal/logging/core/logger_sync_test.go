package core_test

import (
	"bytes"
	"log"
	"log/slog"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func restoreLoggingState(t *testing.T) {
	t.Helper()
	prevWriter := log.Writer()
	prevFlags := log.Flags()

	t.Cleanup(func() {
		// Restore stdlib log settings first: ResetDefault then re-binds log to the fresh handler
		// via slog.SetDefault instead of leaving it pointing at the closed previous logger.
		log.SetOutput(prevWriter)
		log.SetFlags(prevFlags)
		core.ResetDefault()
	})
}

func TestRestoreLoggingState_RebindsStdLogToCurrentDefault(t *testing.T) {
	restoreLoggingState(t)
	t.Setenv("OHM9996_LOGGING_LEVEL", "error")

	buf := new(bytes.Buffer)
	stale := core.NewJSON(buf, config.LevelDebug)
	core.SetDefault(stale)

	t.Run("inner", func(t *testing.T) {
		restoreLoggingState(t)
	})

	require.NotEqual(t, stale.Handler(), slog.Default().Handler())
	log.Print("probe after cleanup")
	stale.Flush()
	assert.Empty(t, buf.String(), "stdlib log must not route to the replaced logger")
}

func TestSetDefault_SynchronizesSlogDefault(t *testing.T) {
	restoreLoggingState(t)

	buf := new(bytes.Buffer)
	customLogger := core.NewJSON(buf, config.LevelDebug)
	core.SetDefault(customLogger)

	assert.Equal(t, customLogger, core.Default())
	assert.Equal(t, customLogger.Handler(), slog.Default().Handler())

	slog.Info("sync slog msg", "key", "val")
	customLogger.Flush()

	assert.Contains(t, buf.String(), "sync slog msg")
	assert.Contains(t, buf.String(), `"key":"val"`)
}

func TestResetDefault_SynchronizesSlogDefault(t *testing.T) {
	restoreLoggingState(t)

	buf := new(bytes.Buffer)
	customLogger := core.NewJSON(buf, config.LevelDebug)
	core.SetDefault(customLogger)
	require.Equal(t, customLogger, core.Default())

	core.ResetDefault()
	assert.NotEqual(t, customLogger, core.Default())
	assert.Equal(t, core.Default().Handler(), slog.Default().Handler())
}

func TestSetDefault_NilDoesNotOverwrite(t *testing.T) {
	restoreLoggingState(t)

	current := core.Default()
	core.SetDefault(nil)
	assert.Equal(t, current, core.Default())
	assert.Equal(t, current.Handler(), slog.Default().Handler())
}
