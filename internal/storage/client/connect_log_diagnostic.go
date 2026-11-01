// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"log/slog"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

func isDebugLoggingEnabled(ctx context.Context, cfg *config.Config, logger *slog.Logger) bool {
	if cfg == nil {
		return false
	}
	return logger.Enabled(ctx, slog.LevelDebug)
}

func resolveNormalizedEndpoint(cfg *config.Config, detection detector.Detection) string {
	if cfg == nil {
		return ""
	}
	ep, err := resolveEndpointWithError(cfg, detection.Provider)
	if err != nil {
		return ""
	}
	return formatLoggedEndpoint(ep)
}

func logDetectionDiagnostics(
	ctx context.Context,
	cfg *config.Config,
	detection detector.Detection,
) {
	logger := central.DefaultWorker().Logger()
	if !isDebugLoggingEnabled(ctx, cfg, logger) {
		return
	}
	logger.DebugContext(ctx, "evaluating storage provider detection",
		"raw_endpoint", sanitizeLoggedURL(cfg.Endpoint),
		"raw_bucket_url", sanitizeLoggedURL(cfg.BucketURL),
		"inspected_host", detection.Host,
		"detected_provider", string(detection.Provider),
		"detection_rule", detection.Rule,
		"extracted_bucket", detection.Bucket,
		"normalized_endpoint", resolveNormalizedEndpoint(cfg, detection),
	)
}
