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

func formatLoggedEndpoint(ep *string) string {
	if ep == nil {
		return ""
	}
	return *ep
}

func validateConfigSettings(cfg *config.Config) error {
	if cfg == nil {
		return contracts.ErrNilConfig
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf(errInvalidConfigFormat, err)
	}
	return nil
}

func resolveConfigRuntime(cfg *config.Config, provider contracts.Provider) (string, *string, error) {
	region, regErr := resolveEffectiveRegion(cfg)
	if regErr != nil {
		return "", nil, fmt.Errorf(errInvalidConfigFormat, regErr)
	}
	ep, epErr := resolveEndpointWithError(cfg, provider)
	if epErr != nil {
		return "", nil, fmt.Errorf(errInvalidConfigFormat, epErr)
	}
	return region, ep, nil
}

type initLogParams struct {
	cfg      *config.Config
	ep       *string
	region   string
	provider contracts.Provider
}

func logInitialization(ctx context.Context, p initLogParams) {
	// False positive: borrowing the central worker singleton for logging is leak-free and safe;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	central.DefaultWorker().Logger().InfoContext(ctx, "initializing storage client",
		"provider", string(p.provider),
		"bucket", p.cfg.Bucket,
		"region", p.region,
		"endpoint", formatLoggedEndpoint(p.ep),
		"use_path_style", shouldUsePathStyle(p.cfg, p.provider),
	)
}

func validateAndLogConfig(ctx context.Context, cfg *config.Config, provider contracts.Provider) error {
	if err := validateConfigSettings(cfg); err != nil {
		return err
	}
	region, ep, err := resolveConfigRuntime(cfg, provider)
	if err != nil {
		return err
	}
	logInitialization(ctx, initLogParams{
		cfg:      cfg,
		provider: provider,
		region:   region,
		ep:       ep,
	})
	return nil
}
