package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProblemResponse(t *testing.T) {
	resp := Error(
		http.StatusTeapot, "IM_A_TEAPOT", "Short and stout",
	)
	assert.Equal(t, http.StatusTeapot, resp.StatusCode())
}
