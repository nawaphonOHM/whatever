// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"net/url"
	"strings"
)

const pathSep = "/"

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
