package core

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"sync"
)

// GCPHandler formats logs into Google Cloud Logging LogEntry JSON structure.
type GCPHandler struct {
	out      io.Writer
	mu       *sync.Mutex
	opts     slog.HandlerOptions
	payload  map[string]any
	topAttrs []slog.Attr
	groups   []string
}

// NewGCPHandler constructs a new GCPHandler writing to the specified writer.
func NewGCPHandler(out io.Writer, opts *slog.HandlerOptions) *GCPHandler {
	targetOut := out
	if targetOut == nil {
		targetOut = os.Stdout
	}
	var hOpts slog.HandlerOptions
	if opts != nil {
		hOpts = *opts
	}
	return &GCPHandler{
		out:     targetOut,
		mu:      &sync.Mutex{},
		opts:    hOpts,
		payload: make(map[string]any),
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *GCPHandler) Enabled(_ctx context.Context, level slog.Level) bool {
	_ = _ctx
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

type handlerUnwrapper interface {
	Unwrap() slog.Handler
}

func unwrapHandler(h slog.Handler) slog.Handler {
	if u, ok := h.(handlerUnwrapper); ok {
		return u.Unwrap()
	}
	return nil
}

// IsGCPHandler reports whether the given handler is or wraps a GCPHandler.
func IsGCPHandler(h slog.Handler) bool {
	for curr := h; curr != nil; curr = unwrapHandler(curr) {
		if _, ok := curr.(*GCPHandler); ok {
			return true
		}
	}
	return false
}

// Handle serializes the slog.Record as GCP Cloud Logging JSON and writes atomically.
func (h *GCPHandler) Handle(_ctx context.Context, record slog.Record) error {
	_ = _ctx
	entry := h.buildEntry(record)
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err = h.out.Write(data)
	return err
}
