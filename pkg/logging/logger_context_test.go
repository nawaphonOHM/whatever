package logging

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPaymentAmount = 100
	testLineCount     = 6
)

func TestLogger_WithAndWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(Config{Output: buf, Level: LevelDebug})

	child := l.With("service", "billing").WithGroup("req")
	child.Debug("payment processed", "amount", testPaymentAmount)

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "billing", logMap["service"])
	reqMap, ok := logMap["req"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(testPaymentAmount), reqMap["amount"])
}

func verifySpanLogLines(t *testing.T, data string) {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(data), "\n")
	require.Len(t, lines, 2)
	for _, line := range lines {
		logMap := parseJSONLog(t, []byte(line))
		assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", logMap["trace_id"])
		assert.Equal(t, "00f067aa0ba902b7", logMap["span_id"])
	}
}

func TestLogger_TraceContextWithSpan(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(Config{Output: buf, Level: LevelTrace})

	ctx := createTestSpanContext(t)
	l.TraceContext(ctx, "trace with context")
	l.InfoContext(ctx, "info with context")

	verifySpanLogLines(t, buf.String())
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
	SetDefault(New(Config{Output: buf, Level: LevelTrace}))

	ctx := createTestSpanContext(t)
	emitPackageLevelLogs(ctx)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, testLineCount)
}
