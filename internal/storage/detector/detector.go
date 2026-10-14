// Package detector provides cloud object storage provider detection from endpoints and bucket URLs.
package detector

import (
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

var gcsPatterns = []string{
	"storage.googleapis.com",
	"storage.googleapis",
	"stroage.googleapis.com",
	"stroage.googleapis",
}

// selectTargetURL returns the primary target URL from endpoint or bucketURL.
func selectTargetURL(endpoint, bucketURL string) string {
	if endpoint != "" {
		return endpoint
	}
	return bucketURL
}

// matchesGCSPattern checks if hostname equals pattern or ends with .pattern.
func matchesGCSPattern(host, pattern string) bool {
	return host == pattern || strings.HasSuffix(host, "."+pattern)
}

// isGCSHost determines if the hostname matches Google Cloud Storage patterns.
func isGCSHost(host string) bool {
	for _, pattern := range gcsPatterns {
		if matchesGCSPattern(host, pattern) {
			return true
		}
	}
	return false
}

// isAWSLabel checks if a dot-separated hostname label indicates AWS S3.
func isAWSLabel(label string) bool {
	return label == "s3" || strings.HasPrefix(label, "s3-")
}

// hasAWSLabel checks if any hostname label matches AWS S3 patterns.
func hasAWSLabel(host string) bool {
	for label := range strings.SplitSeq(host, ".") {
		if isAWSLabel(label) {
			return true
		}
	}
	return false
}

// isAWSHostname checks if hostname ends with AWS domain suffix.
func isAWSHostname(host string) bool {
	return host == "amazonaws.com" || strings.HasSuffix(host, ".amazonaws.com")
}

// isAWSHost determines if hostname matches AWS S3 endpoints.
func isAWSHost(host string) bool {
	if hasAWSLabel(host) {
		return true
	}
	return isAWSHostname(host)
}

// classifyProvider inspects parsed hostname and determines the storage provider.
func classifyProvider(host string) contracts.Provider {
	if isGCSHost(host) {
		return contracts.ProviderGCS
	}
	if isAWSHost(host) {
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
	host := strings.TrimSuffix(strings.ToLower(u.Hostname()), ".")
	return classifyProvider(host), nil
}
