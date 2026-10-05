package mongodb

import (
	"fmt"
	"net/url"
)

// resolveURI constructs the appropriate connection URI string from Options.
func resolveURI(o *Options) string {
	if hasURI(o) {
		return resolveURIFromOptions(o)
	}
	return BuildURI(o, resolveEnableTLS(o))
}

func resolveURIFromOptions(o *Options) string {
	if o == nil {
		return ""
	}
	if o.EnableTLS {
		return injectTLSQueryParam(o.URI)
	}
	return o.URI
}

func resolveEnableTLS(o *Options) bool {
	return o != nil && (o.EnableTLS || o.IsFirestore())
}

// injectTLSQueryParam adds or overwrites the tls=true parameter in the URI.
func injectTLSQueryParam(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	if u.Path == "" {
		u.Path = "/"
	}
	q := u.Query()
	q.Set("tls", "true")
	u.RawQuery = q.Encode()
	return u.String()
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
