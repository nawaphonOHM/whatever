package config

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	loggingconfig "github.com/nawaphonOHM/whatever/v2/internal/logging/config"
	"github.com/nawaphonOHM/whatever/v2/internal/logging/core"
)

func assertLifecycleLogs(t *testing.T, logs string) {
	t.Helper()
	for _, message := range []string{
		"entering OpenTelemetry configuration initialization state",
		"using default OpenTelemetry environment configuration resolution",
		"OpenTelemetry configuration choices",
		"OpenTelemetry configuration initialization complete",
	} {
		require.Contains(t, logs, message)
	}
	require.True(t, strings.Contains(logs, "protocol=grpc"))
}

func TestLoad_LogsLifecycleChoices(t *testing.T) {
	clearEnv(t)
	previous := core.Default()
	var output bytes.Buffer
	testLogger := core.NewText(&output, loggingconfig.LevelInfo)
	core.SetDefault(testLogger)
	t.Cleanup(func() {
		core.SetDefault(previous)
		require.NoError(t, testLogger.Close())
	})

	_, err := Load()
	require.NoError(t, err)
	testLogger.Flush()

	assertLifecycleLogs(t, output.String())
}
