package server

import "testing"

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
