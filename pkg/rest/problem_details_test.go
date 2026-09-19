package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewProblemDetails_MultipleDetails verifies handling of multiple details.
func TestNewProblemDetails_MultipleDetails(t *testing.T) {
	prob := NewProblemDetails(
		http.StatusBadRequest,
		"VALIDATION_ERROR",
		"Multiple issues",
		"issue 1",
		"issue 2",
	)
	assert.Equal(t, http.StatusBadRequest, prob.Status)
	assert.Equal(t, defaultProblemType, prob.Type)
	assert.Equal(t, "Bad Request", prob.Title)
	assert.Equal(t, "VALIDATION_ERROR", prob.Code)
	assert.Equal(t, "Multiple issues", prob.Detail)
	assert.Equal(t, []any{"issue 1", "issue 2"}, prob.Details)
}

func TestProblemDetails_NilAndZeroDefaults(t *testing.T) {
	var nilProb *ProblemDetails
	assert.Equal(t, http.StatusInternalServerError, nilProb.StatusCode())
	nilProb.Write(nil)
	probZero := &ProblemDetails{}
	assert.Equal(t, http.StatusInternalServerError, probZero.StatusCode())
}

func TestProblemDetails_WriteDefaults(t *testing.T) {
	gc, w := newTestContext()
	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	gc.Request = req
	prob := &ProblemDetails{
		Status: http.StatusNotFound,
		Detail: "Resource was not located",
	}
	prob.Write(gc)
	assert.Equal(t, http.StatusNotFound, prob.StatusCode())
	assert.Equal(t, defaultProblemType, prob.Type)
	assert.Equal(t, "Not Found", prob.Title)
	assert.Equal(t, "/test/path", prob.Instance)
	assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
}

func TestProblemDetails_UnknownStatusTitle(t *testing.T) {
	probCustom := NewProblemDetails(
		testUnknownStatus, "CUSTOM_ERR", "custom",
	)
	assert.Equal(t, "Error", probCustom.Title)
}
