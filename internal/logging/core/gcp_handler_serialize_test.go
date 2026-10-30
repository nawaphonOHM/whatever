package core_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	serializeTestCount = 42
	serializeTestPort  = 8080
	serializeErrKey    = "err"
	serializeCountKey  = "count"
)

type nilPointerError struct{ message string }

func (e *nilPointerError) Error() string {
	return e.message
}

type failingMarshalError struct{}

func (failingMarshalError) Error() string {
	return "marshal failed"
}

func (failingMarshalError) MarshalJSON() ([]byte, error) {
	return nil, errors.New("invalid JSON")
}

type panickingMarshalError struct{}

func (panickingMarshalError) Error() string {
	panic("error method failed")
}

func (panickingMarshalError) MarshalJSON() ([]byte, error) {
	panic("marshal method failed")
}

func TestGCPHandler_TypedNilErrorPreservesSiblings(t *testing.T) {
	var buf bytes.Buffer
	var nilErr *nilPointerError
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("nil error", serializeErrKey, nilErr, serializeCountKey, serializeTestCount)
	entry := parseLogMap(t, buf.Bytes())
	assert.Nil(t, entry[serializeErrKey])
	assert.Equal(t, float64(serializeTestCount), entry[serializeCountKey])
}

func TestGCPHandler_FailingMarshalerPreservesSiblings(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("broken error", serializeErrKey, failingMarshalError{}, serializeCountKey, serializeTestCount,
		slog.Group("nested", "port", serializeTestPort, "valid", json.RawMessage(`{"ok":true}`)))
	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, "marshal failed", entry[serializeErrKey])
	assert.Equal(t, float64(serializeTestCount), entry[serializeCountKey])
	nested, ok := entry["nested"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(serializeTestPort), nested["port"])
	assert.Equal(t, map[string]any{"ok": true}, nested["valid"])
}

func TestGCPHandler_PanickingMethodsPreserveSiblings(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("broken methods", serializeErrKey, panickingMarshalError{}, serializeCountKey, serializeTestCount)
	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, "!PANIC: error method failed", entry[serializeErrKey])
	assert.Equal(t, float64(serializeTestCount), entry[serializeCountKey])
}
