package config

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	defaultFirestorePort = 443
)

// resolveFirestorePort returns 443 for default or unspecified ports on Firestore.
func (c *Config) resolveFirestorePort() int {
	if c.Port == 0 || c.Port == defaultPort {
		return defaultFirestorePort
	}
	return c.Port
}

// buildFirestoreQuery constructs query string for Firestore connection URI.
func (c *Config) buildFirestoreQuery() string {
	query := fmt.Sprintf(
		"uuidRepresentation=%s&tls=true&loadBalanced=true&retryWrites=false&authMechanism=SCRAM-SHA-256",
		url.QueryEscape(c.resolveUUIDRep()),
	)
	query = c.appendAuthSourceQuery(query)
	return c.appendAppNameQuery(query)
}

// IsFirestoreURL reports whether the provided host or URI targets Google Cloud Firestore.
func IsFirestoreURL(target string) bool {
	return strings.Contains(strings.ToLower(target), "firestore.goog")
}

// IsFirestore reports whether the configuration targets Google Cloud Firestore.
func (c *Config) IsFirestore() bool {
	if c == nil {
		return false
	}
	return IsFirestoreURL(c.Host)
}
