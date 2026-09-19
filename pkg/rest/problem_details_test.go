package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	intprob "github.com/nawaphonOHM/whatever/internal/rest/problem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func verifyMultipleDetailsResponse(
	t *testing.T,
	body []byte,
) {
	t.Helper()
	var prob intprob.ProblemDetails
	require.NoError(t, json.Unmarshal(body, &prob))
	assert.Equal(t, http.StatusBadRequest, prob.Status)
	assert.Equal(t, intprob.DefaultProblemType, prob.Type)
	assert.Equal(t, "Bad Request", prob.Title)
	assert.Equal(t, "VALIDATION_ERROR", prob.Code)
	assert.Equal(t, "Multiple issues", prob.Detail)
	assert.Equal(t, []any{"issue 1", "issue 2"}, prob.Details)
}

func TestError_MultipleDetails(t *testing.T) {
	resp := Error(
		http.StatusBadRequest,
		"VALIDATION_ERROR",
		"Multiple issues",
		"issue 1",
		"issue 2",
	)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode())

	gc, w := newTestContext()
	resp.Write(gc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
	verifyMultipleDetailsResponse(t, w.Body.Bytes())
}

func verifyWriteDefaultsFields(t *testing.T, body []byte) {
	t.Helper()
	var prob intprob.ProblemDetails
	require.NoError(t, json.Unmarshal(body, &prob))
	assert.Equal(t, http.StatusNotFound, prob.Status)
	assert.Equal(t, intprob.DefaultProblemType, prob.Type)
	assert.Equal(t, "Not Found", prob.Title)
	assert.Equal(t, "/test/path", prob.Instance)
	assert.Equal(t, "Resource was not located", prob.Detail)
}

func TestProblemDetails_WriteDefaults(t *testing.T) {
	gc, w := newTestContext()
	gc.Request = httptest.NewRequest(http.MethodGet, "/test/path", nil)

	resp := NotFound("NOT_FOUND", "Resource was not located")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode())

	resp.Write(gc)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
	verifyWriteDefaultsFields(t, w.Body.Bytes())
}

func TestProblemDetails_UnknownStatusTitle(t *testing.T) {
	gc, w := newTestContext()
	probCustom := Error(
		testUnknownStatus, "CUSTOM_ERR", "custom",
	)
	assert.Equal(t, testUnknownStatus, probCustom.StatusCode())

	probCustom.Write(gc)
	var prob intprob.ProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &prob))
	assert.Equal(t, "Error", prob.Title)
}

func TestProblemDetails_NilContext(t *testing.T) {
	resp := BadRequest("ERR", "detail")
	assert.NotPanics(t, func() {
		resp.Write(nil)
	})
}
