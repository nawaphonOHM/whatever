package rest

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type e2eFailureCase struct {
	name          string
	fixture       string
	expectedError string
}

func getE2EFailureCases() []e2eFailureCase {
	return []e2eFailureCase{
		{"missing file", "nonexistent.json", "failed to read CORS configuration file"},
		{"invalid syntax", "invalid_syntax.json", "failed to parse JSON CORS configuration"},
		{"schema violation", "invalid_schema_disallowed_method.json", schemaValidationErr},
	}
}

func TestStartREST_StartupFailures(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envCORSKey, "")

	for _, tc := range getE2EFailureCases() {
		t.Run(tc.name, func(t *testing.T) {
			corsPath := filepath.Join("testdata", tc.fixture)
			bp := newCORSBlueprintWithFixture(corsPath)
			errCh := startRESTAsync(bp)
			err := waitErr(t, errCh)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}
