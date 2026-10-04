package core

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func verifySpanLogLines(t *testing.T, data string) {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(data), testNewline)
	require.Len(t, lines, 2)
	for _, line := range lines {
		logMap := parseJSONLog(t, []byte(line))
		assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", logMap["trace_id"])
		assert.Equal(t, "00f067aa0ba902b7", logMap["span_id"])
	}
}

func TestLogger_TraceContextWithSpan(t *testing.T) {
	buf := &bytes.Buffer{}
	l := NewJSON(buf, config.LevelTrace)

	ctx := createTestSpanContext(t)
	l.TraceContext(ctx, "trace with context")
	l.InfoContext(ctx, "info with context")

	verifySpanLogLines(t, buf.String())
}

func TestLogger_DisableTraceCorrelation(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := &config.Config{
		Format:                  string(config.FormatJSON),
		Level:                   string(config.LevelInfo),
		DisableTraceCorrelation: true,
	}
	h := BuildHandler(cfg, buf, nil)
	l := NewWithHandler(h)

	ctx := createTestSpanContext(t)
	l.InfoContext(ctx, "hello without trace correlation")

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "hello without trace correlation", logMap[testMsg])
	assert.Nil(t, logMap["trace_id"])
	assert.Nil(t, logMap["span_id"])
}

func emitPackageLevelLogs(ctx context.Context) {
	Trace("global trace")
	TraceContext(ctx, "global trace context")
	Debug("global debug")
	Info("global info")
	Warn("global warn")
	Error("global error")
}

func TestPackageLevel_Functions(t *testing.T) {
	buf := &bytes.Buffer{}
	cleanup := CaptureLogs(buf, config.LevelTrace)
	defer cleanup()

	ctx := createTestSpanContext(t)
	emitPackageLevelLogs(ctx)

	lines := strings.Split(strings.TrimSpace(buf.String()), testNewline)
	assert.Len(t, lines, testLineCount)
}

func emitGlobalContextLogs(ctx context.Context) {
	DebugContext(ctx, "debug ctx msg")
	InfoContext(ctx, "info ctx msg")
	WarnContext(ctx, "warn ctx msg")
	ErrorContext(ctx, "error ctx msg")
}

func TestGlobal_ContextFunctions(t *testing.T) {
	buf := &bytes.Buffer{}
	cleanup := CaptureLogs(buf, config.LevelDebug)
	defer cleanup()

	emitGlobalContextLogs(context.Background())
	out := buf.String()
	assert.Contains(t, out, "debug ctx msg")
	assert.Contains(t, out, "info ctx msg")
	assert.Contains(t, out, "warn ctx msg")
	assert.Contains(t, out, "error ctx msg")
}
