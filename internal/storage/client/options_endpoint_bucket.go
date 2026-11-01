// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

const pathSep = "/"

func joinHostPort(host, port string) string {
	if port == "" {
		return host
	}
	return net.JoinHostPort(host, port)
}

func isMultiTenantHost(host string) bool {
	return detector.IsGCSHost(host) || detector.IsAWSHost(host)
}

func checkHostSubdomainMismatch(parts []string, bucket string) error {
	subdomain, hasBucket := detector.BucketSubdomain(strings.Join(parts, "."))
	if hasBucket && !strings.EqualFold(subdomain, bucket) {
		return fmt.Errorf("bucket %q in configuration does not match bucket %q in URL", bucket, subdomain)
	}
	return nil
}

func hasBucketHostPrefix(hostname, prefix string) bool {
	if len(hostname) <= len(prefix) {
		return false
	}
	return strings.EqualFold(hostname[:len(prefix)], prefix)
}

func stripBucketSubdomain(u *url.URL, bucket, hostname string) bool {
	if bucket == "" {
		return false
	}
	prefix := bucket + "."
	if !hasBucketHostPrefix(hostname, prefix) {
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
