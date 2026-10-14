package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type schemaTestCase struct {
	name    string
	content string
	wantErr string
}

func runSchemaTestCases(t *testing.T, cases []schemaTestCase) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTempConfigFile(t, "cors.json", tc.content)
			cfg, err := LoadAndValidateCORS(path)
			require.Error(t, err)
			assert.Nil(t, cfg)
			assert.Contains(t, err.Error(), schemaValidationErr)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestLoadAndValidateCORS_MethodAndHeaderViolations(t *testing.T) {
	cases := []schemaTestCase{
		{
			name:    "unknown extra property",
			content: `{"allow-origins": ["https://example.com"], "unknown": "val"}`,
			wantErr: "Additional property unknown is not allowed",
		},
		{
			name:    "invalid HTTP method enum",
			content: `{"allow-methods": ["INVALID_METHOD"]}`,
			wantErr: "allow-methods.0 must be one of the following",
		},
		{
			name:    "PATCH method rejected",
			content: `{"allow-methods": ["PATCH"]}`,
			wantErr: "allow-methods.0 must be one of the following",
		},
		{
			name:    "invalid HTTP header enum",
			content: `{"allow-headers": ["Invalid-Header-XYZ"]}`,
			wantErr: "must be one of the following",
		},
		{
			name:    "lowercase header rejected",
			content: `{"allow-headers": ["content-type"]}`,
			wantErr: "must be one of the following",
		},
		{
			name:    "invalid request-headers enum",
			content: `{"request-headers": ["Bogus-Header-123"]}`,
			wantErr: "must be one of the following",
		},
		{
			name:    "invalid request-methods enum",
			content: `{"request-methods": ["UNKNOWN_METHOD"]}`,
			wantErr: "request-methods.0 must be one of the following",
		},
		{
			name:    "mixed valid and invalid allow-methods at index 1",
			content: `{"allow-methods": ["GET", "PATCH"]}`,
			wantErr: "allow-methods.1 must be one of the following",
		},
		{
			name:    "mixed valid and invalid request-methods at index 1",
			content: `{"request-methods": ["POST", "INVALID_METHOD"]}`,
			wantErr: "request-methods.1 must be one of the following",
		},
		{
			name:    "lowercase allow-methods rejected",
			content: `{"allow-methods": ["get"]}`,
			wantErr: "allow-methods.0 must be one of the following",
		},
		{
			name:    "lowercase request-methods rejected",
			content: `{"request-methods": ["post"]}`,
			wantErr: "request-methods.0 must be one of the following",
		},
	}
	runSchemaTestCases(t, cases)
}

func TestLoadAndValidateCORS_OriginAndTypeViolations(t *testing.T) {
	cases := []schemaTestCase{
		{
			name:    "invalid origin format without scheme",
			content: `{"allow-origins": ["invalid-origin-without-scheme"]}`,
			wantErr: patternMismatch,
		},
		{
			name:    "invalid origin format literal dot check",
			content: `{"allow-origins": ["https://exampleXcom"]}`,
			wantErr: patternMismatch,
		},
		{
			name:    "wildcard origin rejected by schema pattern",
			content: `{"allow-origins": ["*"]}`,
			wantErr: patternMismatch,
		},
		{
			name:    "origin with port rejected by schema pattern",
			content: `{"allow-origins": ["https://example.com:8443"]}`,
			wantErr: patternMismatch,
		},
		{
			name:    "allow-methods string instead of array",
			content: `{"allow-methods": "GET"}`,
			wantErr: "Invalid type. Expected: array, given: string",
		},
		{
			name:    "request-methods string instead of array",
			content: `{"request-methods": "GET"}`,
			wantErr: "Invalid type. Expected: array, given: string",
		},
		{
			name:    "allow-methods non-string element type",
			content: `{"allow-methods": [123]}`,
			wantErr: "Invalid type. Expected: string",
		},
		{
			name:    "request-methods non-string element type",
			content: `{"request-methods": [true]}`,
			wantErr: "Invalid type. Expected: string",
		},
	}
	runSchemaTestCases(t, cases)
}
