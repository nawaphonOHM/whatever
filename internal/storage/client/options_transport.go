// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
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

// newTransportClient builds the SDK HTTP client whose TCP dial and TLS handshake are bounded by
// connectTimeout. No overall http.Client.Timeout is set because it would also cut off streamed
// response bodies; request deadlines are applied through contexts instead (see requestContext).
func newTransportClient(connectTimeout time.Duration) aws.HTTPClient {
	if connectTimeout <= 0 {
		return nil
	}
	return awshttp.NewBuildableClient().
		WithDialerOptions(func(d *net.Dialer) { d.Timeout = connectTimeout }).
		WithTransportOptions(func(tr *http.Transport) { tr.TLSHandshakeTimeout = connectTimeout })
}

func applyAWSChecksumDefaults(o *s3.Options) {
	if o.RequestChecksumCalculation == aws.RequestChecksumCalculationUnset {
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenSupported
	}
	if o.ResponseChecksumValidation == aws.ResponseChecksumValidationUnset {
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenSupported
	}
}

// applyProviderCompat configures the SDK's flexible checksums: WhenSupported on AWS
// and WhenRequired on non-AWS providers (which commonly reject aws-chunked trailers, notably GCS).
func applyProviderCompat(o *s3.Options, provider contracts.Provider) {
	if provider == contracts.ProviderAWS {
		applyAWSChecksumDefaults(o)
		return
	}
	o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
	o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
}

// applyProviderSettings applies addressing style and compatibility settings for the provider.
func applyProviderSettings(o *s3.Options, cfg *config.Config, provider contracts.Provider) {
	applyConfiguredPathStyle(o, cfg, provider)
	applyProviderCompat(o, provider)
	o.DisableLogOutputChecksumValidationSkipped = true
}
