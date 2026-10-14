// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"strings"
)

const (
	suffixVPCECN      = ".vpce.amazonaws.com.cn"
	suffixVPCE        = ".vpce.amazonaws.com"
	suffixAWSCN       = ".amazonaws.com.cn"
	suffixAWSStandard = ".amazonaws.com"
)

func trimKnownHostSuffix(host string) (string, bool) {
	lower := strings.ToLower(host)
	suffixes := []string{suffixVPCECN, suffixVPCE, suffixAWSCN, suffixAWSStandard}
	for _, suffix := range suffixes {
		if strings.HasSuffix(lower, suffix) {
			return lower[:len(lower)-len(suffix)], true
		}
	}
	return lower, false
}

func trimS3RegionPrefix(raw string) string {
	prefixes := []string{
		"s3-dualstack-",
		"s3-fips-",
		"s3-website-",
		"s3-website.",
		"s3-accesspoint-",
		"s3-accesspoint.",
		"s3-",
		"s3.",
		"website-",
		"accesspoint-",
	}
	res := raw
	for _, p := range prefixes {
		res = strings.TrimPrefix(res, p)
	}
	return res
}

func isIgnoredRegionKeyword(seg string) bool {
	switch seg {
	case "s3", "dualstack", "fips", "vpce":
		return true
	default:
		return false
	}
}

func cleanRegionSegment(seg string) string {
	cleaned := trimS3RegionPrefix(seg)
	if cleaned == "external-1" {
		return defaultRegion
	}
	if isIgnoredRegionKeyword(cleaned) {
		return ""
	}
	return cleaned
}

func extractAWSRegion(host string) string {
	prefix, ok := trimKnownHostSuffix(host)
	if !ok {
		return ""
	}
	parts := strings.Split(prefix, ".")
	if len(parts) == 0 {
		return ""
	}
	last := parts[len(parts)-1]
	return cleanRegionSegment(last)
}
