// Package response provides standardized JSON API response envelopes and RFC 9457 Problem Details.
package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// MediaTypeProblemJSON is the canonical media type for RFC 9457 Problem Details.
const MediaTypeProblemJSON = "application/problem+json"

// Response represents any HTTP response that can write itself to a Gin context.
type Response interface {
	StatusCode() int
	Write(c *gin.Context)
}

// SuccessResponse represents the standard JSON API response envelope for successful responses.
type SuccessResponse[T any] struct {
	Code      int       `json:"-"`
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
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

// WithMessage sets the message on the SuccessResponse and returns it for chaining.
func (r *SuccessResponse[T]) WithMessage(msg string) *SuccessResponse[T] {
	if r != nil {
		r.Message = msg
	}
	return r
}

// WithTimestamp sets the timestamp on the SuccessResponse and returns it for chaining.
func (r *SuccessResponse[T]) WithTimestamp(t time.Time) *SuccessResponse[T] {
	if r != nil {
		r.Timestamp = t
	}
	return r
}

// Envelope is a type alias for SuccessResponse.
type Envelope[T any] = SuccessResponse[T]

// ProblemDetails represents an RFC 9457 Problem Details object.
type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Code     string `json:"code,omitempty"`
	Details  any    `json:"details,omitempty"`
}

// StatusCode returns the HTTP status code for the ProblemDetails.
func (p *ProblemDetails) StatusCode() int {
	if p == nil || p.Status == 0 {
		return http.StatusInternalServerError
	}
	return p.Status
}

// Write serializes the ProblemDetails object to the Gin context using the application/problem+json media type.
func (p *ProblemDetails) Write(c *gin.Context) {
	if c == nil || p == nil {
		return
	}
	if p.Type == "" {
		p.Type = "about:blank"
	}
	status := p.StatusCode()
	if p.Title == "" {
		p.Title = http.StatusText(status)
		if p.Title == "" {
			p.Title = "Error"
		}
	}
	if p.Instance == "" && c.Request != nil && c.Request.URL != nil {
		p.Instance = c.Request.URL.Path
	}

	c.Header("Content-Type", MediaTypeProblemJSON)
	c.JSON(status, p)
}

// noContentResponse represents an empty HTTP 204 No Content response.
type noContentResponse struct{}

// StatusCode returns 204 No Content.
func (n *noContentResponse) StatusCode() int {
	return http.StatusNoContent
}

// Write writes the 204 status header to the Gin context.
func (n *noContentResponse) Write(c *gin.Context) {
	if c != nil {
		c.Status(http.StatusNoContent)
	}
}

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

// NewProblemDetails creates an RFC 9457 Problem Details pointer object.
func NewProblemDetails(statusCode int, code string, detail string, details ...interface{}) *ProblemDetails {
	title := http.StatusText(statusCode)
	if title == "" {
		title = "Error"
	}
	prob := &ProblemDetails{
		Type:   "about:blank",
		Title:  title,
		Status: statusCode,
		Detail: detail,
		Code:   code,
	}
	if len(details) == 1 {
		prob.Details = details[0]
	} else if len(details) > 1 {
		prob.Details = details
	}
	return prob
}

// JSON creates a Response pointer with the given status code and data.
func JSON(statusCode int, data any, message ...string) Response {
	resp := &SuccessResponse[any]{
		Code:      statusCode,
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	if len(message) > 0 {
		resp.Message = message[0]
	}
	return resp
}

// OK creates a 200 OK Response with the standard success envelope.
func OK(data any, message ...string) Response {
	return JSON(http.StatusOK, data, message...)
}

// Created creates a 201 Created Response with the standard success envelope.
func Created(data any, message ...string) Response {
	return JSON(http.StatusCreated, data, message...)
}

// NoContent creates a 204 No Content Response.
func NoContent() Response {
	return &noContentResponse{}
}

// Problem returns the ProblemDetails as a Response.
func Problem(prob *ProblemDetails) Response {
	return prob
}

// Error creates an RFC 9457 Problem Details Response.
func Error(statusCode int, code string, detail string, details ...any) Response {
	return NewProblemDetails(statusCode, code, detail, details...)
}

// BadRequest creates a 400 Bad Request RFC 9457 Problem Details Response.
func BadRequest(code string, detail string, details ...any) Response {
	return Error(http.StatusBadRequest, code, detail, details...)
}

// Unauthorized creates a 401 Unauthorized RFC 9457 Problem Details Response.
func Unauthorized(code string, detail string, details ...any) Response {
	return Error(http.StatusUnauthorized, code, detail, details...)
}

// Forbidden creates a 403 Forbidden RFC 9457 Problem Details Response.
func Forbidden(code string, detail string, details ...any) Response {
	return Error(http.StatusForbidden, code, detail, details...)
}

// NotFound creates a 404 Not Found RFC 9457 Problem Details Response.
func NotFound(code string, detail string, details ...any) Response {
	return Error(http.StatusNotFound, code, detail, details...)
}

// InternalServerError creates a 500 Internal Server Error RFC 9457 Problem Details Response.
func InternalServerError(code string, detail string, details ...any) Response {
	return Error(http.StatusInternalServerError, code, detail, details...)
}
