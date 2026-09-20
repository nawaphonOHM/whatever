package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

// testSuccessResponse mirrors the public success response JSON shape.
type testSuccessResponse[T any] struct {
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}

// testResponse is a minimal response fixture for server tests.
type testResponse struct {
	body   any
	status int
}

// StatusCode returns the fixture HTTP status code.
func (r *testResponse) StatusCode() int {
	return r.status
}

// Write writes the fixture response to the Gin context.
func (r *testResponse) Write(c *gin.Context) {
	if r.body == nil {
		c.Status(r.status)
		return
	}
	c.JSON(r.status, r.body)
}

func testNoContent() contracts.Response {
	return &testResponse{status: http.StatusNoContent}
}

func testOK(data any) contracts.Response {
	return &testResponse{
		status: http.StatusOK,
		body: testSuccessResponse[any]{
			Data:      data,
			Timestamp: time.Now().UTC(),
			Success:   true,
		},
	}
}
