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
	resp := newSuccessResponse("payload", "ok").withTimestamp(fixedTime)

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))

	assert.Equal(t, true, raw["success"])
	assert.Equal(t, "ok", raw["message"])
	assert.Equal(t, "payload", raw["data"])
	assert.Equal(t, "2026-09-17T02:30:00Z", raw["timestamp"])
}

// assertDTOFields verifies unmarshaling into the public SuccessResponse DTO.
func assertDTOFields(t *testing.T, data []byte, fixedTime time.Time) {
	var dto SuccessResponse[string]
	require.NoError(t, json.Unmarshal(data, &dto))
	assert.True(t, dto.Success)
	assert.Equal(t, "ok", dto.Message)
	assert.Equal(t, "payload", dto.Data)
	assert.Equal(t, fixedTime, dto.Timestamp)
}

// assertEnvelopeFields verifies unmarshaling into the Envelope type alias.
func assertEnvelopeFields(t *testing.T, data []byte, fixedTime time.Time) {
	var env Envelope[string]
	require.NoError(t, json.Unmarshal(data, &env))
	assert.True(t, env.Success)
	assert.Equal(t, "ok", env.Message)
	assert.Equal(t, "payload", env.Data)
	assert.Equal(t, fixedTime, env.Timestamp)
}

// TestSuccessResponse_DTOUnmarshaling tests JSON envelope deserialization
// into the public SuccessResponse and Envelope DTO types.
func TestSuccessResponse_DTOUnmarshaling(t *testing.T) {
	fixedTime := getFixedTime()
	resp := newSuccessResponse("payload", "ok").withTimestamp(fixedTime)
	data, err := json.Marshal(resp)
	require.NoError(t, err)

	assertDTOFields(t, data, fixedTime)
	assertEnvelopeFields(t, data, fixedTime)
}

// TestSuccessResponse_ChainedMethods tests chaining withMessage and
// withTimestamp.
func TestSuccessResponse_ChainedMethods(t *testing.T) {
	fixedTime := getFixedTime()
	resp := newSuccessResponse(testIntVal).
		withMessage("custom").
		withTimestamp(fixedTime)

	assert.Equal(t, "custom", resp.Message)
	assert.Equal(t, testIntVal, resp.Data)
	assert.Equal(t, fixedTime, resp.Timestamp)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
