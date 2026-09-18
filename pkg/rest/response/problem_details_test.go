package response

import (
	"net/http"
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
