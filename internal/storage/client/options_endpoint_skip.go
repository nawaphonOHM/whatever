// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"net/url"
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	labelS3        = "s3"
	labelS3Prefix  = "s3-"
	hostS3Global   = "s3.amazonaws.com"
	maxPrefixParts = 3
)

func hasNonStandardAWSHostKeywords(lower string) bool {
	return strings.Contains(lower, ".vpce.") ||
		strings.Contains(lower, "fips") ||
		strings.Contains(lower, "dualstack")
}

func isS3Label(label string) bool {
	return label == labelS3 || strings.HasPrefix(label, labelS3Prefix)
}

func isStandardPart1(parts []string) bool {
	return len(parts) == 1 && isS3Label(parts[0])
}

func isStandardPart2(parts []string) bool {
	return len(parts) == 2 && (parts[0] == labelS3 || isS3Label(parts[1]))
}

func isStandardPart3(parts []string) bool {
	return len(parts) == maxPrefixParts && isS3Label(parts[1])
}

func isStandardAWSS3PrefixParts(parts []string) bool {
	if isStandardPart1(parts) || isStandardPart2(parts) {
		return true
	}
	return isStandardPart3(parts)
}

func matchesAWSStandardSuffix(lower string) bool {
	return lower == hostS3Global || strings.HasSuffix(lower, suffixAWSStandard)
}

func isStandardAWSS3Host(host string) bool {
	lower := strings.ToLower(strings.TrimSuffix(host, "."))
	if !matchesAWSStandardSuffix(lower) || hasNonStandardAWSHostKeywords(lower) {
		return false
	}
	prefix := strings.TrimSuffix(lower, suffixAWSStandard)
	parts := strings.Split(prefix, ".")
	return isStandardAWSS3PrefixParts(parts)
}

func isStandardAWSS3Endpoint(endpoint string) bool {
	u, err := url.Parse(endpoint)
	if err != nil {
		return false
	}
	return isStandardAWSS3Host(u.Hostname())
}

func shouldSkipAWSBaseEndpoint(provider contracts.Provider, endpoint string) bool {
	if provider != contracts.ProviderAWS {
		return false
	}
	return isStandardAWSS3Endpoint(endpoint)
}

func finalizeEndpoint(formatted *string, provider contracts.Provider) *string {
	if formatted == nil {
		return nil
	}
	if shouldSkipAWSBaseEndpoint(provider, *formatted) {
		return nil
	}
	return formatted
}
