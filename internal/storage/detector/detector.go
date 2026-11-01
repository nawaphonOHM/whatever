// Package detector provides cloud object storage provider detection from endpoints and bucket URLs.
package detector

import (
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

// Storage provider detection rule constants.
const (
	RuleDefaultAWS           = "default_aws"
	RuleAWSHostname          = "aws_hostname"
	RuleGCSHostname          = "gcs_hostname"
	RuleS3CompatibleFallback = "s3_compatible_fallback"
)

// Detection holds evaluated cloud storage provider and classification reasoning.
type Detection struct {
	Provider contracts.Provider
	Host     string
	Rule     string
	Bucket   string
}

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

// classifyDetection determines the provider, matched rule, and bucket for a normalized host.
func classifyDetection(host string) Detection {
	if IsGCSHost(host) {
		bucket, _ := gcsBucketSubdomain(host)
		return Detection{
			Provider: contracts.ProviderGCS,
			Host:     host,
			Rule:     RuleGCSHostname,
			Bucket:   bucket,
		}
	}
	if IsAWSHost(host) {
		bucket, _ := awsBucketSubdomain(host)
		return Detection{
			Provider: contracts.ProviderAWS,
			Host:     host,
			Rule:     RuleAWSHostname,
			Bucket:   bucket,
		}
	}
	return Detection{
		Provider: contracts.ProviderS3Compatible,
		Host:     host,
		Rule:     RuleS3CompatibleFallback,
	}
}

// Detect evaluates endpoint and bucket URL, returning full Detection information.
func Detect(endpoint, bucketURL string) (Detection, error) {
	target := selectTargetURL(endpoint, bucketURL)
	if target == "" {
		return Detection{
			Provider: contracts.ProviderAWS,
			Rule:     RuleDefaultAWS,
		}, nil
	}
	u, err := ParseTarget(target)
	if err != nil {
		return Detection{}, err
	}
	return classifyDetection(normalizeHost(u.Hostname())), nil
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
