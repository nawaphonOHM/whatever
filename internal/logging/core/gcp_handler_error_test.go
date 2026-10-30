package core_test

import (
	"bytes"
	"errors"
	"log/slog"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

func TestGCPHandler_ErrorAttribute(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, nil)
	logger := slog.New(handler)

	testErr := errors.New("connection reset by peer")
	logger.Error("request failed", "err", testErr)

	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, core.GCPSeverityError, entry[testKeySeverity])
	assert.Equal(t, "connection reset by peer", entry["err"])
}

func verifySpecialFloats(t *testing.T, data []byte) {
	t.Helper()
	entry := parseLogMap(t, data)
	nanStr, ok := entry["nan_val"].(string)
	require.True(t, ok)
	assert.Contains(t, nanStr, "NaN")

	infStr, ok := entry["inf_val"].(string)
	require.True(t, ok)
	assert.Contains(t, infStr, "Inf")
}

func TestGCPHandler_SpecialFloatValues(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))

	logger.Info("metrics", "nan_val", math.NaN(), "inf_val", math.Inf(1))
	verifySpecialFloats(t, buf.Bytes())
}

func TestGCPHandler_UnserializableFallback(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))

	ch := make(chan int)
	logger.Info("chan test", "channel", ch)

	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, "chan test", entry[testKeyMessage])
	assert.Equal(t, "!ERROR:json: unsupported type: chan int", entry["channel"])
}
