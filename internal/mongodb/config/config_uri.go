package config

import (
	"fmt"
	"net/url"
	"strings"
)

// buildUserInfo constructs user credentials for connection URI.
func (c *Config) buildUserInfo() *url.Userinfo {
	if c.Username != "" || c.Password != "" {
		return url.UserPassword(c.Username, c.Password)
	}
	return nil
}

// resolvePort returns effective connection port for URI.
func (c *Config) resolvePort() int {
	if c.IsFirestore() {
		return c.resolveFirestorePort()
	}
	return c.Port
}

func (c *Config) isHostWithPort() bool {
	return c.Protocol == ProtocolMongoDBSrv || strings.Contains(c.Host, ":")
}

// buildHost constructs host and optional port for connection URI.
func (c *Config) buildHost() string {
	if c.isHostWithPort() || c.resolvePort() == 0 {
		return c.Host
	}
	return fmt.Sprintf("%s:%d", c.Host, c.resolvePort())
}

// appendAuthSourceQuery adds optional authSource parameter.
func (c *Config) appendAuthSourceQuery(query string) string {
	if c.AuthSource != "" {
		return fmt.Sprintf("%s&authSource=%s", query, url.QueryEscape(c.AuthSource))
	}
	return query
}

// appendAppNameQuery adds optional appName parameter.
func (c *Config) appendAppNameQuery(query string) string {
	if c.AppName != "" {
		return fmt.Sprintf("%s&appName=%s", query, url.QueryEscape(c.AppName))
	}
	return query
}

// resolveUUIDRep returns effective UUID representation.
func (c *Config) resolveUUIDRep() string {
	if c.UUIDRepresentation == "" {
		return UUIDRepresentationUnspecified
	}
	return c.UUIDRepresentation
}

// buildQuery constructs query string for connection URI.
func (c *Config) buildQuery(enableTLS bool) string {
	if c.IsFirestore() {
		return c.buildFirestoreQuery()
	}
	query := fmt.Sprintf("uuidRepresentation=%s&tls=%t", url.QueryEscape(c.resolveUUIDRep()), enableTLS)
	query = c.appendAuthSourceQuery(query)
	return c.appendAppNameQuery(query)
}

// buildURIFromFields constructs a standard MongoDB connection URI string based on
// the configuration parameters and the enableTLS flag.
func (c *Config) buildURIFromFields(enableTLS bool) string {
	proto := c.Protocol
	if proto == "" {
		proto = ProtocolMongoDB
	}
	u := &url.URL{
		Scheme:   proto,
		User:     c.buildUserInfo(),
		Host:     c.buildHost(),
		Path:     "/",
		RawQuery: c.buildQuery(enableTLS),
	}
	return u.String()
}

// BuildURI constructs a standard MongoDB connection URI string based on
// the configuration parameters and the enableTLS flag.
func (c *Config) BuildURI(enableTLS bool) string {
	return c.buildURIFromFields(enableTLS)
}
