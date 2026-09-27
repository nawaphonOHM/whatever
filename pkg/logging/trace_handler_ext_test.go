package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraceHandler_NoSpanOrNilContext(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := slog.New(NewTraceHandler(slog.NewJSONHandler(buf, nil)))
	logger.Info("no span message")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "no span message", logMap["msg"])
	assert.Nil(t, logMap["trace_id"])
}

func verifyGroupLog(t *testing.T, data []byte) {
	t.Helper()
	logMap := parseJSONLog(t, data)
	assert.Equal(t, "auth", logMap["component"])
	meta, ok := logMap["meta"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", meta["trace_id"])
}

func TestTraceHandler_WithAttrsAndGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	th := NewTraceHandler(slog.NewJSONHandler(buf, nil))
	withAttr := th.WithAttrs([]slog.Attr{slog.String("component", "auth")})
	logger := slog.New(withAttr.WithGroup("meta"))

	ctx := createTestSpanContext(t)
	logger.InfoContext(ctx, "grouped message", "attempt", 1)

	verifyGroupLog(t, buf.Bytes())
}

func TestTraceHandler_EnabledAndUnwrap(t *testing.T) {
	base := slog.NewJSONHandler(&bytes.Buffer{}, &slog.HandlerOptions{Level: slog.LevelWarn})
	th := NewTraceHandler(base)

	assert.False(t, th.Enabled(context.Background(), slog.LevelInfo))
	assert.True(t, th.Enabled(context.Background(), slog.LevelWarn))
	assert.Equal(t, base, th.Unwrap())
}

func TestTraceHandler_NilSafety_NilHandler(t *testing.T) {
	var th *TraceHandler
	assert.False(t, th.Enabled(context.Background(), slog.LevelInfo))
	assert.NoError(t, th.Handle(context.Background(), slog.Record{}))
	assert.Nil(t, th.WithAttrs(nil))
	assert.Nil(t, th.WithGroup("test"))
	assert.Nil(t, th.Unwrap())
}

func TestTraceHandler_NilSafety_NilInnerHandler(t *testing.T) {
	thNilInner := NewTraceHandler(nil)
	assert.False(t, thNilInner.Enabled(context.Background(), slog.LevelInfo))
	assert.NoError(t, thNilInner.Handle(context.Background(), slog.Record{}))
	assert.NotNil(t, thNilInner.WithAttrs(nil))
	assert.NotNil(t, thNilInner.WithGroup("test"))
	assert.Nil(t, thNilInner.Unwrap())
}
