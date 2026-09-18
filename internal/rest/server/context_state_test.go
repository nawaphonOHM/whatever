package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	stateUserID = "user-123"
	stateCount  = 42
	stateMins   = 5
)

// TestContext_StateStore verifies typed getters for stored values.
func TestContext_StateStore(t *testing.T) {
	// Arrange a route that stores and reads context values.
	r := gin.New()
	r.GET("/state-test", func(gc *gin.Context) {
		c := NewContext(gc)
		c.Set("user_id", stateUserID)
		c.Set("is_admin", true)
		c.Set("count", stateCount)
		c.Set("duration", stateMins*time.Minute)

		assert.Equal(t, stateUserID, c.GetString("user_id"))
		assert.True(t, c.GetBool("is_admin"))
		assert.Equal(t, stateCount, c.GetInt("count"))
		assert.Equal(t, stateMins*time.Minute, c.GetDuration("duration"))
		c.Status(http.StatusOK)
	})

	// Act
	req := httptest.NewRequest(http.MethodGet, "/state-test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}
