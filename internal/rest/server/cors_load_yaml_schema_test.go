package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAndValidateCORS_YAMLSchemaViolations(t *testing.T) {
	cases := []schemaTestCase{
		{
			name:    "invalid method enum",
			content: "allow-methods:\n  - PATCH\n",
			wantErr: "allow-methods.0 must be one of the following",
		},
		{
			name:    "mixed valid and invalid allow-methods at index 1",
			content: "allow-methods:\n  - GET\n  - PATCH\n",
			wantErr: "allow-methods.1 must be one of the following",
		},
		{
			name:    "mixed valid and invalid request-methods at index 1",
			content: "request-methods:\n  - POST\n  - INVALID_METHOD\n",
			wantErr: "request-methods.1 must be one of the following",
		},
		{
			name:    "lowercase allow-methods rejected",
			content: "allow-methods:\n  - get\n",
			wantErr: "allow-methods.0 must be one of the following",
		},
		{
			name:    "lowercase request-methods rejected",
			content: "request-methods:\n  - post\n",
			wantErr: "request-methods.0 must be one of the following",
		},
		{
			name:    "non-string item in allow-methods",
			content: "allow-methods:\n  - 123\n",
			wantErr: "Invalid type. Expected: string",
		},
		{
			name:    "non-string item in request-methods",
			content: "request-methods:\n  - true\n",
			wantErr: "Invalid type. Expected: string",
		},
		{
			name:    "unknown property",
			content: "allow-origins:\n  - https://example.com\nunknown: val\n",
			wantErr: "Additional property unknown is not allowed",
		},
		{
			name:    "bad origin pattern",
			content: "allow-origins:\n  - invalid-origin\n",
			wantErr: patternMismatch,
		},
		{
			name:    "lowercase header",
			content: "allow-headers:\n  - content-type\n",
			wantErr: "must be one of the following",
		},
		{
			name:    "wrong type for allow-methods",
			content: "allow-methods: GET\n",
			wantErr: "Invalid type. Expected: array, given: string",
		},
		{
			name:    "wrong type for request-methods",
			content: "request-methods: GET\n",
			wantErr: "Invalid type. Expected: array, given: string",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTempConfigFile(t, "cors.yaml", tc.content)
			cfg, err := LoadAndValidateCORS(path)
			require.Error(t, err)
			assert.Nil(t, cfg)
			assert.Contains(t, err.Error(), schemaValidationErr)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
