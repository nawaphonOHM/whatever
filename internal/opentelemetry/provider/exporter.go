package provider

import (
	"context"
	"fmt"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
	"github.com/nawaphonOHM/whatever/pkg/logging"
)

// newExporter constructs an OTLP SpanExporter according to the configured protocol.
func newExporter(ctx context.Context, cfg *config.Config) (sdktrace.SpanExporter, error) {
	if cfg.IsGRPC() {
		logging.InfoContext(ctx, "OpenTelemetry exporter choice",
			"protocol", config.ProtocolGRPC,
			"transport", "grpc",
			"endpoint", cfg.Endpoint,
		)
		return newGRPCExporter(ctx, cfg)
	}
	if cfg.IsHTTP() {
		logging.InfoContext(ctx, "OpenTelemetry exporter choice",
			"protocol", cfg.Protocol,
			"transport", "http",
			"endpoint", cfg.Endpoint,
		)
		return newHTTPExporter(ctx, cfg)
	}
	return nil, fmt.Errorf("%w: %s", config.ErrInvalidProtocol, cfg.Protocol)
}
