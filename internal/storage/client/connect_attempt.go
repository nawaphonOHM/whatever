// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func handleConnectionError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	msg := fmt.Sprintf("failed to connect to storage: %v; exiting abnormally", err)
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
