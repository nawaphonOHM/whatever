package core_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

type invalidJSONMarshaler struct{}

func (invalidJSONMarshaler) MarshalJSON() ([]byte, error) {
	return []byte("{"), nil
}

type panickingJSONMarshaler struct{}

func (panickingJSONMarshaler) MarshalJSON() ([]byte, error) {
	panic("JSON failed")
}

type failingJSONMarshaler struct{}

func (failingJSONMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New("custom marshal error")
}

type countingJSONMarshaler struct{ calls *int }

func (m countingJSONMarshaler) MarshalJSON() ([]byte, error) {
	*m.calls++
	return []byte(`{"ok":true}`), nil
}

func TestGCPHandler_NonErrorMarshalerFallbackPreservesSiblings(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("broken marshaler", "invalid", invalidJSONMarshaler{}, "panic", panickingJSONMarshaler{},
		"fail", failingJSONMarshaler{}, serializeCountKey, serializeTestCount)
	entry := parseLogMap(t, buf.Bytes())
	expectedInvalid := "!ERROR:json: error calling MarshalJSON for type *core_test.invalidJSONMarshaler: " +
		"unexpected end of JSON input"
	assert.Equal(t, expectedInvalid, entry["invalid"])
	assert.Equal(t, "!PANIC: JSON failed", entry["panic"])
	expectedFail := "!ERROR:json: error calling MarshalJSON for type *core_test.failingJSONMarshaler: " +
		"custom marshal error"
	assert.Equal(t, expectedFail, entry["fail"])
	assert.Equal(t, float64(serializeTestCount), entry[serializeCountKey])
}

func TestGCPHandler_ValidMarshalerCalledOnce(t *testing.T) {
	var buf bytes.Buffer
	calls := 0
	marshaler := countingJSONMarshaler{calls: &calls}
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("valid marshaler", "custom", marshaler)
	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, map[string]any{"ok": true}, entry["custom"])
	assert.Equal(t, 1, calls)
}
