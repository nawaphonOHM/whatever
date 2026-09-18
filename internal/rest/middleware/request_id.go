package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HeaderXRequestID is the standard HTTP header key for request tracing.
const HeaderXRequestID = "X-Request-ID"

// RequestIDKey is the gin.Context key where the request ID is stored.
const RequestIDKey = "RequestID"

// RequestID returns a middleware that injects or propagates an
// X-Request-ID header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Set(RequestIDKey, reqID)
		c.Writer.Header().Set(HeaderXRequestID, reqID)

		c.Next()
	}
}

// GetRequestID extracts the request ID from the Gin context, or empty
// string if not found.
func GetRequestID(c *gin.Context) string {
	if val, exists := c.Get(RequestIDKey); exists {
		if reqID, ok := val.(string); ok {
			return reqID
		}
	}
	// Fallback to reading from header directly if middleware was not executed
	return c.GetHeader(HeaderXRequestID)
}
