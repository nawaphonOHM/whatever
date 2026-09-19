package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProblemResponse(t *testing.T) {
	prob := NewProblemDetails(
		http.StatusTeapot, "IM_A_TEAPOT", "Short and stout",
	)
	resp := Problem(prob)
	assert.Equal(t, http.StatusTeapot, resp.StatusCode())
}
