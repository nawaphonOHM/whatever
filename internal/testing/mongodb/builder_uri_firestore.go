package mongodb

import (
	"fmt"
	"net/url"
)

// resolveFirestorePort returns 443 for default or unspecified ports on Firestore.
func resolveFirestorePort(o *Options) int {
	if o.Port == 0 || o.Port == DefaultPort {
		return DefaultFirestorePort
	}
	return o.Port
}

// buildFirestoreQuery constructs query string for Firestore connection URI.
func buildFirestoreQuery(o *Options) string {
	query := fmt.Sprintf(
		"uuidRepresentation=%s&tls=true&loadBalanced=true&retryWrites=false&authMechanism=SCRAM-SHA-256",
		url.QueryEscape(resolveUUIDRep(o)),
	)
	query = appendAuthSourceQuery(query, o)
	return appendAppNameQuery(query, o)
}
