package mongodb

import (
	"fmt"
	"net/url"
)

// hasCredentials reports whether authentication credentials are provided.
func hasCredentials(o *Options) bool {
	if o == nil {
		return false
	}
	if o.Username != "" {
		return true
	}
	return o.Password != ""
}

// buildUserInfo constructs user credentials for connection URI if provided.
func buildUserInfo(o *Options) *url.Userinfo {
	if !hasCredentials(o) {
		return nil
	}
	return url.UserPassword(o.Username, o.Password)
}

func resolveStandardPort(port int) int {
	if port <= 0 {
		return DefaultPort
	}
	return port
}

// resolvePort returns effective connection port for URI.
func resolvePort(o *Options) int {
	if o == nil {
		return DefaultPort
	}
	if o.IsFirestore() {
		return resolveFirestorePort(o)
	}
	return resolveStandardPort(o.Port)
}

// resolveHost returns the host address from options or the default host.
func resolveHost(o *Options) string {
	if o == nil || o.Host == "" {
		return DefaultHost
	}
	return o.Host
}

// isSrvOrZeroPort reports whether the connection uses mongodb+srv or zero port.
func isSrvOrZeroPort(o *Options) bool {
	if o == nil || o.IsFirestore() {
		return false
	}
	return isDefaultOrSrv(o)
}

func isDefaultOrSrv(o *Options) bool {
	return o.Port == 0 || o.Protocol == "mongodb+srv"
}

// buildHost constructs host and optional port for connection URI.
func buildHost(o *Options) string {
	host := resolveHost(o)
	if isSrvOrZeroPort(o) {
		return host
	}
	return fmt.Sprintf("%s:%d", host, resolvePort(o))
}

// resolveUUIDRep returns effective UUID representation for URI.
func resolveUUIDRep(o *Options) string {
	if o == nil || o.UUIDRepresentation == "" {
		return DefaultUUIDRepresentation
	}
	return o.UUIDRepresentation
}

// buildQuery constructs query string for connection URI.
func buildQuery(o *Options, enableTLS bool) string {
	if o != nil && o.IsFirestore() {
		return buildFirestoreQuery(o)
	}
	query := fmt.Sprintf("uuidRepresentation=%s&tls=%t", url.QueryEscape(resolveUUIDRep(o)), enableTLS)
	query = appendAuthSourceQuery(query, o)
	query = appendAppNameQuery(query, o)
	return appendDirectConnQuery(query, o)
}

// buildPath constructs the URI path from configured database.
func buildPath(o *Options) string {
	if o != nil && o.Database != "" {
		return fmt.Sprintf("/%s", url.PathEscape(o.Database))
	}
	return "/"
}

// BuildURI constructs a standard MongoDB connection URI string based on Options and enableTLS.
func BuildURI(o *Options, enableTLS bool) string {
	proto := DefaultProtocol
	if o != nil && o.Protocol != "" {
		proto = o.Protocol
	}
	u := &url.URL{
		Scheme: proto,
		// False positive: nil User is intentional when credentials are absent; proof:
		// TestBuildURI_URLStructNilUserProof in builder_nil_proof_test.go.
		User:     buildUserInfo(o),
		Host:     buildHost(o),
		Path:     buildPath(o),
		RawQuery: buildQuery(o, enableTLS),
	}
	return u.String()
}
