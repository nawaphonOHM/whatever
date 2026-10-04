package core

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

const (
	testPaymentAmount = 100
	testLineCount     = 6
	testMsg           = "msg"
	testNewline       = "\n"
	testFiveCount     = 5
)

func createTestSpanContext(t *testing.T) context.Context {
	t.Helper()
	traceID, err := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex("00f067aa0ba902b7")
	require.NoError(t, err)

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(context.Background(), sc)
}

func parseJSONLog(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var result map[string]any
	require.NoError(t, json.Unmarshal(data, &result))
	return result
}

type levelFilterCase struct {
	logAction     func(*Logger)
	name          string
	expectedLevel string
	cfgLevel      config.Level
	expectLog     bool
}

func levelFilterCases() []*levelFilterCase {
	return []*levelFilterCase{
		{
			logAction: func(l *Logger) { l.Trace("trace") },
			name:      "trace suppressed at info",
			cfgLevel:  config.LevelInfo,
			expectLog: false,
		},
		{
			logAction:     func(l *Logger) { l.Info("info") },
			name:          "info emitted at info",
			expectedLevel: "INFO",
			cfgLevel:      config.LevelInfo,
			expectLog:     true,
		},
		{
			logAction:     func(l *Logger) { l.Trace("trace") },
			name:          "trace emitted at trace",
			expectedLevel: "TRACE",
			cfgLevel:      config.LevelTrace,
			expectLog:     true,
		},
		{
			logAction:     func(l *Logger) { l.Warn("warn") },
			name:          "warn emitted at warn",
			expectedLevel: "WARN",
			cfgLevel:      config.LevelWarn,
			expectLog:     true,
		},
	}
}
