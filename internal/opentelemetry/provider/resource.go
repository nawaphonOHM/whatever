package provider

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
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
