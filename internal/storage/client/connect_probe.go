// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
)

// SetExitFunc overrides the process exit hook used by the central logger worker
// and returns the previous hook.
func SetExitFunc(fn func(int)) func(int) {
	// False positive: borrowing the central worker singleton to configure exit hook is leak-free;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	prev := central.DefaultWorker().ExitFunc()
	central.DefaultWorker().SetExitFunc(fn)
	return prev
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

func maybeProbeStartup(ctx context.Context, c *Client, options *Options) error {
	if options.SkipProbe {
		return nil
	}
	return probeStartupClient(ctx, c)
}
