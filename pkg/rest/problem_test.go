package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	intprob "github.com/nawaphonOHM/whatever/internal/rest/problem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertProblemResponse verifies the unmarshaled problem details object.
func assertProblemResponse(
	t *testing.T,
	prob intprob.ProblemDetails,
	tc problemTestCase,
) {
	assert.Equal(t, tc.expectedStatus, prob.Status)
	assert.Equal(t, intprob.DefaultProblemType, prob.Type)
	assert.Equal(t, tc.expectedTitle, prob.Title)
	assert.Equal(t, tc.expectedCode, prob.Code)
	assert.Equal(t, tc.expectedDetail, prob.Detail)
}

// performProblemRequest executes a GET request against an endpoint with
// the given handler.
func performProblemRequest(
	t *testing.T,
	endpoint string,
	handler gin.HandlerFunc,
) *httptest.ResponseRecorder {
	r := setupTestRouter()
	r.GET(endpoint, handler)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	require.NoError(t, err)
	r.ServeHTTP(w, req)
	return w
}

// runProblemTest executes a single problem test case against a test router.
func runProblemTest(t *testing.T, tc problemTestCase) {
	w := performProblemRequest(t, tc.endpoint, tc.handler)
	assert.Equal(t, tc.expectedStatus, w.Code)
	assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
	var prob intprob.ProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &prob))
	assertProblemResponse(t, prob, tc)
}

// TestErrorResponses tests helper functions generating problem responses.
func TestErrorResponses(t *testing.T) {
	for _, tc := range buildProblemCases() {
		t.Run(tc.name, func(t *testing.T) {
			runProblemTest(t, tc)
		})
	}
}
