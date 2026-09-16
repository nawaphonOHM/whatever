package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// HeaderXRequestID is the standard HTTP header key for request tracing.
const HeaderXRequestID = "X-Request-ID"

// RequestIDKey is the gin.Context key where the request ID is stored.
const RequestIDKey = "RequestID"

// GetRequestID extracts the request ID from the Gin context, or empty string if not found.
func GetRequestID(c *gin.Context) string {
	if val, exists := c.Get(RequestIDKey); exists {
		if reqID, ok := val.(string); ok {
			return reqID
		}
	}
	return c.GetHeader(HeaderXRequestID)
}

// LoggerConfig defines the configuration options for the structured Logger middleware.
type LoggerConfig struct {
	Logger    *slog.Logger
	SkipPaths []string
}

// Logger returns a structured logging middleware using the default slog logger.
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Logger: slog.Default(),
	})
}

// LoggerWithLogger returns a structured logging middleware using a specific slog.Logger instance.
func LoggerWithLogger(logger *slog.Logger) gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Logger: logger,
	})
}

// LoggerWithConfig returns a structured logging middleware configured with custom options.
func LoggerWithConfig(cfg LoggerConfig) gin.HandlerFunc {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	skipMap := make(map[string]bool, len(cfg.SkipPaths))
	for _, path := range cfg.SkipPaths {
		skipMap[path] = true
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if skipMap[path] {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		status := c.Writer.Status()
		reqID := GetRequestID(c)

		attrs := []slog.Attr{
			slog.String("request_id", reqID),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.Int("status", status),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.Duration("latency", latency),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Int("bytes_out", c.Writer.Size()),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}

		msg := "HTTP Request"
		ctx := c.Request.Context()

		switch {
		case status >= 500:
			logger.LogAttrs(ctx, slog.LevelError, msg, attrs...)
		case status >= 400:
			logger.LogAttrs(ctx, slog.LevelWarn, msg, attrs...)
		default:
			logger.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
		}
	}
}
