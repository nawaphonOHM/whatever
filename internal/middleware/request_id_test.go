package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRequestID_GeneratesNewID(t *testing.T) {
	r := gin.New()
	r.Use(RequestID())

	var capturedCtxID string
	r.GET("/test", func(c *gin.Context) {
		capturedCtxID = GetRequestID(c)
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	respHeaderID := w.Header().Get(HeaderXRequestID)
	require.NotEmpty(t, respHeaderID)
	assert.Equal(t, respHeaderID, capturedCtxID)

	// Validate it is a valid UUID
	_, err := uuid.Parse(respHeaderID)
	assert.NoError(t, err)
}

func TestRequestID_PropagatesExistingID(t *testing.T) {
	r := gin.New()
	r.Use(RequestID())

	customID := "custom-request-id-12345"
	var capturedCtxID string
	r.GET("/test", func(c *gin.Context) {
		capturedCtxID = GetRequestID(c)
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(HeaderXRequestID, customID)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, customID, w.Header().Get(HeaderXRequestID))
	assert.Equal(t, customID, capturedCtxID)
}

func TestGetRequestID_Fallback(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(HeaderXRequestID, "fallback-id")
	c.Request = req

	assert.Equal(t, "fallback-id", GetRequestID(c))
}
