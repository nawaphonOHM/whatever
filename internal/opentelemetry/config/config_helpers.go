package config

import (
	"strings"
)

// isNonEmptyMatch returns true if skip is non-empty and equals path.
func isNonEmptyMatch(skip, path string) bool {
	if skip == "" {
		return false
	}
	return skip == path
}

// matchPath searches skipPaths for a match against path.
func matchPath(skipPaths []string, path string) bool {
	for _, skip := range skipPaths {
		if isNonEmptyMatch(skip, path) {
			return true
		}
	}
	return false
}

// ShouldSkip returns true if the given request path matches any path in SkipPaths.
func (c *Config) ShouldSkip(path string) bool {
	if c == nil {
		return false
	}
	return matchPath(c.SkipPaths, path)
}

// IsGRPC returns true if the configured protocol is gRPC.
func (c *Config) IsGRPC() bool {
	if c == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(c.Protocol), ProtocolGRPC)
}

// isHTTPProtocol returns true if protocol string matches any supported HTTP variant.
func isHTTPProtocol(p string) bool {
	return p == ProtocolHTTP || p == ProtocolHTTPProtobuf || p == ProtocolHTTPJSON
}

// IsHTTP returns true if the configured protocol is any HTTP-based protocol.
func (c *Config) IsHTTP() bool {
	if c == nil {
		return false
	}
	return isHTTPProtocol(strings.ToLower(strings.TrimSpace(c.Protocol)))
}
