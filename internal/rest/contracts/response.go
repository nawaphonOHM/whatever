package contracts

import "github.com/gin-gonic/gin"

// MediaTypeProblemJSON is the canonical media type for RFC 9457 Problem
// Details.
const MediaTypeProblemJSON = "application/problem+json"

// Response represents any HTTP response that can write itself to a Gin
// context.
type Response interface {
	// StatusCode returns the HTTP status code of the response.
	StatusCode() int
	// Write renders the response to the given Gin context.
	Write(*gin.Context)
}
