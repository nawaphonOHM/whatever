// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/callstack"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
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

// SetExitFunc overrides the process exit hook used by the central logger worker
// and returns the previous hook.
func SetExitFunc(fn func(int)) func(int) {
	// False positive: borrowing the central worker singleton to configure exit hook is leak-free;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	prev := central.DefaultWorker().ExitFunc()
	central.DefaultWorker().SetExitFunc(fn)
	return prev
}

func resolveProvider(cfg *config.Config) (contracts.Provider, error) {
	if cfg == nil {
		return "", contracts.ErrNilConfig
	}
	provider, err := detector.DetectProvider(cfg.Endpoint, cfg.BucketURL)
	if err != nil {
		return "", fmt.Errorf(errDetectProviderFormat, err)
	}
	return provider, nil
}

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

// ConnectWithConfig initializes a managed storage Client using the provided Config and Options.
func ConnectWithConfig(ctx context.Context, cfg *config.Config, opts ...Option) (*Client, error) {
	fn := callstack.DecorateContextFunc("ConnectWithConfig", func(ctx context.Context) (*Client, error) {
		provider, err := resolveProvider(cfg)
		if err != nil {
			return nil, err
		}
		if err := validateAndLogConfig(ctx, cfg, provider); err != nil {
			return nil, err
		}
		return initAndProbeClient(ctx, cfg, provider, opts...)
	})
	return fn(ctx)
}

func handleConnectionError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	msg := "failed to connect to storage; exiting abnormally"
	central.ExitWithAbnormal(ctx, msg, err, nil)
	return err
}

func probeStartupClient(ctx context.Context, c *Client) error {
	return handleConnectionError(ctx, c.probeBucket(ctx))
}

func logStartupComplete(ctx context.Context, provider contracts.Provider, bucket string) {
	// False positive: borrowing the central worker singleton for logging is leak-free and safe;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	central.DefaultWorker().Logger().InfoContext(ctx, "storage client initialization complete",
		"provider", string(provider),
		"bucket", bucket,
	)
}

func maybeProbeStartup(ctx context.Context, c *Client, options *Options) error {
	if options.SkipProbe {
		return nil
	}
	return probeStartupClient(ctx, c)
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
	if err := maybeProbeStartup(ctx, c, options); err != nil {
		return nil, err
	}
	logStartupComplete(ctx, provider, cfg.Bucket)
	return c, nil
}
