// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

const (
	defaultRegion     = "us-east-1"
	errRegionMismatch = "region mismatch: config region %q conflicts with endpoint region %q"
)

func detectURLRegion(target string) string {
	if target == "" {
		return ""
	}
	u, err := detector.ParseTarget(target)
	if err != nil {
		return ""
	}
	return extractAWSRegion(u.Hostname())
}

func selectDetectedURLRegion(cfg *config.Config) string {
	if reg := detectURLRegion(cfg.Endpoint); reg != "" {
		return reg
	}
	if reg := detectURLRegion(cfg.BucketURL); reg != "" {
		return reg
	}
	return ""
}

func isRegionCompatible(configured, detected string) bool {
	return configured == "" || configured == defaultRegion || configured == detected
}

func checkRegionConflict(configured, detected string) error {
	if isRegionCompatible(configured, detected) {
		return nil
	}
	return fmt.Errorf(errRegionMismatch, configured, detected)
}

func resolveRegionWithDetected(configured, detected string) (string, error) {
	if detected == "" {
		return resolveRegion(configured), nil
	}
	if err := checkRegionConflict(configured, detected); err != nil {
		return "", err
	}
	return detected, nil
}

func resolveEffectiveRegion(cfg *config.Config) (string, error) {
	if cfg == nil {
		return defaultRegion, nil
	}
	return resolveRegionWithDetected(cfg.Region, selectDetectedURLRegion(cfg))
}
