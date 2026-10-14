// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

func isNativeAWSDomain(host, domain string) bool {
	return host == domain || strings.HasSuffix(host, "."+domain)
}

func isNativeAWSHost(host string) bool {
	h := strings.ToLower(host)
	return isNativeAWSDomain(h, "amazonaws.com") || isNativeAWSDomain(h, "amazonaws.com.cn")
}

func extractTargetFromConfig(cfg *config.Config) string {
	if cfg.Endpoint != "" {
		return cfg.Endpoint
	}
	return cfg.BucketURL
}

func extractTargetHost(cfg *config.Config) string {
	target := extractTargetFromConfig(cfg)
	if target == "" {
		return ""
	}
	u, err := detector.ParseTarget(target)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func isCompatibleProvider(provider contracts.Provider) bool {
	return provider == contracts.ProviderGCS || provider == contracts.ProviderS3Compatible
}

func isPathStyleConfigured(cfg *config.Config, provider contracts.Provider) bool {
	return cfg.UsePathStyle || isCompatibleProvider(provider)
}

func shouldUsePathStyle(cfg *config.Config, provider contracts.Provider) bool {
	if isPathStyleConfigured(cfg, provider) {
		return true
	}
	host := extractTargetHost(cfg)
	return host != "" && !isNativeAWSHost(host)
}

func applyConfiguredPathStyle(opts *s3.Options, cfg *config.Config, provider contracts.Provider) {
	opts.UsePathStyle = shouldUsePathStyle(cfg, provider)
}
