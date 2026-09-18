// Package middleware provides HTTP middlewares for request tracking,
// recovery, and CORS.
package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	wildcardOrigin = "*"
	commaSep       = ", "
	headerAllowOrg = "Access-Control-Allow-Origin"
	headerAllowMth = "Access-Control-Allow-Methods"
	headerAllowHdr = "Access-Control-Allow-Headers"
	headerExposeHd = "Access-Control-Expose-Headers"
	headerAllowCrd = "Access-Control-Allow-Credentials"
	headerMaxAge   = "Access-Control-Max-Age"
)

// CORSConfig defines configuration options for CORS middleware.
type CORSConfig struct {
	MaxAge           string
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
}

// DefaultCORSConfig returns a permissive CORS configuration suitable for APIs.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{wildcardOrigin},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
			HeaderXRequestID,
		},
		ExposeHeaders:    []string{"Content-Length", HeaderXRequestID},
		AllowCredentials: false,
		MaxAge:           "86400",
	}
}

// setOriginHeader sets the Access-Control-Allow-Origin header on response.
func setOriginHeader(c *gin.Context, origins []string, origin string) {
	if slices.Contains(origins, wildcardOrigin) {
		c.Writer.Header().Set(headerAllowOrg, wildcardOrigin)
		return
	}
	if slices.Contains(origins, origin) {
		c.Writer.Header().Set(headerAllowOrg, origin)
	}
}

// setHeaderIfNonEmpty sets a response header if value is non-empty.
func setHeaderIfNonEmpty(w http.ResponseWriter, key, val string) {
	if val != "" {
		w.Header().Set(key, val)
	}
}

// setConfigHeaders sets configured CORS header values on response writer.
func setConfigHeaders(w http.ResponseWriter, cfg *CORSConfig) {
	setHeaderIfNonEmpty(
		w, headerAllowMth, strings.Join(cfg.AllowMethods, commaSep),
	)
	setHeaderIfNonEmpty(
		w, headerAllowHdr, strings.Join(cfg.AllowHeaders, commaSep),
	)
	setHeaderIfNonEmpty(
		w, headerExposeHd, strings.Join(cfg.ExposeHeaders, commaSep),
	)
	if cfg.AllowCredentials {
		w.Header().Set(headerAllowCrd, "true")
	}
	setHeaderIfNonEmpty(w, headerMaxAge, cfg.MaxAge)
}

// setCORSHeaders sets all CORS headers on response.
func setCORSHeaders(c *gin.Context, cfg *CORSConfig) {
	origin := c.Request.Header.Get("Origin")
	if origin == "" && !slices.Contains(cfg.AllowOrigins, wildcardOrigin) {
		return
	}
	setOriginHeader(c, cfg.AllowOrigins, origin)
	setConfigHeaders(c.Writer, cfg)
}

// CORS returns a Cross-Origin Resource Sharing middleware.
func CORS(cfg CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		setCORSHeaders(c, &cfg)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
