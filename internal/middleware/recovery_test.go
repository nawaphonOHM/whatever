package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/go-boilerplate/pkg/rest/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecovery_HandlesPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	r := gin.New()
	r.Use(RequestID())
	r.Use(RecoveryWithLogger(logger))

	r.GET("/panic", func(c *gin.Context) {
		panic("something went critically wrong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/panic", nil)
	req.Header.Set(HeaderXRequestID, "panic-req-id")

	// Ensure the test does not crash from panic
	assert.NotPanics(t, func() {
		r.ServeHTTP(w, req)
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))

	var prob response.ProblemDetails
	err := json.Unmarshal(w.Body.Bytes(), &prob)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, prob.Status)
	assert.Equal(t, "about:blank", prob.Type)
	assert.Equal(t, "Internal Server Error", prob.Title)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", prob.Code)
	assert.Equal(t, "An unexpected internal server error occurred", prob.Detail)
	assert.Equal(t, "/panic", prob.Instance)

	// Verify structured log captured the panic
	logOutput := buf.String()
	assert.Contains(t, logOutput, "panic recovered during request processing")
	assert.Contains(t, logOutput, "something went critically wrong")
	assert.Contains(t, logOutput, "panic-req-id")
}
