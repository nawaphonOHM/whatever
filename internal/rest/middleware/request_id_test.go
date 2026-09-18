// Package middleware provides unit tests for request ID propagation middleware.
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

const (
	testRouteReqID = "/test"
	fallbackTestID = "fallback-id"
)

// init initializes the gin test mode.
func init() {
	gin.SetMode(gin.TestMode)
}

// setupTestIDRouter creates router with request ID middleware.
func setupTestIDRouter(capturedCtxID *string) *gin.Engine {
	r := gin.New()
	r.Use(RequestID())
	r.GET(testRouteReqID, func(c *gin.Context) {
		*capturedCtxID = GetRequestID(c)
		c.String(http.StatusOK, "ok")
	})
	return r
}

// executeRequestIDRoute sets up a test engine and returns recorded response
// and captured ID.
func executeRequestIDRoute(
	t *testing.T,
	headerVal string,
) (*httptest.ResponseRecorder, string) {
	var capturedCtxID string
	r := setupTestIDRouter(&capturedCtxID)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, testRouteReqID, nil)
	require.NoError(t, err)
	if headerVal != "" {
		req.Header.Set(HeaderXRequestID, headerVal)
	}
	r.ServeHTTP(w, req)
	return w, capturedCtxID
}

// TestRequestID_GeneratesNewID tests generating a fresh UUID.
func TestRequestID_GeneratesNewID(t *testing.T) {
	w, capturedCtxID := executeRequestIDRoute(t, "")

	assert.Equal(t, http.StatusOK, w.Code)
	respHeaderID := w.Header().Get(HeaderXRequestID)
	require.NotEmpty(t, respHeaderID)
	assert.Equal(t, respHeaderID, capturedCtxID)

	// Validate it is a valid UUID
	_, err := uuid.Parse(respHeaderID)
	assert.NoError(t, err)
}

// TestRequestID_PropagatesExistingID tests propagating custom ID header.
func TestRequestID_PropagatesExistingID(t *testing.T) {
	customID := "custom-request-id-12345"
	w, capturedCtxID := executeRequestIDRoute(t, customID)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, customID, w.Header().Get(HeaderXRequestID))
	assert.Equal(t, customID, capturedCtxID)
}

// TestGetRequestID_Fallback tests reading ID from header if context key absent.
func TestGetRequestID_Fallback(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, err := http.NewRequest(http.MethodGet, testRouteReqID, nil)
	require.NoError(t, err)
	req.Header.Set(HeaderXRequestID, fallbackTestID)
	c.Request = req

	assert.Equal(t, fallbackTestID, GetRequestID(c))
}
