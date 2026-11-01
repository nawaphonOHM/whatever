// Package urlutil provides shared URL scheme inspection, validation, and security sanitization
// for cloud object storage endpoints and configurations.
package urlutil

import (
	"net/url"
	"strings"
)

const (
	// UnparseableURLPlaceholder is returned when a URL cannot be parsed safely or contains ambiguous malformed structures.
	UnparseableURLPlaceholder = "<unparseable>"

	schemeSep           = "://"
	defaultSchemePrefix = "https://"
)

// IsValidScheme checks if the provided scheme complies with RFC 3986 section 3.1:
// scheme = ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
func IsValidScheme(scheme string) bool {
	if scheme == "" {
		return false
	}
	for i := 0; i < len(scheme); i++ {
		c := scheme[i]
		if i == 0 {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
				return false
			}
		} else {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '+' || c == '-' || c == '.') {
				return false
			}
		}
	}
	return true
}

// HasSchemePrefix returns true if raw begins with a valid RFC 3986 scheme followed by "://".
func HasSchemePrefix(raw string) bool {
	idx := strings.Index(raw, schemeSep)
	if idx <= 0 {
		return false
	}
	return IsValidScheme(raw[:idx])
}

// HasSchemeSingleSlash checks if raw begins with a valid RFC 3986 scheme followed by ":/" but NOT "://".
// Examples: "https:/host", "http:/admin:pass@minio:9000", "s3:/bucket/x", "x:/host".
func HasSchemeSingleSlash(raw string) bool {
	idx := strings.Index(raw, ":/")
	if idx <= 0 {
		return false
	}
	return IsValidScheme(raw[:idx]) && !strings.HasPrefix(raw[idx:], schemeSep)
}

func sanitizeParsedURL(u *url.URL) {
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	u.RawFragment = ""
}

// isAmbiguousMalformed checks whether parsed components represent an ambiguous malformed input
// (such as empty-authority, protocol-relative, single-slash schemes, or malformed pseudo-schemes in path)
// that might leak userinfo.
func isAmbiguousMalformed(raw string, u *url.URL) bool {
	// Case 1: Single-slash schemes (e.g. "https:/host", "http:/admin:pass@minio:9000", "s3:/user:pass@host/x").
	if HasSchemeSingleSlash(raw) {
		return true
	}

	// Case 2: Leading slash or backslash in raw input (e.g. "//admin:pass@host", "/user:pass@host", "\\host").
	if strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "\\") || strings.Contains(raw, "\\") {
		return true
	}

	// Case 3: Host starts or ends with colon (e.g. "https:", "host:").
	if strings.HasSuffix(u.Host, ":") || strings.HasPrefix(u.Host, ":") {
		return true
	}

	// Case 4: Host is bare scheme name resulting from schemeless single-slash or colon parsing.
	if u.Host == "http:" || u.Host == "https:" || u.Hostname() == "http" || u.Hostname() == "https" {
		return true
	}

	hasAt := strings.Contains(raw, "@")
	if !hasAt {
		return false
	}

	// Case 5: Empty host when userinfo was present.
	if u.Host == "" {
		return true
	}

	// Case 6: Path, RawPath, or Opaque contains "@" (userinfo leaked into path).
	if strings.Contains(u.Path, "@") || strings.Contains(u.RawPath, "@") || strings.Contains(u.Opaque, "@") {
		return true
	}

	// Case 7: Host contains "@" (multiple @ or malformed authority).
	if strings.Contains(u.Host, "@") {
		return true
	}

	// Case 8: Path contains scheme delimiter "://" along with "@" (e.g. "a/b://user:pass@host").
	if strings.Contains(u.Path, schemeSep) {
		return true
	}

	return false
}

func sanitizeSchemelessURL(raw string) string {
	if strings.HasPrefix(raw, "/") && strings.Contains(raw, "@") {
		return UnparseableURLPlaceholder
	}

	u, err := url.Parse(defaultSchemePrefix + raw)
	if err != nil {
		return UnparseableURLPlaceholder
	}

	if isAmbiguousMalformed(raw, u) {
		return UnparseableURLPlaceholder
	}

	sanitizeParsedURL(u)
	res := u.String()
	if strings.HasPrefix(res, defaultSchemePrefix) {
		return strings.TrimPrefix(res, defaultSchemePrefix)
	}
	return strings.TrimPrefix(res, "https:")
}

func sanitizeAbsoluteURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return UnparseableURLPlaceholder
	}

	if isAmbiguousMalformed(raw, u) {
		return UnparseableURLPlaceholder
	}

	sanitizeParsedURL(u)
	return u.String()
}

func sanitizePopulatedURL(raw string) string {
	if HasSchemePrefix(raw) {
		return sanitizeAbsoluteURL(raw)
	}
	return sanitizeSchemelessURL(raw)
}

// SanitizeURL redacts userinfo credentials, query strings, and fragments from URLs.
// If raw is empty, it returns "". If parsing fails, invalid whitespace is present,
// or ambiguous malformed input is detected, it returns "<unparseable>".
func SanitizeURL(raw string) string {
	if raw == "" {
		return ""
	}
	if strings.ContainsAny(raw, " \t\r\n") {
		return UnparseableURLPlaceholder
	}
	return sanitizePopulatedURL(raw)
}
