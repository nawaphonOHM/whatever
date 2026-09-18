package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	bindJSONBody = `{"name":"Alice","age":30}`
	bindAge      = 30
)

// bodyPayload is a sample JSON binding target.
type bodyPayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TestContext_Binding verifies ShouldBindJSON behavior.
func TestContext_Binding(t *testing.T) {
	// Arrange
	r := gin.New()
	var bound bodyPayload

	r.POST("/bind-test", func(gc *gin.Context) {
		c := NewContext(gc)
		if err := c.ShouldBindJSON(&bound); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusCreated)
	})

	// Act
	req := httptest.NewRequest(
		http.MethodPost,
		"/bind-test",
		strings.NewReader(bindJSONBody),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "Alice", bound.Name)
	assert.Equal(t, bindAge, bound.Age)
}
