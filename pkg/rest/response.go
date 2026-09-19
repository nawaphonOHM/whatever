package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MediaTypeProblemJSON is the canonical media type for RFC 9457 Problem
// Details.
const MediaTypeProblemJSON = "application/problem+json"

// Response represents any HTTP response that can write itself to a Gin context.
type Response interface {
	// StatusCode returns the HTTP status code of the response.
	StatusCode() int
	// Write renders the response to the given Gin context.
	Write(*gin.Context)
}

// noContentResponse represents an empty HTTP 204 No Content response.
type noContentResponse struct{}

// StatusCode returns 204 No Content.
func (*noContentResponse) StatusCode() int {
	return http.StatusNoContent
}

// Write writes the 204 status header to the Gin context.
func (*noContentResponse) Write(c *gin.Context) {
	if c != nil {
		c.Status(http.StatusNoContent)
	}
}

// JSON creates a Response pointer with the given status code and data.
func JSON(statusCode int, data any, message ...string) Response {
	resp := newSuccessResponse(data, message...)
	resp.Code = statusCode
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
