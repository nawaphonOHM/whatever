package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// testSuccessResponse mirrors the public success response JSON shape.
type testSuccessResponse[T any] struct {
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}

// testResponse is a minimal response fixture for server tests.
type testResponse struct {
	headers map[string]string
	body    any
	status  int
}

// StatusCode returns the fixture HTTP status code.
func (r *testResponse) StatusCode() int {
	return r.status
}

// WithHeader attaches a custom key-value header to the response.
func (r *testResponse) WithHeader(key, value string) contracts.Response {
	if r == nil {
		return r
	}
	contracts.SetHeader(&r.headers, key, value)
	return r
}

// WithHeaderInt attaches a custom header whose integer value is converted to string.
func (r *testResponse) WithHeaderInt(key string, value int) contracts.Response {
	if r == nil {
		return r
	}
	return r.WithHeader(key, strconv.Itoa(value))
}

// WithHeaderStringer attaches a custom header whose value is obtained from fmt.Stringer.
func (r *testResponse) WithHeaderStringer(key string, value fmt.Stringer) contracts.Response {
	if r == nil {
		return r
	}
	return applyTestHeaderStringer(r, key, value)
}

func applyTestHeaderStringer(r *testResponse, key string, value fmt.Stringer) contracts.Response {
	str, ok := contracts.SafeStringer(value)
	if !ok {
		return r
	}
	return r.WithHeader(key, str)
}

// Write writes the fixture response to the Gin context.
func (r *testResponse) Write(c *gin.Context) {
	if !canWriteTestResponse(c, r) {
		return
	}
	contracts.WriteHeaders(c, r.headers)
	writeTestResponseBody(c, r.status, r.body)
}

func canWriteTestResponse(c *gin.Context, r *testResponse) bool {
	if c == nil {
		return false
	}
	return r != nil
}

func writeTestResponseBody(c *gin.Context, status int, body any) {
	if body == nil {
		c.Status(status)
		return
	}
	c.JSON(status, body)
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
