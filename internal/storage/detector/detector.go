// Package detector provides cloud object storage provider detection from endpoints and bucket URLs.
package detector

import (
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const labelSep = "."

// Provider detection classification rules.
const (
	RuleDefaultAWS           = "RuleDefaultAWS"
	RuleAWSHostname          = "RuleAWSHostname"
	RuleGCSHostname          = "RuleGCSHostname"
	RuleS3CompatibleFallback = "RuleS3CompatibleFallback"
)

// Detection encapsulates provider classification and inspected metadata.
type Detection struct {
	Provider contracts.Provider
	Host     string
	Rule     string
	Bucket   string
}

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

// Detect evaluates endpoint and bucket URL to determine the active storage provider and associated metadata.
// If endpoint is non-empty, it is used as the classification target; otherwise bucketURL is used.
// If both are empty, it defaults to ProviderAWS under RuleDefaultAWS.
func Detect(endpoint, bucketURL string) (Detection, error) {
	target := selectTargetURL(endpoint, bucketURL)
	if target == "" {
		return Detection{
			Provider: contracts.ProviderAWS,
			Host:     "",
			Rule:     RuleDefaultAWS,
			Bucket:   "",
		}, nil
	}

	u, err := ParseTarget(target)
	if err != nil {
		return Detection{}, err
	}

	host := normalizeHost(u.Hostname())
	if IsGCSHost(host) {
		bucket, _ := gcsBucketSubdomain(host)
		return Detection{
			Provider: contracts.ProviderGCS,
			Host:     host,
			Rule:     RuleGCSHostname,
			Bucket:   bucket,
		}, nil
	}

	if IsAWSHost(host) {
		bucket, _ := awsBucketSubdomain(host)
		return Detection{
			Provider: contracts.ProviderAWS,
			Host:     host,
			Rule:     RuleAWSHostname,
			Bucket:   bucket,
		}, nil
	}

	return Detection{
		Provider: contracts.ProviderS3Compatible,
		Host:     host,
		Rule:     RuleS3CompatibleFallback,
		Bucket:   "",
	}, nil
}

// DetectProvider identifies the cloud storage provider based on endpoint and bucket URL patterns.
// If both arguments are empty, it defaults to AWS S3.
func DetectProvider(endpoint, bucketURL string) (contracts.Provider, error) {
	d, err := Detect(endpoint, bucketURL)
	if err != nil {
		return "", err
	}
	return d.Provider, nil
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
