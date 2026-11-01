// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

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

func validateConfig(cfg *config.Config, provider contracts.Provider) error {
	if err := validateConfigSettings(cfg); err != nil {
		return err
	}
	_, _, err := resolveConfigRuntime(cfg, provider)
	return err
}

func formatEffectiveChecksumMode(calc aws.RequestChecksumCalculation) string {
	if calc == aws.RequestChecksumCalculationWhenRequired {
		return "when_required"
	}
	return "when_supported"
}

type lifecycleLogParams struct {
	provider contracts.Provider
	bucket   string
	s3Opts   s3.Options
}

func logLifecycleInfo(
	ctx context.Context,
	msg string,
	p lifecycleLogParams,
) {
	// False positive: borrowing the central worker singleton for logging is leak-free and safe;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	central.DefaultWorker().Logger().InfoContext(ctx, msg,
		"provider", string(p.provider),
		"bucket", p.bucket,
		"region", p.s3Opts.Region,
		"endpoint", formatLoggedEndpoint(p.s3Opts.BaseEndpoint),
		"use_path_style", p.s3Opts.UsePathStyle,
		"checksum_mode", formatEffectiveChecksumMode(p.s3Opts.RequestChecksumCalculation),
	)
}

func logInitialization(ctx context.Context, p lifecycleLogParams) {
	logLifecycleInfo(ctx, "initializing storage client", p)
}

func logStartupComplete(ctx context.Context, p lifecycleLogParams) {
	logLifecycleInfo(ctx, "storage client initialization complete", p)
}
