package detector

import "strings"

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
