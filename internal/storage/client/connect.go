// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/callstack"
	_ "github.com/nawaphonOHM/whatever/v3/internal/logging/core" // side-effect: eagerly initialize structured logging
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

const (
	errInvalidConfigFormat  = "invalid storage config: %w"
	errDetectProviderFormat = "failed to detect storage provider: %w"
)

// Compile-time assertion that Client implements contracts.Storage.
var _ contracts.Storage = (*Client)(nil)

// Connect loads storage configuration from environment variables (OHM9996_STORAGE_*),
// establishes an S3 client, optionally probes bucket connectivity, and returns a managed Client.
func Connect(ctx context.Context, opts ...Option) (*Client, error) {
	fn := callstack.DecorateContextFunc("Connect", func(ctx context.Context) (*Client, error) {
		cfg, err := config.LoadConfig()
		if err != nil {
			return nil, err
		}
		return ConnectWithConfig(ctx, cfg, opts...)
	})
	return fn(ctx)
}

func resolveAndLogDetection(ctx context.Context, cfg *config.Config) (detector.Detection, error) {
	if cfg == nil {
		return detector.Detection{}, contracts.ErrNilConfig
	}
	detection, err := detector.Detect(cfg.Endpoint, cfg.BucketURL)
	if err != nil {
		return detector.Detection{}, fmt.Errorf(errDetectProviderFormat, err)
	}
	logDetectionDiagnostics(ctx, cfg, detection)
	return detection, nil
}

// ConnectWithConfig initializes a managed storage Client using the provided Config and Options.
func ConnectWithConfig(ctx context.Context, cfg *config.Config, opts ...Option) (*Client, error) {
	fn := callstack.DecorateContextFunc("ConnectWithConfig", func(ctx context.Context) (*Client, error) {
		detection, err := resolveAndLogDetection(ctx, cfg)
		if err != nil {
			return nil, err
		}
		if err := validateConfig(cfg, detection.Provider); err != nil {
			return nil, err
		}
		return initAndProbeClient(ctx, cfg, detection.Provider, opts...)
	})
	return fn(ctx)
}

func initAndProbeClient(
	ctx context.Context,
	cfg *config.Config,
	provider contracts.Provider,
	opts ...Option,
) (*Client, error) {
	options := NewOptions(opts...)
	c, err := buildClientInstance(ctx, cfg, provider, options)
	if err != nil {
		return nil, err
	}
	s3Opts := c.RawClient().Options()
	logClientInitAndCompat(ctx, compatLogParams{
		cfg:      cfg,
		provider: provider,
		opts:     options,
		s3Opts:   s3Opts,
	})
	if err := maybeProbeStartup(ctx, c, options); err != nil {
		return nil, err
	}
	logStartupComplete(ctx, lifecycleLogParams{
		provider: provider,
		bucket:   cfg.Bucket,
		s3Opts:   s3Opts,
	})
	return c, nil
}
