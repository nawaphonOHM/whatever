// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

type s3ClientParams struct {
	cfg      *config.Config
	opts     *Options
	provider contracts.Provider
}

func hasStaticCredentials(cfg *config.Config) bool {
	return cfg != nil && cfg.AccessKeyID != "" && cfg.SecretAccessKey != ""
}

func applyRuntimeBaseEndpoint(o *s3.Options, p s3ClientParams) {
	ep, err := resolveEndpointWithError(p.cfg, p.provider)
	if err == nil && ep != nil {
		o.BaseEndpoint = ep
	}
}

func applyRuntimeOptionsSafe(o *s3.Options, p s3ClientParams) {
	applyConfiguredCreds(o, p.cfg)
	applyRuntimeBaseEndpoint(o, p)
	applyConfiguredHTTPClient(o, p.cfg, p.opts)
	applyProviderSettings(o, p.cfg, p.provider)
	p.opts.applyDriver(o)
}

func bindS3Options(p s3ClientParams) func(*s3.Options) {
	return func(o *s3.Options) {
		applyRuntimeOptionsSafe(o, p)
	}
}

func newStaticS3Client(p s3ClientParams) (*s3.Client, error) {
	opts, err := BuildS3Options(p.cfg, p.provider, p.opts)
	if err != nil {
		return nil, err
	}
	return s3.New(opts), nil
}

func newDefaultS3Client(ctx context.Context, p s3ClientParams) (*s3.Client, error) {
	awsCfg, err := loadAWSConfig(ctx, p.cfg, p.opts)
	if err != nil {
		return nil, fmt.Errorf(errLoadAWSConfigFmt, err)
	}
	return s3.NewFromConfig(awsCfg, bindS3Options(p)), nil
}

func newS3Client(ctx context.Context, p s3ClientParams) (*s3.Client, error) {
	if hasStaticCredentials(p.cfg) {
		return newStaticS3Client(p)
	}
	return newDefaultS3Client(ctx, p)
}

func presignOptions(options *Options) []func(*s3.PresignOptions) {
	if options == nil {
		return nil
	}
	return options.PresignOptions
}

func buildClientInstance(
	ctx context.Context,
	cfg *config.Config,
	provider contracts.Provider,
	options *Options,
) (*Client, error) {
	raw, err := newS3Client(ctx, s3ClientParams{cfg: cfg, provider: provider, opts: options})
	if err != nil {
		return nil, err
	}
	presign := s3.NewPresignClient(raw, presignOptions(options)...)
	return NewClient(ClientParams{
		Raw:            raw,
		Presign:        presign,
		Provider:       provider,
		DefaultBucket:  cfg.Bucket,
		PresignExpiry:  cfg.PresignExpiry,
		RequestTimeout: cfg.RequestTimeout,
		PingTimeout:    cfg.PingTimeout,
	}), nil
}
