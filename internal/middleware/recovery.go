package middleware

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/nawaphonOHM/whatever/pkg/rest/response"
	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from any panics, logs the error and stack trace,
// and returns a structured 500 Internal Server Error JSON response.
func Recovery() gin.HandlerFunc {
	return RecoveryWithLogger(slog.Default())
}

// RecoveryWithLogger returns a panic recovery middleware using a custom slog.Logger.
func RecoveryWithLogger(logger *slog.Logger) gin.HandlerFunc {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				reqID := GetRequestID(c)
				stack := string(debug.Stack())

				logger.ErrorContext(
					c.Request.Context(),
					"panic recovered during request processing",
					slog.String("request_id", reqID),
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
					slog.Any("panic", r),
					slog.String("stack", stack),
				)

				c.Abort()
				resp := response.InternalServerError("INTERNAL_SERVER_ERROR", "An unexpected internal server error occurred")
				resp.Write(c)
			}
		}()

		c.Next()
	}
}
