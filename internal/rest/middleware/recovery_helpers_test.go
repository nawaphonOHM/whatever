package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/problem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	panicMsg        = "something went critically wrong"
	panicPath       = "/panic"
	panicReqID      = "panic-req-id"
	expectedLogMsg  = "panic recovered during request processing"
	expectedProbErr = "INTERNAL_SERVER_ERROR"
)

func setupPanicRouter(logger *slog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(RequestID())
	r.Use(RecoveryWithLogger(logger))
	r.GET(panicPath, func(*gin.Context) {
		panic(panicMsg)
	})
	return r
}

func setupDefaultRecoveryRouter() *gin.Engine {
	r := gin.New()
	r.Use(RequestID())
	r.Use(Recovery())
	r.GET(panicPath, func(*gin.Context) {
		panic(panicMsg)
	})
	return r
}

func executePanicRequest(t *testing.T, r *gin.Engine) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, panicPath, nil)
	require.NoError(t, err)
	req.Header.Set(HeaderXRequestID, panicReqID)
	r.ServeHTTP(w, req)
	return w
}

func assertProblemFields(t *testing.T, prob problem.PProblemDetails) {
	t.Helper()
	assert.Equal(t, http.StatusInternalServerError, prob.Status)
	assert.Equal(t, "about:blank", prob.Type)
	assert.Equal(t, "Internal Server Error", prob.Title)
	assert.Equal(t, expectedProbErr, prob.Code)
	assert.Equal(t, panicPath, prob.Instance)
}

func verifyPanicResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))

	var prob problem.PProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &prob))
	assertProblemFields(t, prob)
}

func parseJSONLogEntry(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var entry map[string]any
	err := json.Unmarshal(bytes.TrimSpace(data), &entry)
	require.NoError(t, err)
	return entry
}
