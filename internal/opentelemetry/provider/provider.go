// Package provider manages the OpenTelemetry TracerProvider lifecycle,
// OTLP trace exporters, resources, and global propagators.
package provider

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/callstack"
	"github.com/nawaphonOHM/whatever/v2/internal/opentelemetry/config"
	"github.com/nawaphonOHM/whatever/v2/pkg/logging"
)

// ShutdownFunc defines the callback function to flush and terminate the TracerProvider.
type ShutdownFunc = func(context.Context) error

// noopShutdown is returned when telemetry is disabled.
func noopShutdown(context.Context) error {
	return nil
}

// isTelemetryDisabled returns true if configuration is nil or explicitly disabled.
func isTelemetryDisabled(cfg *config.Config) bool {
	return cfg == nil || !cfg.Enabled
}

// validateConfig validates the OpenTelemetry configuration.
func validateConfig(ctx context.Context, cfg *config.Config) error {
	if err := cfg.Validate(); err != nil {
		logging.ErrorContext(ctx, "OpenTelemetry provider configuration validation failed", "error", err)
		return fmt.Errorf("invalid opentelemetry config: %w", err)
	}
	return nil
}

// logProviderConfig logs the provider configuration values.
func logProviderConfig(ctx context.Context, cfg *config.Config) {
	logging.InfoContext(ctx, "OpenTelemetry provider configuration choices",
		"service_name", cfg.ServiceName,
		"endpoint", cfg.Endpoint,
		"protocol", cfg.Protocol,
		"insecure", cfg.Insecure,
		"sample_rate", cfg.SampleRate,
	)
}

// initAndRegisterTracer registers the configured tracer provider globally.
func initAndRegisterTracer(
	ctx context.Context,
	exporter sdktrace.SpanExporter,
	res *resource.Resource,
	sampleRate float64,
) ShutdownFunc {
	logging.InfoContext(ctx, "OpenTelemetry sampling rate selected", "sample_rate", sampleRate)
	tp := buildTracerProvider(exporter, res, sampleRate)
	otel.SetTracerProvider(tp)
	logging.InfoContext(ctx, "OpenTelemetry provider initialization complete", "enabled", true)
	return tp.Shutdown
}

// setupActiveProvider validates configuration, initializes resources/exporters, and configures global tracer.
func setupActiveProvider(ctx context.Context, cfg *config.Config) (ShutdownFunc, error) {
	if err := validateConfig(ctx, cfg); err != nil {
		return nil, err
	}
	logProviderConfig(ctx, cfg)
	exporter, res, err := createExporterAndResource(ctx, cfg)
	if err != nil {
		logging.ErrorContext(ctx, "OpenTelemetry provider initialization failed", "error", err)
		return nil, err
	}
	return initAndRegisterTracer(ctx, exporter, res, cfg.SampleRate), nil
}

// InitTracerProvider initializes the OpenTelemetry TracerProvider, registers global propagators,
// and returns a graceful shutdown function.
func InitTracerProvider(ctx context.Context, cfg *config.Config) (ShutdownFunc, error) {
	decorated := callstack.DecorateContextFunc(
		"opentelemetry.provider.InitTracerProvider",
		func(ctx context.Context) (ShutdownFunc, error) {
			logging.InfoContext(ctx, "entering OpenTelemetry provider initialization state")
			RegisterPropagators()
			if isTelemetryDisabled(cfg) {
				logging.InfoContext(
					ctx,
					"OpenTelemetry disabled; falling back to noop tracer provider",
					"enabled",
					false,
				)
				logging.InfoContext(ctx, "OpenTelemetry provider initialization complete", "enabled", false)
				return noopShutdown, nil
			}
			return setupActiveProvider(ctx, cfg)
		},
	)
	return decorated(ctx)
}
