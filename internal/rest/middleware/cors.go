// Package middleware provides HTTP middlewares for request tracking,
// recovery, and CORS.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// DefaultCORSConfig returns a disabled CORS configuration with no allowed origins.
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: nil,
	}
}

func resolveCORSConfig(cfg *CORSConfig) *CORSConfig {
	if cfg == nil {
		return DefaultCORSConfig()
	}
	return cfg
}

// CORS returns a Cross-Origin Resource Sharing middleware.
func CORS(cfg *CORSConfig) gin.HandlerFunc {
	resolved := resolveCORSConfig(cfg)
	return func(c *gin.Context) {
		applied := setCORSHeaders(c, resolved)
		if applied && c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
