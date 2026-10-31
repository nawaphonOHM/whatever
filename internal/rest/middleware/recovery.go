package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/problem"
)

// Recovery returns a middleware that recovers from any panics,
// logs the error and stack trace, and returns a structured 500 JSON response.
func Recovery() gin.HandlerFunc {
	return RecoveryWithLogger(nil)
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
		slog.String("error", fmt.Sprint(r)),
		slog.String("stack_trace", stack),
	)

	c.Abort()
	resp := problem.New(
		http.StatusInternalServerError,
		"INTERNAL_SERVER_ERROR",
		"An unexpected internal server error occurred",
	)
	resp.Write(c)
}

func resolveRecoveryLogger(logger *slog.Logger) *slog.Logger {
	if logger != nil {
		return logger
	}
	return core.Default().Slog()
}

// RecoveryWithLogger returns a panic recovery middleware using custom logger.
func RecoveryWithLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				l := resolveRecoveryLogger(logger)
				logAndAbortPanic(c, l, r)
			}
		}()

		c.Next()
	}
}
