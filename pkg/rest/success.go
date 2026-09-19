package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SuccessResponse represents the standard JSON API response envelope for
// successful responses.
type SuccessResponse[T any] struct {
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message,omitempty"`
	Code      int       `json:"-"`
	Success   bool      `json:"success"`
}

// StatusCode returns the HTTP status code for the response.
func (r *SuccessResponse[T]) StatusCode() int {
	if r == nil || r.Code == 0 {
		return http.StatusOK
	}
	return r.Code
}

// Write writes the JSON response envelope to the Gin context.
func (r *SuccessResponse[T]) Write(c *gin.Context) {
	if c != nil && r != nil {
		c.JSON(r.StatusCode(), r)
	}
}

// WithMessage sets the message on the SuccessResponse and returns it for
// chaining.
func (r *SuccessResponse[T]) WithMessage(msg string) *SuccessResponse[T] {
	if r != nil {
		r.Message = msg
	}
	return r
}

// WithTimestamp sets the timestamp on the SuccessResponse and returns it for
// chaining.
func (r *SuccessResponse[T]) WithTimestamp(t time.Time) *SuccessResponse[T] {
	if r != nil {
		r.Timestamp = t
	}
	return r
}

// Envelope is a type alias for SuccessResponse.
type Envelope[T any] = SuccessResponse[T]

// NewSuccessResponse creates a successful SuccessResponse envelope pointer.
func NewSuccessResponse[T any](data T, message ...string) *SuccessResponse[T] {
	resp := &SuccessResponse[T]{
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
