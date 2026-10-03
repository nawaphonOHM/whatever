package core

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWithHandler(t *testing.T) {
	t.Run("valid handler", func(t *testing.T) {
		buf := &bytes.Buffer{}
		h := slog.NewJSONHandler(buf, nil)
		l := NewWithHandler(h)
		require.NotNil(t, l)

		l.Info("message via custom handler")
		logMap := parseJSONLog(t, buf.Bytes())
		assert.Equal(t, "message via custom handler", logMap["msg"])
	})

	t.Run("nil handler defaults safely", func(t *testing.T) {
		l := NewWithHandler(nil)
		require.NotNil(t, l)
		assert.NotNil(t, l.Slog())
	})
}

func TestLogger_WithAndWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	l := NewJSON(buf, config.LevelDebug)

	child := l.With("service", "billing").WithGroup("req")
	child.Debug("payment processed", "amount", testPaymentAmount)

	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "billing", logMap["service"])
	reqMap, ok := logMap["req"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(testPaymentAmount), reqMap["amount"])
}
