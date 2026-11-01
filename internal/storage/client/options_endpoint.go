// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"net"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

func formatTargetURL(u *url.URL, port int) string {
	if port > 0 && u.Port() == "" {
		u.Host = net.JoinHostPort(u.Hostname(), strconv.Itoa(port))
	}
	return u.String()
}

func applyInsecureScheme(u *url.URL, cfg *config.Config) {
	if isInsecure(cfg) {
		u.Scheme = "http"
	}
}

func isInsecure(cfg *config.Config) bool {
	return cfg != nil && (cfg.Insecure || cfg.DisableSSL)
}

func normalizeEndpointTargetHost(u *url.URL, bucket string) error {
	if isMultiTenantHost(u.Hostname()) {
		return normalizeMultiTenantEndpoint(u, bucket)
	}
	return nil
}

func parseEndpointURL(target, bucket string, cfg *config.Config) (*url.URL, error) {
	u, err := detector.ParseTarget(target)
	if err != nil {
		return nil, err
	}
	u.User = nil
	applyInsecureScheme(u, cfg)
	if err := normalizeEndpointTargetHost(u, bucket); err != nil {
		return nil, err
	}
	return u, nil
}

func parseBucketURL(target, bucket string, cfg *config.Config) (*url.URL, error) {
	u, err := detector.ParseTarget(target)
	if err != nil {
		return nil, err
	}
	u.User = nil
	applyInsecureScheme(u, cfg)
	if err := normalizeBucketURL(u, bucket); err != nil {
		return nil, err
	}
	return u, nil
}

func defaultProviderEndpoint(provider contracts.Provider) *string {
	if provider == contracts.ProviderGCS {
		return aws.String("https://storage.googleapis.com")
	}
	return nil
}

func finalizeEndpointURL(cfg *config.Config, provider contracts.Provider) (*string, error) {
	u, err := parseEndpointURL(cfg.Endpoint, cfg.Bucket, cfg)
	if err != nil {
		return nil, err
	}
	formatted := aws.String(formatTargetURL(u, cfg.Port))
	return finalizeEndpoint(formatted, provider), nil
}

func finalizeBucketURL(cfg *config.Config, provider contracts.Provider) (*string, error) {
	u, err := parseBucketURL(cfg.BucketURL, cfg.Bucket, cfg)
	if err != nil {
		return nil, err
	}
	formatted := aws.String(formatTargetURL(u, cfg.Port))
	return finalizeEndpoint(formatted, provider), nil
}

func resolveEndpointWithError(cfg *config.Config, provider contracts.Provider) (*string, error) {
	if cfg == nil {
		return defaultProviderEndpoint(provider), nil
	}
	if cfg.Endpoint != "" {
		return finalizeEndpointURL(cfg, provider)
	}
	if cfg.BucketURL != "" {
		return finalizeBucketURL(cfg, provider)
	}
	return defaultProviderEndpoint(provider), nil
}
