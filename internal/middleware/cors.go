// Package middleware provides HTTP middlewares for request tracking, recovery, and CORS.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig defines configuration options for CORS middleware.
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           string
}

// DefaultCORSConfig returns a permissive CORS configuration suitable for APIs.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", HeaderXRequestID},
		ExposeHeaders:    []string{"Content-Length", HeaderXRequestID},
		AllowCredentials: false,
		MaxAge:           "86400",
	}
}

// CORS returns a Cross-Origin Resource Sharing middleware.
func CORS(cfg CORSConfig) gin.HandlerFunc {
	allowOrigins := strings.Join(cfg.AllowOrigins, ", ")
	allowMethods := strings.Join(cfg.AllowMethods, ", ")
	allowHeaders := strings.Join(cfg.AllowHeaders, ", ")
	exposeHeaders := strings.Join(cfg.ExposeHeaders, ", ")

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" || allowOrigins == "*" {
			if allowOrigins == "*" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				for _, allowed := range cfg.AllowOrigins {
					if allowed == "*" || allowed == origin {
						c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			if allowMethods != "" {
				c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethods)
			}
			if allowHeaders != "" {
				c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders)
			}
			if exposeHeaders != "" {
				c.Writer.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
			}
			if cfg.AllowCredentials {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if cfg.MaxAge != "" {
				c.Writer.Header().Set("Access-Control-Max-Age", cfg.MaxAge)
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
