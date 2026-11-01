// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"fmt"
	"net/url"
	"strings"
)

func stripEndpointBucketHost(u *url.URL, bucket string) bool {
	if stripBucketSubdomain(u, bucket, u.Hostname()) {
		u.Path = ""
		u.RawPath = ""
		return true
	}
	return false
}

func handleEndpointHostSubdomain(u *url.URL, bucket string) error {
	if stripEndpointBucketHost(u, bucket) {
		return nil
	}
	return checkHostSubdomainMismatch(strings.Split(u.Hostname(), "."), bucket)
}

func normalizeMultiTenantEndpoint(u *url.URL, bucket string) error {
	if bucket == "" {
		return nil
	}
	done, err := checkStrippedOrError(stripBucketFromPath(u, bucket))
	if done {
		return err
	}
	return handleEndpointHostSubdomain(u, bucket)
}

func checkMultiTenantPathMismatch(u *url.URL, bucket string) error {
	if strings.Trim(u.Path, pathSep) != "" {
		if strings.Contains(u.Path, "@") || u.User != nil {
			return fmt.Errorf("bucket %q in configuration does not match bucket in URL path", bucket)
		}
		return fmt.Errorf("bucket %q in configuration does not match bucket in URL path %q", bucket, u.Path)
	}
	return nil
}

func handleBucketURLHostMismatch(u *url.URL, bucket string) error {
	if !isMultiTenantHost(u.Hostname()) {
		return nil
	}
	if err := checkHostSubdomainMismatch(strings.Split(u.Hostname(), "."), bucket); err != nil {
		return err
	}
	return checkMultiTenantPathMismatch(u, bucket)
}

func handleBucketURLHost(u *url.URL, bucket string) error {
	if stripBucketSubdomain(u, bucket, u.Hostname()) {
		u.Path = ""
		u.RawPath = ""
		return nil
	}
	return handleBucketURLHostMismatch(u, bucket)
}

func handlePathStrippedHostCheck(u *url.URL, bucket string) error {
	if isMultiTenantHost(u.Hostname()) {
		return checkHostSubdomainMismatch(strings.Split(u.Hostname(), "."), bucket)
	}
	return nil
}

func processBucketURLStripping(u *url.URL, bucket string) error {
	stripped, err := stripBucketFromPath(u, bucket)
	if err != nil {
		return err
	}
	if stripped {
		return handlePathStrippedHostCheck(u, bucket)
	}
	return handleBucketURLHost(u, bucket)
}

func normalizeBucketURL(u *url.URL, bucket string) error {
	if bucket == "" {
		return nil
	}
	return processBucketURLStripping(u, bucket)
}
