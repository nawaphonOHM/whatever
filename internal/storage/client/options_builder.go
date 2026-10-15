// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func resolveCredentials(cfg *config.Config) aws.CredentialsProvider {
	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		return credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken)
	}
	return nil
}

func resolveOptionsHTTPClient(opts *Options) *http.Client {
	if opts != nil {
		return opts.HTTPClient
	}
	return nil
}

func resolveHTTPClient(cfg *config.Config, opts *Options) aws.HTTPClient {
	if client := resolveOptionsHTTPClient(opts); client != nil {
		return client
	}
	return newTransportClient(cfg.ConnectTimeout)
}

func resolveRegion(region string) string {
	if region != "" {
		return region
	}
	return defaultRegion
}

func applyConfiguredCreds(opts *s3.Options, cfg *config.Config) {
	if creds := resolveCredentials(cfg); creds != nil {
		opts.Credentials = creds
	}
}

func applyConfiguredHTTPClient(opts *s3.Options, cfg *config.Config, clientOpts *Options) {
	if client := resolveHTTPClient(cfg, clientOpts); client != nil {
		opts.HTTPClient = client
	}
}

func applyConfiguredBaseEndpoint(o *s3.Options, cfg *config.Config, provider contracts.Provider) error {
	ep, err := resolveEndpointWithError(cfg, provider)
	if err != nil {
		return err
	}
	if ep != nil {
		o.BaseEndpoint = ep
	}
	return nil
}

func applyRuntimeConfigSettings(o *s3.Options, cfg *config.Config, opts *Options) {
	applyConfiguredCreds(o, cfg)
	applyConfiguredHTTPClient(o, cfg, opts)
}

func applyS3RuntimeOptions(o *s3.Options, cfg *config.Config, provider contracts.Provider, opts *Options) error {
	reg, err := resolveEffectiveRegion(cfg)
	if err != nil {
		return err
	}
	o.Region = reg
	applyRuntimeConfigSettings(o, cfg, opts)
	if epErr := applyConfiguredBaseEndpoint(o, cfg, provider); epErr != nil {
		return epErr
	}
	applyProviderSettings(o, cfg, provider)
	opts.applyDriver(o)
	return nil
}

// BuildS3Options constructs AWS SDK v2 S3 options from configuration and functional options.
func BuildS3Options(cfg *config.Config, provider contracts.Provider, opts *Options) (s3.Options, error) {
	var s3Opts s3.Options
	if err := applyS3RuntimeOptions(&s3Opts, cfg, provider, opts); err != nil {
		return s3.Options{}, err
	}
	return s3Opts, nil
}
