package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

func appendContextAttrs(c *gin.Context, base []slog.Attr) []slog.Attr {
	res := append([]slog.Attr(nil), base...)
	res = append(res, traceAttrs(c)...)
	if len(c.Errors) > 0 {
		res = append(res, slog.String("errors", c.Errors.String()))
	}
	return res
}

func buildStandardBaseAttrs(c *gin.Context, latency time.Duration) []slog.Attr {
	return []slog.Attr{
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
}

func requestURL(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	return req.URL.RequestURI()
}

func buildGCPBaseAttrs(c *gin.Context, latency time.Duration) []slog.Attr {
	attrs := []slog.Attr{
		slog.String("request_id", GetRequestID(c)),
		slog.String("method", c.Request.Method),
		slog.String("request_url", requestURL(c.Request)),
		slog.Int("status", c.Writer.Status()),
		slog.Duration("latency", latency),
		slog.String("client_ip", c.ClientIP()),
		slog.String("user_agent", c.Request.UserAgent()),
	}
	if c.Writer.Size() >= 0 {
		attrs = append(attrs, slog.Int("bytes_out", c.Writer.Size()))
	}
	if c.Request.Proto != "" {
		attrs = append(attrs, slog.String("protocol", c.Request.Proto))
	}
	return attrs
}

func selectBaseAttrs(c *gin.Context, latency time.Duration, h slog.Handler) []slog.Attr {
	if core.IsGCPHandler(h) {
		return buildGCPBaseAttrs(c, latency)
	}
	return buildStandardBaseAttrs(c, latency)
}

func buildLogAttrsForHandler(c *gin.Context, latency time.Duration, h slog.Handler) []slog.Attr {
	if c == nil || c.Request == nil {
		return nil
	}
	return appendContextAttrs(c, selectBaseAttrs(c, latency, h))
}
