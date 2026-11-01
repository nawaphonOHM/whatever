package detector

import (
	"slices"
	"strings"
)

const (
	gcsRegionalDomain    = "rep.googleapis.com"
	gcsPSCDomain         = "p.googleapis.com"
	gcsRegionalMinLabels = 3
	gcsPSCMinLabels      = 2
)

var gcsSuffixPatterns = []string{
	"storage.googleapis.com",
	"storage.googleapis",
	"stroage.googleapis.com",
	"stroage.googleapis",
	"commondatastorage.googleapis.com",
	"storage.cloud.google.com",
}

func isGCSServiceLabel(label string) bool {
	return label == "storage" || label == "stroage"
}

func hasGCSPSCPrefix(label string) bool {
	return strings.HasPrefix(label, "storage-") || strings.HasPrefix(label, "stroage-")
}

func isGCSPSCLabel(label string) bool {
	return isGCSServiceLabel(label) || hasGCSPSCPrefix(label)
}

func hasGCSSuffix(host string) bool {
	return slices.ContainsFunc(gcsSuffixPatterns, func(pattern string) bool {
		return matchesDomainOrSubdomain(host, pattern)
	})
}

// isGCSRegionalHost matches [bucket.]storage.<region>.rep.googleapis.com.
func isGCSRegionalHost(host string) bool {
	if !matchesDomainOrSubdomain(host, gcsRegionalDomain) {
		return false
	}
	labels := prefixLabels(host, gcsRegionalDomain)
	return len(labels) >= 2 && isGCSServiceLabel(labels[len(labels)-2])
}

// isGCSPSCHost matches [bucket.]storage[-name].p.googleapis.com.
func isGCSPSCHost(host string) bool {
	if !matchesDomainOrSubdomain(host, gcsPSCDomain) {
		return false
	}
	labels := prefixLabels(host, gcsPSCDomain)
	return len(labels) > 0 && isGCSPSCLabel(labels[len(labels)-1])
}

func firstLabel(labels []string) (string, bool) {
	if len(labels) == 0 {
		return "", false
	}
	return labels[0], true
}

func gcsSuffixOf(host string) (string, bool) {
	idx := slices.IndexFunc(gcsSuffixPatterns, func(suffix string) bool {
		return matchesDomainOrSubdomain(host, suffix)
	})
	if idx < 0 {
		return "", false
	}
	return gcsSuffixPatterns[idx], true
}

func gcsRegionalBucket(host string) (string, bool) {
	labels := prefixLabels(host, gcsRegionalDomain)
	if len(labels) < gcsRegionalMinLabels || !isGCSServiceLabel(labels[len(labels)-2]) {
		return "", false
	}
	return labels[0], true
}

func gcsPSCBucket(host string) (string, bool) {
	labels := prefixLabels(host, gcsPSCDomain)
	if len(labels) < gcsPSCMinLabels || !isGCSPSCLabel(labels[len(labels)-1]) {
		return "", false
	}
	return labels[0], true
}

func gcsDedicatedBucket(host string) (string, bool) {
	if matchesDomainOrSubdomain(host, gcsRegionalDomain) {
		return gcsRegionalBucket(host)
	}
	if matchesDomainOrSubdomain(host, gcsPSCDomain) {
		return gcsPSCBucket(host)
	}
	return "", false
}

func gcsBucketSubdomain(host string) (string, bool) {
	if suffix, ok := gcsSuffixOf(host); ok {
		return firstLabel(prefixLabels(host, suffix))
	}
	return gcsDedicatedBucket(host)
}

// IsGCSHost determines if the hostname matches Google Cloud Storage patterns.
func IsGCSHost(host string) bool {
	h := normalizeHost(host)
	return hasGCSSuffix(h) || isGCSRegionalHost(h) || isGCSPSCHost(h)
}
