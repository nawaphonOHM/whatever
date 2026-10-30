package core_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const serializeNestedKey = "nested"

func TestGCPHandler_StoredAndGroupedFallbackPreservesTypes(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil)).
		With(serializeErrKey, failingMarshalError{}, serializeCountKey, serializeTestCount).
		WithGroup(serializeNestedKey).With(serializeErrKey, (*nilPointerError)(nil), "port", serializeTestPort)
	logger.Info("stored and grouped")
	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, "marshal failed", entry[serializeErrKey])
	assert.Equal(t, float64(serializeTestCount), entry[serializeCountKey])
	nested, ok := entry[serializeNestedKey].(map[string]any)
	require.True(t, ok)
	assert.Nil(t, nested[serializeErrKey])
	assert.Equal(t, float64(serializeTestPort), nested["port"])
}

func TestGCPHandler_GroupedFailingMarshalerPreservesTypes(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("grouped error", slog.Group(serializeNestedKey, serializeErrKey, failingMarshalError{},
		serializeCountKey, serializeTestCount))
	entry := parseLogMap(t, buf.Bytes())
	nested, ok := entry[serializeNestedKey].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "marshal failed", nested[serializeErrKey])
	assert.Equal(t, float64(serializeTestCount), nested[serializeCountKey])
}
