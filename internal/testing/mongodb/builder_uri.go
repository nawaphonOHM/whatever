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

// resolvePort returns effective connection port for URI.
func resolvePort(o *Options) int {
	if o == nil || o.Port <= 0 {
		return DefaultPort
	}
	return o.Port
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
	if o == nil {
		return false
	}
	return o.Protocol == "mongodb+srv" || o.Port == 0
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

// appendAuthSourceQuery adds optional authSource parameter.
func appendAuthSourceQuery(query string, o *Options) string {
	if o != nil && o.AuthSource != "" {
		return fmt.Sprintf("%s&authSource=%s", query, url.QueryEscape(o.AuthSource))
	}
	return query
}

// appendAppNameQuery adds optional appName parameter.
func appendAppNameQuery(query string, o *Options) string {
	if o != nil && o.AppName != "" {
		return fmt.Sprintf("%s&appName=%s", query, url.QueryEscape(o.AppName))
	}
	return query
}

// appendDirectConnQuery adds directConnection parameter if enabled.
func appendDirectConnQuery(query string, o *Options) string {
	if o != nil && o.DirectConnection {
		return fmt.Sprintf("%s&directConnection=true", query)
	}
	return query
}

// buildQuery constructs query string for connection URI.
func buildQuery(o *Options, enableTLS bool) string {
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
		Scheme:   proto,
		User:     buildUserInfo(o),
		Host:     buildHost(o),
		Path:     buildPath(o),
		RawQuery: buildQuery(o, enableTLS),
	}
	return u.String()
}
