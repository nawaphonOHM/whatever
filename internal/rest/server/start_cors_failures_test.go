package server

import (
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type serverFailureCase struct {
	name          string
	fixture       string
	expectedError string
}

func getStartupFailureCases() []serverFailureCase {
	return []serverFailureCase{
		{"missing file", "nonexistent.json", "failed to read CORS configuration file"},
		{"unsupported extension", "invalid.txt", "unsupported CORS configuration file format"},
		{"invalid json syntax", "invalid_syntax.json", "failed to parse JSON CORS configuration"},
		{"invalid yaml syntax", "invalid_syntax.yaml", "failed to parse YAML CORS configuration"},
		{"multi-doc yaml", "invalid_multidoc.yaml", "multiple documents or trailing content detected"},
		{"disallowed method", "invalid_schema_disallowed_method.json", schemaValidationErr},
		{"bad origin format", "invalid_schema_bad_origin.json", schemaValidationErr},
		{"bad header format", "invalid_schema_bad_header.json", schemaValidationErr},
		{"unknown property", "invalid_unknown_property.json", schemaValidationErr},
		{"yaml disallowed method", "invalid_schema_disallowed_method.yaml", schemaValidationErr},
		{"yaml bad origin format", "invalid_schema_bad_origin.yaml", schemaValidationErr},
	}
}

// TestNewFromBluePrint_StartupFailures verifies startup errors for invalid files.
func TestNewFromBluePrint_StartupFailures(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envCORSPathKey, "")

	for _, tc := range getStartupFailureCases() {
		t.Run(tc.name, func(t *testing.T) {
			corsPath := filepath.Join("testdata", tc.fixture)
			srv, err := NewFromBluePrint(createTestBlueprintWithRoute(corsPath))
			require.Error(t, err)
			assert.Nil(t, srv)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}
