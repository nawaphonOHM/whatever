package provider

import (
	"context"
	"fmt"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
)

// newExporter constructs an OTLP SpanExporter according to the configured protocol.
func newExporter(ctx context.Context, cfg *config.Config) (sdktrace.SpanExporter, error) {
	if cfg.IsGRPC() {
		return newGRPCExporter(ctx, cfg)
	}
	if cfg.IsHTTP() {
		return newHTTPExporter(ctx, cfg)
	}
	return nil, fmt.Errorf("%w: %s", config.ErrInvalidProtocol, cfg.Protocol)
}
