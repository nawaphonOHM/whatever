package core_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

func TestIsGCPHandler_DirectAndWrapped(t *testing.T) {
	gcpHandler := core.NewGCPHandler(io.Discard, nil)
	assert.True(t, core.IsGCPHandler(gcpHandler))

	traceWrapped := core.NewTraceHandler(gcpHandler)
	assert.True(t, core.IsGCPHandler(traceWrapped))

	nestedTrace := core.NewTraceHandler(traceWrapped)
	assert.True(t, core.IsGCPHandler(nestedTrace))

	withAttrs := gcpHandler.WithAttrs([]slog.Attr{slog.String("key", "val")})
	assert.True(t, core.IsGCPHandler(withAttrs))

	withGroup := gcpHandler.WithGroup("test_group")
	assert.True(t, core.IsGCPHandler(withGroup))
}

func TestIsGCPHandler_NonGCPHandlers(t *testing.T) {
	assert.False(t, core.IsGCPHandler(nil))

	textHandler := slog.NewTextHandler(io.Discard, nil)
	assert.False(t, core.IsGCPHandler(textHandler))

	jsonHandler := slog.NewJSONHandler(io.Discard, nil)
	assert.False(t, core.IsGCPHandler(jsonHandler))

	traceWrappedStd := core.NewTraceHandler(jsonHandler)
	assert.False(t, core.IsGCPHandler(traceWrappedStd))
}
