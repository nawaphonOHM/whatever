// Package detector provides cloud object storage provider detection from endpoints and bucket URLs.
package detector

import (
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const labelSep = "."

// selectTargetURL returns the primary target URL from endpoint or bucketURL.
func selectTargetURL(endpoint, bucketURL string) string {
	if endpoint != "" {
		return endpoint
	}
	return bucketURL
}

// normalizeHost lowercases host and trims a trailing root-label dot.
func normalizeHost(host string) string {
	return strings.ToLower(strings.TrimSuffix(host, labelSep))
}

// matchesDomainOrSubdomain checks if host equals domain or ends with .domain.
func matchesDomainOrSubdomain(host, domain string) bool {
	return host == domain || strings.HasSuffix(host, labelSep+domain)
}

// prefixLabels returns the dot-separated labels preceding the given domain suffix.
func prefixLabels(host, domain string) []string {
	prefix := strings.TrimSuffix(strings.TrimSuffix(host, domain), labelSep)
	if prefix == "" {
		return nil
	}
	return strings.Split(prefix, labelSep)
}

// classifyProvider inspects parsed hostname and determines the storage provider.
func classifyProvider(host string) contracts.Provider {
	if IsGCSHost(host) {
		return contracts.ProviderGCS
	}
	if IsAWSHost(host) {
		return contracts.ProviderAWS
	}
	return contracts.ProviderS3Compatible
}

// DetectProvider identifies the cloud storage provider based on endpoint and bucket URL patterns.
// If both arguments are empty, it defaults to AWS S3.
func DetectProvider(endpoint, bucketURL string) (contracts.Provider, error) {
	target := selectTargetURL(endpoint, bucketURL)
	if target == "" {
		return contracts.ProviderAWS, nil
	}
	u, err := ParseTarget(target)
	if err != nil {
		return "", err
	}
	return classifyProvider(normalizeHost(u.Hostname())), nil
}

// BucketSubdomain extracts the virtual-hosted bucket name from a multi-tenant provider hostname
// (e.g. AWS S3 or Google Cloud Storage) if present. It returns ("", false) if the hostname is a base service
// endpoint or not a recognized multi-tenant storage host.
func BucketSubdomain(host string) (string, bool) {
	h := normalizeHost(host)
	if IsGCSHost(h) {
		return gcsBucketSubdomain(h)
	}
	if IsAWSHost(h) {
		return awsBucketSubdomain(h)
	}
	return "", false
}
