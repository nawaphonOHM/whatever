package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/problem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	missingPath      = "/does-not-exist"
	problemTypeJSON  = "application/problem+json"
	codeNotFoundTest = "NOT_FOUND"
	codeMethodTest   = "METHOD_NOT_ALLOWED"
)

// newMiddlewareEngine builds an engine with default middlewares.
func newMiddlewareEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	srv := New(&Config{Mode: gin.TestMode})
	srv.SetupDefaultMiddlewares()
	srv.Engine.GET("/only-get", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return srv.Engine
}

// decodeProblem unmarshals an RFC 9457 problem body.
func decodeProblem(
	t *testing.T,
	body []byte,
) problem.ProblemDetails {
	t.Helper()
	var prob problem.ProblemDetails
	require.NoError(t, json.Unmarshal(body, &prob))
	return prob
}

// TestSetupDefaultMiddlewares_NoRoute returns RFC 9457 404.
func TestSetupDefaultMiddlewares_NoRoute(t *testing.T) {
	// Arrange
	engine := newMiddlewareEngine()
	w := httptest.NewRecorder()

	// Act
	engine.ServeHTTP(
		w,
		httptest.NewRequest(http.MethodGet, missingPath, nil),
	)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, problemTypeJSON, w.Header().Get("Content-Type"))
	prob := decodeProblem(t, w.Body.Bytes())
	assert.Equal(t, codeNotFoundTest, prob.Code)
}

// TestSetupDefaultMiddlewares_NoMethod returns RFC 9457 405.
func TestSetupDefaultMiddlewares_NoMethod(t *testing.T) {
	// Arrange
	engine := newMiddlewareEngine()
	w := httptest.NewRecorder()

	// Act
	engine.ServeHTTP(
		w,
		httptest.NewRequest(http.MethodPost, "/only-get", nil),
	)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, problemTypeJSON, w.Header().Get("Content-Type"))
	prob := decodeProblem(t, w.Body.Bytes())
	assert.Equal(t, codeMethodTest, prob.Code)
}
