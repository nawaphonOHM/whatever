package core_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

func verifyStringMsg(t *testing.T, input string) {
	t.Helper()
	got := core.ParseGCPMessage(input)
	str, ok := got.(string)
	require.True(t, ok)
	assert.Equal(t, input, str)
}

func verifyJSONMsg(t *testing.T, input string) {
	t.Helper()
	got := core.ParseGCPMessage(input)
	raw, ok := got.(json.RawMessage)
	require.True(t, ok)
	assert.True(t, json.Valid(raw))
}

func TestParseGCPMessageString(t *testing.T) {
	tests := []string{
		"plain text message",
		"12345",
		"{ incomplete json",
		"{ invalid: json }",
		"[ invalid array",
		"",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			verifyStringMsg(t, tt)
		})
	}
}

func TestParseGCPMessageJSON(t *testing.T) {
	tests := []string{
		`{"action": "login", "user": "alice"}`,
		`[1, 2, 3]`,
		`  {"nested": {"key": true}}  `,
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			verifyJSONMsg(t, tt)
		})
	}
}
