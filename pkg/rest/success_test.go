package rest

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Constants for fixed timestamp fields and integer payload.
const (
	testYear   = 2026
	testMonth  = 9
	testDay    = 17
	testHour   = 2
	testMinute = 30
	testIntVal = 123
)

// getFixedTime returns a deterministic UTC timestamp for serialization tests.
func getFixedTime() time.Time {
	return time.Date(
		testYear, testMonth, testDay, testHour, testMinute, 0, 0, time.UTC,
	)
}

// TestSuccessResponse_ShapeSerialization tests JSON envelope serialization.
func TestSuccessResponse_ShapeSerialization(t *testing.T) {
	fixedTime := getFixedTime()
	resp := NewSuccessResponse("payload", "ok").WithTimestamp(fixedTime)

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))

	assert.Equal(t, true, raw["success"])
	assert.Equal(t, "ok", raw["message"])
	assert.Equal(t, "payload", raw["data"])
	assert.Equal(t, "2026-09-17T02:30:00Z", raw["timestamp"])
}

// TestSuccessResponse_ChainedMethods tests chaining WithMessage and
// WithTimestamp.
func TestSuccessResponse_ChainedMethods(t *testing.T) {
	fixedTime := getFixedTime()
	resp := NewSuccessResponse(testIntVal).
		WithMessage("custom").
		WithTimestamp(fixedTime)

	assert.Equal(t, "custom", resp.Message)
	assert.Equal(t, testIntVal, resp.Data)
	assert.Equal(t, fixedTime, resp.Timestamp)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
