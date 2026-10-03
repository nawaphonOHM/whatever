package core

import (
	"context"
	"log/slog"
)

// TraceHandler is an slog.Handler wrapper that enriches records with active
// OpenTelemetry distributed trace_id and span_id attributes.
type TraceHandler struct {
	handler slog.Handler
}

// NewTraceHandler wraps an existing slog.Handler with trace correlation support.
func NewTraceHandler(handler slog.Handler) *TraceHandler {
	return &TraceHandler{handler: handler}
}

// Unwrap returns the underlying wrapped slog.Handler.
func (h *TraceHandler) Unwrap() slog.Handler {
	if h == nil {
		return nil
	}
	return h.handler
}

// Enabled reports whether the handler emits records at the given level.
func (h *TraceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h == nil || h.handler == nil {
		return false
	}
	return h.handler.Enabled(ctx, level)
}

// Handle enriches the record with trace identifiers and delegates to the underlying handler.
func (h *TraceHandler) Handle(ctx context.Context, record slog.Record) error {
	if h == nil || h.handler == nil {
		return nil
	}
	enrichTraceAttrs(ctx, &record)
	return h.handler.Handle(ctx, record)
}

// WithAttrs returns a new TraceHandler with the given attributes added.
func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h == nil || h.handler == nil {
		return h
	}
	return &TraceHandler{handler: h.handler.WithAttrs(attrs)}
}

// WithGroup returns a new TraceHandler with the given group name.
func (h *TraceHandler) WithGroup(name string) slog.Handler {
	if h == nil || h.handler == nil {
		return h
	}
	return &TraceHandler{handler: h.handler.WithGroup(name)}
}
