package mongodb

import (
	"net/url"
	"strings"
)

// sanitizeURIPath extracts raw path string from URI.
func sanitizeURIPath(uri string) string {
	if uri == "" {
		return ""
	}
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(u.Path, "/")
}

// extractDatabaseFromURI extracts the database name from a URI path if present.
func extractDatabaseFromURI(uri string) string {
	raw := sanitizeURIPath(uri)
	if raw == "" {
		return ""
	}
	unescaped, err := url.PathUnescape(raw)
	if err != nil {
		return raw
	}
	return unescaped
}
