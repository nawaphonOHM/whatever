package core_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

func verifyMixedMsg(t *testing.T, input, expectedHeader, expectedJSON string) {
	t.Helper()
	got := core.ParseGCPMessage(input)
	mixed, ok := got.(core.MixedMessage)
	require.True(t, ok, "expected MixedMessage for input: %q", input)
	assert.Equal(t, expectedHeader, mixed.Header)
	assert.True(t, json.Valid(mixed.Content))
	assert.JSONEq(t, expectedJSON, string(mixed.Content))
}

type mixedTestCase struct {
	name           string
	input          string
	expectedHeader string
	expectedJSON   string
}

func mixedTestCases() []mixedTestCase {
	return []mixedTestCase{
		{
			name:           "colon separated object",
			input:          `received response: {"status": "ok"}`,
			expectedHeader: "received response:",
			expectedJSON:   `{"status": "ok"}`,
		},
		{
			name:           "space separated array",
			input:          `status [1, 2, 3]`,
			expectedHeader: "status",
			expectedJSON:   `[1, 2, 3]`,
		},
		{
			name:           "space separated nested object",
			input:          `event {"nested":{"k":"v"}}`,
			expectedHeader: "event",
			expectedJSON:   `{"nested":{"k":"v"}}`,
		},
		{
			name:           "payload with id",
			input:          `payload: {"action":"create","id":123}`,
			expectedHeader: "payload:",
			expectedJSON:   `{"action":"create","id":123}`,
		},
		{
			name:           "braces inside header text",
			input:          `prefix with {brace} inside: {"valid": true}`,
			expectedHeader: "prefix with {brace} inside:",
			expectedJSON:   `{"valid": true}`,
		},
		{
			name:           "brackets inside header text",
			input:          `brackets [1] and obj: {"key": "value"}`,
			expectedHeader: "brackets [1] and obj:",
			expectedJSON:   `{"key": "value"}`,
		},
		{
			name:           "colon separated array",
			input:          `error: [1, 2, 3]`,
			expectedHeader: "error:",
			expectedJSON:   `[1, 2, 3]`,
		},
		{
			name:           "nested mixed structures",
			input:          `nested mixed: {"arr": [1, 2], "sub": {"a": true}}`,
			expectedHeader: "nested mixed:",
			expectedJSON:   `{"arr": [1, 2], "sub": {"a": true}}`,
		},
		{
			name:           "whitespace around message",
			input:          "   msg: {\"x\": 1}   \n",
			expectedHeader: "msg:",
			expectedJSON:   `{"x": 1}`,
		},
	}
}

func TestParseGCPMessageMixed(t *testing.T) {
	for _, tt := range mixedTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			verifyMixedMsg(t, tt.input, tt.expectedHeader, tt.expectedJSON)
		})
	}
}
