package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// SuccessResponse represents the standard JSON API response envelope DTO
// for successful responses, used for unmarshaling in clients and tests.
type SuccessResponse[T any] struct {
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message,omitempty"`
	Success   bool      `json:"success"`
}

// Envelope is a type alias for SuccessResponse.
type Envelope[T any] = SuccessResponse[T]

// successResponse represents the internal response renderer implementing
// rest.Response.
type successResponse[T any] struct {
	headers   map[string]string
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message,omitempty"`
	Code      int       `json:"-"`
	Success   bool      `json:"success"`
}

// StatusCode returns the HTTP status code for the response.
func (r *successResponse[T]) StatusCode() int {
	if r == nil || r.Code == 0 {
		return http.StatusOK
	}
	return r.Code
}

// Write writes the JSON response envelope to the Gin context.
func (r *successResponse[T]) Write(c *gin.Context) {
	if c != nil && r != nil {
		contracts.WriteHeaders(c, r.headers)
		c.JSON(r.StatusCode(), r)
	}
}

// WithHeader attaches a custom key-value header to the response.
func (r *successResponse[T]) WithHeader(key, value string) contracts.Response {
	if r == nil {
		return r
	}
	contracts.SetHeader(&r.headers, key, value)
	return r
}

// WithHeaderInt attaches a custom header whose integer value is converted to string.
func (r *successResponse[T]) WithHeaderInt(key string, value int) contracts.Response {
	if r == nil {
		return r
	}
	return r.WithHeader(key, strconv.Itoa(value))
}

// WithHeaderStringer attaches a custom header whose value is obtained from fmt.Stringer.
func (r *successResponse[T]) WithHeaderStringer(key string, value fmt.Stringer) contracts.Response {
	if r == nil {
		return r
	}
	return applySuccessHeaderStringer(r, key, value)
}

func applySuccessHeaderStringer(r contracts.Response, key string, value fmt.Stringer) contracts.Response {
	str, ok := contracts.SafeStringer(value)
	if !ok {
		return r
	}
	return r.WithHeader(key, str)
}

// withMessage sets the message on the successResponse and returns it for
// chaining.
func (r *successResponse[T]) withMessage(msg string) *successResponse[T] {
	if r != nil {
		r.Message = msg
	}
	return r
}

// withTimestamp sets the timestamp on the successResponse and returns it for
// chaining.
func (r *successResponse[T]) withTimestamp(t time.Time) *successResponse[T] {
	if r != nil {
		r.Timestamp = t
	}
	return r
}

// newSuccessResponse creates a successful successResponse envelope pointer.
func newSuccessResponse[T any](data T, message ...string) *successResponse[T] {
	resp := &successResponse[T]{
		Code:      http.StatusOK,
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	if len(message) > 0 {
		resp.Message = message[0]
	}
	return resp
}
