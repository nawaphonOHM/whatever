package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
	"github.com/stretchr/testify/assert"
)

const (
	abortCode   = "INVALID_REQUEST"
	abortDetail = "Bad parameter"
	problemCT   = "application/problem+json"
)

// TestContext_AbortWithResponse ensures the chain stops on abort.
func TestContext_AbortWithResponse(t *testing.T) {
	// Arrange handlers where the second must not run.
	r := gin.New()
	r.GET("/abort-test", func(gc *gin.Context) {
		c := NewContext(gc)
		c.AbortWithResponse(
			response.BadRequest(abortCode, abortDetail),
		)
	}, func(*gin.Context) {
		// Should not be called after abort.
		t.Fatal("subsequent handler was called after abort")
	})

	// Act
	req := httptest.NewRequest(http.MethodGet, "/abort-test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert problem details response.
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, problemCT, w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), abortCode)
}

// TestContext_Metadata verifies request metadata helpers.
func TestContext_Metadata(t *testing.T) {
	// Arrange
	r := gin.New()
	r.POST("/meta-test", func(gc *gin.Context) {
		c := NewContext(gc)
		assert.NotEmpty(t, c.ClientIP())
		assert.Equal(t, "application/json", c.ContentType())
		assert.Equal(t, "/meta-test", c.FullPath())
		assert.NotNil(t, c.Context())
		c.Status(http.StatusOK)
	})

	// Act
	req := httptest.NewRequest(
		http.MethodPost,
		"/meta-test",
		strings.NewReader(`{}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}
