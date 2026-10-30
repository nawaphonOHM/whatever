// Package logger provides structured HTTP request logging middleware.
package logger

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// HeaderXRequestID is the standard HTTP header key for request tracing.
const HeaderXRequestID = "X-Request-ID"

// RequestIDKey is the gin.Context key where the request ID is stored.
const RequestIDKey = "RequestID"

const (
	statusServerError = 500
	statusClientError = 400
	httpRequestMsg    = "HTTP Request"
)

// GetRequestID extracts the request ID from the Gin context, or empty string.
func GetRequestID(c *gin.Context) string {
	if val, exists := c.Get(RequestIDKey); exists {
		if reqID, ok := val.(string); ok {
			return reqID
		}
	}
	return c.GetHeader(HeaderXRequestID)
}

// Config defines options for the structured Logger middleware.
type Config struct {
	Logger    *slog.Logger
	SkipPaths []string
}

// Logger returns a structured logging middleware using the default logger.
func Logger() gin.HandlerFunc {
	return WithConfig(&Config{Logger: slog.Default()})
}

// WithLogger returns a logging middleware using a specific slog.Logger.
func WithLogger(logger *slog.Logger) gin.HandlerFunc {
	return WithConfig(&Config{Logger: logger})
}

// buildSkipMap creates a lookup map for skipped paths.
func buildSkipMap(paths []string) map[string]bool {
	skipMap := make(map[string]bool, len(paths))
	for _, p := range paths {
		skipMap[p] = true
	}
	return skipMap
}

// determineLogLevel returns appropriate slog.Level for given HTTP status.
func determineLogLevel(status int) slog.Level {
	if status >= statusServerError {
		return slog.LevelError
	}
	if status >= statusClientError {
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

// logRequest writes the log entry at appropriate level according to status.
func logRequest(
	ctx context.Context,
	logger *slog.Logger,
	status int,
	attrs []slog.Attr,
) {
	level := determineLogLevel(status)
	logger.LogAttrs(ctx, level, httpRequestMsg, attrs...)
}

func resolveLoggerConfig(cfg *Config) (*slog.Logger, map[string]bool) {
	if cfg == nil {
		return slog.Default(), make(map[string]bool)
	}
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return logger, buildSkipMap(cfg.SkipPaths)
}

// WithConfig returns a structured logging middleware configured with options.
func WithConfig(cfg *Config) gin.HandlerFunc {
	logger, skipMap := resolveLoggerConfig(cfg)

	return func(c *gin.Context) {
		if skipMap[c.Request.URL.Path] {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		attrs := buildLogAttrsForHandler(c, latency, logger.Handler())
		logRequest(c.Request.Context(), logger, c.Writer.Status(), attrs)
	}
}
