package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	paramRoute   = "/users/:id"
	paramRequest = "/users/42?active=true"
	paramID      = "42"
	queryActive  = "true"
	queryMissing = "no"
)

// paramQueryResult holds captured path/query values.
type paramQueryResult struct {
	param string
	query string
	def   string
}

// captureParamQuery reads param and query values from context.
func captureParamQuery(c *Context) paramQueryResult {
	return paramQueryResult{
		param: c.Param("id"),
		query: c.Query("active"),
		def:   c.DefaultQuery("missing", queryMissing),
	}
}

// TestContext_ParamAndQuery verifies path and query helpers.
func TestContext_ParamAndQuery(t *testing.T) {
	// Arrange
	r := gin.New()
	var got paramQueryResult
	r.GET(paramRoute, func(gc *gin.Context) {
		got = captureParamQuery(NewContext(gc))
	})

	// Act
	req := httptest.NewRequest(http.MethodGet, paramRequest, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, paramID, got.param)
	assert.Equal(t, queryActive, got.query)
	assert.Equal(t, queryMissing, got.def)
}
