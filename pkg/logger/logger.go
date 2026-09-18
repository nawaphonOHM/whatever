// Package logger provides structured HTTP request logging middleware.
package logger

import (
	"context"
	"log/slog"
	"os"
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
	return WithConfig(Config{Logger: slog.Default()})
}

// WithLogger returns a logging middleware using a specific slog.Logger.
func WithLogger(logger *slog.Logger) gin.HandlerFunc {
	return WithConfig(Config{Logger: logger})
}

// buildSkipMap creates a lookup map for skipped paths.
func buildSkipMap(paths []string) map[string]bool {
	skipMap := make(map[string]bool, len(paths))
	for _, p := range paths {
		skipMap[p] = true
	}
	return skipMap
}

// buildLogAttrs constructs standard log attributes for a request.
func buildLogAttrs(c *gin.Context, latency time.Duration) []slog.Attr {
	attrs := []slog.Attr{
		slog.String("request_id", GetRequestID(c)),
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.String("query", c.Request.URL.RawQuery),
		slog.Int("status", c.Writer.Status()),
		slog.Int64("latency_ms", latency.Milliseconds()),
		slog.Duration("latency", latency),
		slog.String("client_ip", c.ClientIP()),
		slog.String("user_agent", c.Request.UserAgent()),
		slog.Int("bytes_out", c.Writer.Size()),
	}
	if len(c.Errors) > 0 {
		attrs = append(attrs, slog.String("errors", c.Errors.String()))
	}
	return attrs
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

// WithConfig returns a structured logging middleware configured with options.
func WithConfig(cfg Config) gin.HandlerFunc {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	skipMap := buildSkipMap(cfg.SkipPaths)

	return func(c *gin.Context) {
		if skipMap[c.Request.URL.Path] {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		attrs := buildLogAttrs(c, latency)
		logRequest(c.Request.Context(), logger, c.Writer.Status(), attrs)
	}
}
