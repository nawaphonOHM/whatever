// Package storage provides the public entrypoint and abstract contract
// for cloud object storage operations across AWS S3 and S3-compatible providers.
package storage

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

// Client wraps the official AWS S3 client and provides unified cloud object storage operations.
type Client = client.Client

// Compile-time assertion that Client implements Storage.
var _ Storage = (*Client)(nil)

// Option defines a functional option for customizing storage client options.
type Option = client.Option

// Config encapsulates object storage configuration parameters.
type Config = config.Config

// DefaultConfig returns storage configuration populated with standard defaults.
func DefaultConfig() *Config {
	return config.DefaultConfig()
}

// SetExitFunc overrides the process exit hook used during connection failures.
func SetExitFunc(fn func(int)) func(int) {
	return client.SetExitFunc(fn)
}

// Connect loads storage configuration from environment variables (OHM9996_STORAGE_*),
// establishes an S3 client, optionally probes bucket connectivity, and returns a managed Client.
func Connect(ctx context.Context, opts ...Option) (*Client, error) {
	return client.Connect(ctx, opts...)
}

// ConnectWithConfig initializes a managed storage Client using the provided Config and Options.
func ConnectWithConfig(ctx context.Context, cfg *Config, opts ...Option) (*Client, error) {
	return client.ConnectWithConfig(ctx, cfg, opts...)
}

// WithDriverOptions appends raw AWS SDK s3.Options modifiers.
func WithDriverOptions(opts ...func(*s3.Options)) Option {
	return client.WithDriverOptions(opts...)
}

// WithPresignOptions appends raw AWS SDK PresignOptions modifiers.
func WithPresignOptions(opts ...func(*s3.PresignOptions)) Option {
	return client.WithPresignOptions(opts...)
}

// WithHTTPClient overrides the HTTP client used by the AWS SDK client.
func WithHTTPClient(httpClient *http.Client) Option {
	return client.WithHTTPClient(httpClient)
}

// WithSkipProbe controls whether startup connectivity probing is skipped.
func WithSkipProbe(skip bool) Option {
	return client.WithSkipProbe(skip)
}
