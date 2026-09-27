// Package provider manages the OpenTelemetry TracerProvider lifecycle,
// OTLP trace exporters, resources, and global propagators.
package provider

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
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

// buildTracerProvider constructs a new SDK TracerProvider with batch processor, resource, and sampler.
func buildTracerProvider(
	exporter sdktrace.SpanExporter,
	res *resource.Resource,
	rate float64,
) *sdktrace.TracerProvider {
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	sampler := buildSampler(rate)
	return sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)
}

// createExporterAndResource initializes the configured exporter and resource.
func createExporterAndResource(
	ctx context.Context,
	cfg *config.Config,
) (sdktrace.SpanExporter, *resource.Resource, error) {
	exporter, err := newExporter(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize exporter: %w", err)
	}
	res, err := newResource(ctx, cfg.ServiceName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}
	return exporter, res, nil
}

// setupActiveProvider validates configuration, initializes resources/exporters, and configures global tracer.
func setupActiveProvider(ctx context.Context, cfg *config.Config) (ShutdownFunc, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid opentelemetry config: %w", err)
	}
	exporter, res, err := createExporterAndResource(ctx, cfg)
	if err != nil {
		return nil, err
	}
	tp := buildTracerProvider(exporter, res, cfg.SampleRate)
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

// InitTracerProvider initializes the OpenTelemetry TracerProvider, registers global propagators,
// and returns a graceful shutdown function.
func InitTracerProvider(ctx context.Context, cfg *config.Config) (ShutdownFunc, error) {
	RegisterPropagators()
	if isTelemetryDisabled(cfg) {
		return noopShutdown, nil
	}
	return setupActiveProvider(ctx, cfg)
}
