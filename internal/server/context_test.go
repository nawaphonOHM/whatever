package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/example/go-boilerplate/pkg/rest/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestContext_ParamAndQuery(t *testing.T) {
	r := gin.New()
	var capturedParam string
	var capturedQuery string
	var capturedDefault string

	r.GET("/users/:id", func(gc *gin.Context) {
		c := NewContext(gc)
		capturedParam = c.Param("id")
		capturedQuery = c.Query("filter")
		capturedDefault = c.DefaultQuery("sort", "asc")
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42?filter=active", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "42", capturedParam)
	assert.Equal(t, "active", capturedQuery)
	assert.Equal(t, "asc", capturedDefault)
}

func TestContext_HeadersAndCookies(t *testing.T) {
	r := gin.New()
	var headerVal string

	r.GET("/test-headers", func(gc *gin.Context) {
		c := NewContext(gc)
		headerVal = c.GetHeader("X-Custom-Header")
		c.SetHeader("X-Response-Header", "processed")
		c.SetCookie("session_id", "abc123xyz", 3600, "/", "example.com", true, true)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test-headers", nil)
	req.Header.Set("X-Custom-Header", "test-val")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test-val", headerVal)
	assert.Equal(t, "processed", w.Header().Get("X-Response-Header"))
	assert.Contains(t, w.Header().Get("Set-Cookie"), "session_id=abc123xyz")
}

func TestContext_Binding(t *testing.T) {
	type BodyPayload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	r := gin.New()
	var boundPayload BodyPayload

	r.POST("/bind-test", func(gc *gin.Context) {
		c := NewContext(gc)
		if err := c.ShouldBindJSON(&boundPayload); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodPost, "/bind-test", strings.NewReader(`{"name":"Alice","age":30}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "Alice", boundPayload.Name)
	assert.Equal(t, 30, boundPayload.Age)
}

func TestContext_StateStore(t *testing.T) {
	r := gin.New()

	r.GET("/state-test", func(gc *gin.Context) {
		c := NewContext(gc)
		c.Set("user_id", "user-123")
		c.Set("is_admin", true)
		c.Set("count", 42)
		c.Set("duration", 5*time.Minute)

		assert.Equal(t, "user-123", c.GetString("user_id"))
		assert.True(t, c.GetBool("is_admin"))
		assert.Equal(t, 42, c.GetInt("count"))
		assert.Equal(t, 5*time.Minute, c.GetDuration("duration"))
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/state-test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestContext_AbortWithResponse(t *testing.T) {
	r := gin.New()

	r.GET("/abort-test", func(gc *gin.Context) {
		c := NewContext(gc)
		c.AbortWithResponse(response.BadRequest("INVALID_REQUEST", "Bad parameter"))
	}, func(gc *gin.Context) {
		// Should not be called
		t.Fatal("subsequent handler was called after AbortWithResponse")
	})

	req := httptest.NewRequest(http.MethodGet, "/abort-test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
}

func TestContext_Metadata(t *testing.T) {
	r := gin.New()

	r.POST("/meta-test", func(gc *gin.Context) {
		c := NewContext(gc)
		assert.NotEmpty(t, c.ClientIP())
		assert.Equal(t, "application/json", c.ContentType())
		assert.Equal(t, "/meta-test", c.FullPath())
		assert.NotNil(t, c.Context())
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/meta-test", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
