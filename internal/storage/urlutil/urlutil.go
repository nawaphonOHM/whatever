// Package urlutil provides shared URL scheme inspection, validation, and security sanitization
// for cloud object storage endpoints and configurations.
package urlutil

import (
	"strings"
)

const (
	// UnparseableURLPlaceholder is returned when a URL cannot be parsed safely or contains ambiguous structures.
	UnparseableURLPlaceholder = "<unparseable>"

	schemeSep           = "://"
	defaultSchemePrefix = "https://"
	atChar              = "@"
)

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

func sanitizePopulatedURL(raw string) string {
	if HasSchemePrefix(raw) {
		return sanitizeAbsoluteURL(raw)
	}
	return sanitizeSchemelessURL(raw)
}
