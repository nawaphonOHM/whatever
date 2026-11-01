// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"net/url"
	"strings"
)

const (
	unparseableURLPlaceholder = "<unparseable>"
	schemeSep                 = "://"
	defaultSchemePrefix       = "https://"
)

func sanitizeParsedURL(u *url.URL) {
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	u.RawFragment = ""
}

func sanitizeSchemelessURL(raw string) string {
	u, err := url.Parse(defaultSchemePrefix + raw)
	if err != nil {
		return unparseableURLPlaceholder
	}
	sanitizeParsedURL(u)
	return strings.TrimPrefix(u.String(), defaultSchemePrefix)
}

func sanitizeAbsoluteURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return unparseableURLPlaceholder
	}
	sanitizeParsedURL(u)
	return u.String()
}

func sanitizePopulatedURL(raw string) string {
	if strings.Contains(raw, schemeSep) {
		return sanitizeAbsoluteURL(raw)
	}
	return sanitizeSchemelessURL(raw)
}

// sanitizeLoggedURL redacts userinfo credentials, query strings, and fragments
// from logged URLs. If raw is empty, it returns "". If parsing fails or invalid
// whitespace is present, it returns "<unparseable>".
func sanitizeLoggedURL(raw string) string {
	if raw == "" {
		return ""
	}
	if strings.ContainsAny(raw, " \t\r\n") {
		return unparseableURLPlaceholder
	}
	return sanitizePopulatedURL(raw)
}

func formatLoggedEndpoint(ep *string) string {
	if ep == nil {
		return ""
	}
	return sanitizeLoggedURL(*ep)
}
