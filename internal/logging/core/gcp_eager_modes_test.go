package core_test

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envModesHelper = "TEST_MODES_EAGER_HELPER"
	msgNonGCP      = "non-gcp plain log event"
	msgOtelActive  = "otel-active log event"
)

func runModesSubprocess(t *testing.T, extraEnv []string) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestModesEagerHelperProcess$", "-test.v")
	cmd.Env = append(os.Environ(), envModesHelper+"=1")
	cmd.Env = append(cmd.Env, extraEnv...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	require.NoError(t, cmd.Run())
	return out.String()
}

func TestModesEagerHelperProcess(t *testing.T) {
	if os.Getenv(envModesHelper) != "1" {
		t.Skip("not running as modes helper process")
	}
	if os.Getenv("OHM9996_OTEL_ENABLED") == "true" {
		slog.Info(msgOtelActive, "tracer", "active")
	} else {
		slog.Info(msgNonGCP, "mode", "standard")
	}
	core.Flush()
}

func TestEagerInit_NonGCPProviderProducesTextFormat(t *testing.T) {
	output := runModesSubprocess(t, []string{
		"OHM9996_CLOUD_PROVIDER=",
		"OHM9996_LOGGING_OUTPUT=stdout",
		"OHM9996_LOGGING_LEVEL=info",
	})
	assert.Contains(t, output, msgNonGCP)
	assert.Contains(t, output, "level=INFO")
	assert.NotContains(t, output, `"severity":`)
}

func TestEagerInit_OTelActiveDisablesGCPHandler(t *testing.T) {
	output := runModesSubprocess(t, []string{
		"OHM9996_CLOUD_PROVIDER=gcp",
		"OHM9996_OTEL_ENABLED=true",
		"OHM9996_LOGGING_OUTPUT=stdout",
		"OHM9996_LOGGING_LEVEL=info",
	})
	assert.Contains(t, output, msgOtelActive)
	assert.Contains(t, output, "level=INFO")
	assert.NotContains(t, output, `"severity":`)
}
