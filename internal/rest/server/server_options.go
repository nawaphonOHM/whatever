package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// applyTrustedProxies configures trusted proxies on gin engine.
func applyTrustedProxies(engine *gin.Engine, proxies []string) error {
	if proxies == nil {
		return nil
	}
	if err := engine.SetTrustedProxies(proxies); err != nil {
		return fmt.Errorf("failed to set trusted proxies: %w", err)
	}
	return nil
}

// applyProxyOptions configures trusted proxy and client IP headers.
func applyProxyOptions(engine *gin.Engine, p ProxyFields) error {
	engine.ForwardedByClientIP = p.ForwardedByClientIP
	if len(p.RemoteIPHeaders) > 0 {
		engine.RemoteIPHeaders = p.RemoteIPHeaders
	}
	return applyTrustedProxies(engine, p.TrustedProxies)
}

// applyHTTPOptions configures HTTP routing and path flags.
func applyHTTPOptions(engine *gin.Engine, h HTTPFields) {
	engine.RedirectTrailingSlash = h.RedirectTrailingSlash
	engine.RedirectFixedPath = h.RedirectFixedPath
	engine.HandleMethodNotAllowed = h.HandleMethodNotAllowed
	engine.UseRawPath = h.UseRawPath
	engine.UnescapePathValues = h.UnescapePathValues
	engine.RemoveExtraSlash = h.RemoveExtraSlash
}

// applyServerOptions applies configured proxy, HTTP, and resource options to engine.
func applyServerOptions(engine *gin.Engine, c *Config) error {
	if err := applyProxyOptions(engine, c.ProxyFields); err != nil {
		return err
	}
	applyHTTPOptions(engine, c.HTTPFields)
	if c.MaxBodySize > 0 {
		engine.MaxMultipartMemory = c.MaxBodySize
	}
	return nil
}
