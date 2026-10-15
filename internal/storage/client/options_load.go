// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

const errLoadAWSConfigFmt = "failed to load AWS default config: %w"

func appendHTTPClientLoadOption(
	loadOpts []func(*awsconfig.LoadOptions) error,
	cfg *config.Config,
	opts *Options,
) []func(*awsconfig.LoadOptions) error {
	if client := resolveHTTPClient(cfg, opts); client != nil {
		return append(loadOpts, awsconfig.WithHTTPClient(client))
	}
	return loadOpts
}

func awsLoadOptions(cfg *config.Config, opts *Options) ([]func(*awsconfig.LoadOptions) error, error) {
	reg, err := resolveEffectiveRegion(cfg)
	if err != nil {
		return nil, err
	}
	loadOpts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(reg),
	}
	return appendHTTPClientLoadOption(loadOpts, cfg, opts), nil
}

func loadAWSConfig(ctx context.Context, cfg *config.Config, opts *Options) (aws.Config, error) {
	loadOpts, err := awsLoadOptions(cfg, opts)
	if err != nil {
		return aws.Config{}, err
	}
	return awsconfig.LoadDefaultConfig(ctx, loadOpts...)
}
