package middleware

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest"
)

// Recovery returns a middleware that recovers from any panics,
// logs the error and stack trace, and returns a structured 500 JSON response.
func Recovery() gin.HandlerFunc {
	return RecoveryWithLogger(slog.Default())
}

// logAndAbortPanic logs panic details and writes an error response.
func logAndAbortPanic(c *gin.Context, l *slog.Logger, r any) {
	reqID := GetRequestID(c)
	stack := string(debug.Stack())

	l.ErrorContext(
		c.Request.Context(),
		"panic recovered during request processing",
		slog.String("request_id", reqID),
		slog.String("path", c.Request.URL.Path),
		slog.String("method", c.Request.Method),
		slog.Any("panic", r),
		slog.String("stack", stack),
	)

	c.Abort()
	resp := rest.InternalServerError(
		"INTERNAL_SERVER_ERROR",
		"An unexpected internal server error occurred",
	)
	resp.Write(c)
}

// RecoveryWithLogger returns a panic recovery middleware using custom logger.
func RecoveryWithLogger(logger *slog.Logger) gin.HandlerFunc {
	l := logger
	if l == nil {
		l = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}

	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				logAndAbortPanic(c, l, r)
			}
		}()

		c.Next()
	}
}
