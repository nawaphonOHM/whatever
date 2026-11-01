// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"github.com/nawaphonOHM/whatever/v3/internal/storage/urlutil"
)

// sanitizeLoggedURL redacts userinfo credentials, query strings, and fragments
// from logged URLs. If raw is empty, it returns "". If parsing fails, invalid
// whitespace is present, or ambiguous malformed input is detected, it returns "<unparseable>".
func sanitizeLoggedURL(raw string) string {
	return urlutil.SanitizeURL(raw)
}

func formatLoggedEndpoint(ep *string) string {
	if ep == nil {
		return ""
	}
	return urlutil.SanitizeURL(*ep)
}
