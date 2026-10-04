package provider

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	loggingconfig "github.com/nawaphonOHM/whatever/internal/logging/config"
	"github.com/nawaphonOHM/whatever/internal/logging/core"
	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
)

func captureProviderLogs(t *testing.T, fn func()) string {
	t.Helper()
	previous := core.Default()
	var output bytes.Buffer
	testLogger := core.NewText(&output, loggingconfig.LevelInfo)
	core.SetDefault(testLogger)
	t.Cleanup(func() {
		core.SetDefault(previous)
		require.NoError(t, testLogger.Close())
	})

	fn()
	testLogger.Flush()
	return output.String()
}

func TestInitTracerProvider_LogsDisabledLifecycle(t *testing.T) {
	logs := captureProviderLogs(t, func() {
		shutdown, err := InitTracerProvider(context.Background(), &config.Config{Enabled: false})
		require.NoError(t, err)
		require.NoError(t, shutdown(context.Background()))
	})

	for _, message := range []string{
		"entering OpenTelemetry provider initialization state",
		"OpenTelemetry disabled; falling back to noop tracer provider",
		"OpenTelemetry provider initialization complete",
	} {
		require.Contains(t, logs, message)
	}
}

func TestInitTracerProvider_LogsActiveChoices(t *testing.T) {
	logs := captureProviderLogs(t, func() {
		shutdown, err := InitTracerProvider(context.Background(), &config.Config{
			ServiceName: "lifecycle-service",
			Endpoint:    "127.0.0.1:4317",
			Protocol:    config.ProtocolGRPC,
			SampleRate:  0.5,
			Enabled:     true,
			Insecure:    true,
		})
		require.NoError(t, err)
		require.NoError(t, shutdown(context.Background()))
	})

	for _, message := range []string{
		"OpenTelemetry provider configuration choices",
		"OpenTelemetry exporter choice",
		"OpenTelemetry sampling rate selected",
		"OpenTelemetry provider initialization complete",
	} {
		require.Contains(t, logs, message)
	}
	require.True(t, strings.Contains(logs, "protocol=grpc"))
	require.True(t, strings.Contains(logs, "sample_rate=0.5"))
}
