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
			content: "allow-methods: PATCH\n",
			wantErr: "allow-methods must be one of the following",
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
			name:    "wrong type for methods",
			content: "request-methods:\n  - GET\n",
			wantErr: "Invalid type. Expected: string, given: array",
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
