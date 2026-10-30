package core_test

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDecimalBase = 10

func parseLogMapNumbers(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var entry map[string]any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	err := dec.Decode(&entry)
	require.NoError(t, err)
	return entry
}

func matchSigned(t *testing.T, expected any, num json.Number) bool {
	t.Helper()
	if n, ok := expected.(int64); ok {
		assert.Equal(t, strconv.FormatInt(n, testDecimalBase), num.String())
		return true
	}
	if n, ok := expected.(int); ok {
		assert.Equal(t, strconv.Itoa(n), num.String())
		return true
	}
	return false
}

func matchUnsigned(t *testing.T, expected any, num json.Number) bool {
	t.Helper()
	if u, ok := expected.(uint64); ok {
		assert.Equal(t, strconv.FormatUint(u, testDecimalBase), num.String())
		return true
	}
	return false
}

func matchFloatNumber(t *testing.T, expected any, num json.Number) bool {
	t.Helper()
	if exp, ok := expected.(float64); ok {
		f, err := num.Float64()
		assert.NoError(t, err)
		assert.Equal(t, exp, f)
		return true
	}
	return false
}

func matchIntegerNumber(t *testing.T, expected any, num json.Number) bool {
	t.Helper()
	return matchSigned(t, expected, num) || matchUnsigned(t, expected, num)
}

func matchJSONNumber(t *testing.T, expected any, num json.Number) {
	t.Helper()
	if matchIntegerNumber(t, expected, num) {
		return
	}
	if matchFloatNumber(t, expected, num) {
		return
	}
	assert.Equal(t, expected, num.String())
}

func matchPreservedDuration(t *testing.T, expected, actual any) bool {
	t.Helper()
	if d, ok := expected.(time.Duration); ok {
		if s, ok := actual.(string); ok {
			assert.Equal(t, d.String(), s)
			return true
		}
	}
	return false
}

func matchPreservedValue(t *testing.T, expected, actual any) {
	t.Helper()
	if matchPreservedDuration(t, expected, actual) {
		return
	}
	if num, ok := actual.(json.Number); ok {
		matchJSONNumber(t, expected, num)
		return
	}
	assert.Equal(t, expected, actual)
}
