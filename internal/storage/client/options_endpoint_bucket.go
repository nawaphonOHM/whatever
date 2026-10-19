// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

const pathSep = "/"

func joinHostPort(host, port string) string {
	if port == "" {
		return host
	}
	return net.JoinHostPort(host, port)
}

func isExactNonBucketLabel(label string) bool {
	switch label {
	case "s3", "storage", "stroage", "bucket":
		return true
	default:
		return false
	}
}

func isNonBucketLabel(label string) bool {
	lower := strings.ToLower(label)
	if isExactNonBucketLabel(lower) {
		return true
	}
	return strings.HasPrefix(lower, "s3-") || strings.Contains(lower, "vpce")
}

func isMultiTenantHost(host string) bool {
	lower := strings.ToLower(host)
	if strings.Contains(lower, "storage.googleapis") || strings.Contains(lower, "stroage.googleapis") {
		return true
	}
	return isNativeAWSHost(lower)
}

func isAllowedHostSubdomain(subdomain, bucket string) bool {
	return subdomain == bucket || isNonBucketLabel(subdomain)
}

func checkHostSubdomainMismatch(parts []string, bucket string) error {
	if len(parts) <= 1 || isAllowedHostSubdomain(parts[0], bucket) {
		return nil
	}
	return fmt.Errorf("bucket %q in configuration does not match bucket %q in URL", bucket, parts[0])
}

func stripBucketSubdomain(u *url.URL, bucket, hostname string) bool {
	if bucket == "" {
		return false
	}
	prefix := bucket + "."
	if !strings.HasPrefix(hostname, prefix) {
		return false
	}
	remainder := hostname[len(prefix):]
	u.Host = joinHostPort(remainder, u.Port())
	return true
}

func findBucketIndexInSegments(segments []string, bucket string) int {
	for i, seg := range segments {
		if seg == bucket {
			return i
		}
	}
	return -1
}

func applyPrefixSegments(u *url.URL, segments []string, idx int) {
	prefixSegments := segments[:idx]
	if len(prefixSegments) == 0 {
		u.Path = ""
	} else {
		u.Path = pathSep + strings.Join(prefixSegments, pathSep)
	}
	u.RawPath = ""
}

func checkStrippedOrError(stripped bool, err error) (bool, error) {
	if err != nil {
		return true, err
	}
	return stripped, nil
}

func stripBucketFromPath(u *url.URL, bucket string) (bool, error) {
	trimmed := strings.Trim(u.Path, pathSep)
	if trimmed == "" {
		u.Path = ""
		return false, nil
	}
	segments := strings.Split(trimmed, pathSep)
	idx := findBucketIndexInSegments(segments, bucket)
	if idx >= 0 {
		applyPrefixSegments(u, segments, idx)
		return true, nil
	}
	return false, nil
}
