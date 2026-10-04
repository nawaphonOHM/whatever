package provider

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nawaphonOHM/whatever/v2/internal/opentelemetry/config"
)

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
