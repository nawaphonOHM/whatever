package core_test

import (
	"bytes"
	"log/slog"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const testStructAValue = 5

type panickingErrorOnly struct{}

func (panickingErrorOnly) Error() string {
	panic("error method panic")
}

type structWithFunc struct {
	Fn func()
	A  int
}

func verifyUnsupportedEntries(t *testing.T, entry map[string]any) {
	t.Helper()
	assert.Equal(t, "!ERROR:json: unsupported type: func()", entry["fn"])
	assert.Equal(t, "!ERROR:json: unsupported type: chan int", entry["map"])
	assert.Equal(t, "!ERROR:json: unsupported value: NaN", entry["slice"])
	assert.Equal(t, "!ERROR:json: unsupported type: func()", entry["struct"])
	assert.Equal(t, float64(serializeTestCount), entry[serializeCountKey])
}

func TestGCPHandler_UnsupportedFunctionsAndContainers(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	mapVal := map[string]any{"x": 1, "c": make(chan int)}
	sliceVal := []float64{math.NaN(), 1}
	structVal := structWithFunc{A: testStructAValue, Fn: func() {}}

	logger.Info("unsupported", "fn", func() {}, "map", mapVal, "slice", sliceVal, "struct", structVal,
		serializeCountKey, serializeTestCount)
	entry := parseLogMap(t, buf.Bytes())
	verifyUnsupportedEntries(t, entry)
}

func TestGCPHandler_PanickingErrorOnlyPreservesSiblings(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("panicking error", serializeErrKey, panickingErrorOnly{}, serializeCountKey, serializeTestCount)
	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, "!PANIC: error method panic", entry[serializeErrKey])
	assert.Equal(t, float64(serializeTestCount), entry[serializeCountKey])
}
