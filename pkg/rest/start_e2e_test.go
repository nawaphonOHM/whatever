package rest

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestE2EBlueprint creates a blueprint with sample versioned and unversioned routes.
func newTestE2EBlueprint() *BluePrint {
	regPing := &RRestAPIRegistration{
		Version: 1, Prefix: "/api",
		Apis: []*ExportableAPI{{
			Path: "/ping", Method: GET,
			Handler: func(Context) Response { return OK("pong") },
		}},
	}
	regItems := &RRestAPIRegistration{
		Version: 1, Prefix: "/items",
		Apis: []*ExportableAPI{{
			Path: "/:id", Method: GET,
			Handler: func(c Context) Response { return OK(c.Param("id")) },
		}},
	}
	regStatus := &RRestAPIRegistration{
		Version: 0, Prefix: "/status",
		Apis: []*ExportableAPI{{
			Path: "", Method: GET,
			Handler: func(Context) Response { return OK("healthy") },
		}},
	}
	return NewBluePrint().WithAPIs(regPing, regItems, regStatus)
}

// TestStartREST_Success boots, serves a registered route, then exits.
func TestStartREST_Success(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(newTestE2EBlueprint())

	assertStatusOK(t, fmt.Sprintf("http://%s:%d/api/v1/ping", testHost, port))
	assertStatusOK(t, fmt.Sprintf("http://%s:%d/api/v1/items/42", testHost, port))
	assertStatusOK(t, fmt.Sprintf("http://%s:%d/api/status", testHost, port))

	signalAndWait(t, errCh)
}

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
