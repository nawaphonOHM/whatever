package provider

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/nawaphonOHM/whatever/v3/internal/opentelemetry/config"
)

// DefaultServiceVersion is the fallback service version if not specified elsewhere.
const DefaultServiceVersion = "1.0.0"

// buildResourceOptions creates options for resource detection and attributes.
func buildResourceOptions(serviceName string) []resource.Option {
	return []resource.Option{
		resource.WithHost(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(DefaultServiceVersion),
		),
	}
}

// isFatalResourceError returns true if err is not nil and not ErrPartialResource.
func isFatalResourceError(err error) bool {
	if err == nil {
		return false
	}
	return !errors.Is(err, resource.ErrPartialResource)
}

// newResource constructs an OpenTelemetry Resource with host metadata, telemetry SDK, and service attributes.
func newResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	opts := buildResourceOptions(serviceName)
	res, err := resource.New(ctx, opts...)
	if isFatalResourceError(err) {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	return res, nil
}

// buildSampler returns an appropriate Sampler based on the configured sample rate.
func buildSampler(rate float64) sdktrace.Sampler {
	if rate <= config.MinSampleRate {
		return sdktrace.ParentBased(sdktrace.NeverSample())
	}
	if rate >= config.MaxSampleRate {
		return sdktrace.ParentBased(sdktrace.AlwaysSample())
	}
	return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(rate))
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
