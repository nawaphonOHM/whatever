package provider

import (
	"bytes"
	"context"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/internal/opentelemetry/config"
)

type providerLogMode struct {
	capture  func(*bytes.Buffer) func()
	levelKey string
}

func jsonProviderLogMode() providerLogMode {
	return providerLogMode{
		capture:  func(output *bytes.Buffer) func() { return core.CaptureLogs(output) },
		levelKey: fieldLevel,
	}
}

func gcpProviderLogMode() providerLogMode {
	return providerLogMode{capture: captureGCPProviderLogs, levelKey: fieldSeverity}
}

func captureGCPProviderLogs(output *bytes.Buffer) func() {
	handler := core.NewGCPHandler(output, nil)
	return core.SetTestLogger(core.NewWithHandler(handler))
}

func captureProviderLogEntries(
	t *testing.T,
	cfg *config.Config,
	mode providerLogMode,
) ([]map[string]any, error) {
	t.Helper()
	output := new(bytes.Buffer)
	reset := mode.capture(output)
	defer reset()
	shutdown, err := InitTracerProvider(context.Background(), cfg)
	if err != nil {
		return parseLogEntries(t, output.Bytes()), err
	}
	if err := shutdown(context.Background()); err != nil {
		return parseLogEntries(t, output.Bytes()), err
	}
	return parseLogEntries(t, output.Bytes()), nil
}

func activeProviderTestConfig() *config.Config {
	return &config.Config{
		ServiceName: "test-service",
		Endpoint:    "127.0.0.1:4317",
		Protocol:    config.ProtocolGRPC,
		SampleRate:  testSampleRateHalf,
		Enabled:     true,
		Insecure:    true,
	}
}
