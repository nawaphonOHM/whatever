package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const defaultProblemType = "about:blank"

// ProblemDetails represents an RFC 9457 Problem Details object.
type ProblemDetails struct {
	Details  any    `json:"details,omitempty"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Code     string `json:"code,omitempty"`
	Status   int    `json:"status"`
}

// StatusCode returns the HTTP status code for the ProblemDetails.
func (p *ProblemDetails) StatusCode() int {
	if p == nil || p.Status == 0 {
		return http.StatusInternalServerError
	}
	return p.Status
}

// Write serializes the ProblemDetails object to the Gin context using the
// application/problem+json media type.
func (p *ProblemDetails) Write(c *gin.Context) {
	if c == nil || p == nil {
		return
	}
	p.ensureDefaults(c)
	c.Header("Content-Type", MediaTypeProblemJSON)
	c.JSON(p.StatusCode(), p)
}

func (p *ProblemDetails) ensureDefaults(c *gin.Context) {
	if p.Type == "" {
		p.Type = defaultProblemType
	}
	p.ensureTitle()
	p.ensureInstance(c)
}

func (p *ProblemDetails) ensureTitle() {
	if p.Title != "" {
		return
	}
	p.Title = http.StatusText(p.StatusCode())
	if p.Title == "" {
		p.Title = "Error"
	}
}

func (p *ProblemDetails) ensureInstance(c *gin.Context) {
	if p.Instance != "" || c.Request == nil {
		return
	}
	p.extractURLPath(c)
}

func (p *ProblemDetails) extractURLPath(c *gin.Context) {
	if c.Request.URL != nil {
		p.Instance = c.Request.URL.Path
	}
}

func (p *ProblemDetails) setDetails(details ...any) {
	if len(details) == 1 {
		p.Details = details[0]
	} else if len(details) > 1 {
		p.Details = details
	}
}

// NewProblemDetails creates an RFC 9457 Problem Details pointer object.
func NewProblemDetails(
	statusCode int,
	code, detail string,
	details ...any,
) *ProblemDetails {
	title := http.StatusText(statusCode)
	if title == "" {
		title = "Error"
	}
	prob := &ProblemDetails{
		Type:   defaultProblemType,
		Title:  title,
		Status: statusCode,
		Detail: detail,
		Code:   code,
	}
	prob.setDetails(details...)
	return prob
}

// Error creates an RFC 9457 Problem Details Response.
func Error(statusCode int, code, detail string, details ...any) Response {
	return NewProblemDetails(statusCode, code, detail, details...)
}

// BadRequest creates a 400 Bad Request RFC 9457 Problem Details Response.
func BadRequest(code, detail string, details ...any) Response {
	return Error(http.StatusBadRequest, code, detail, details...)
}

// Unauthorized creates a 401 Unauthorized RFC 9457 Problem Details Response.
func Unauthorized(code, detail string, details ...any) Response {
	return Error(http.StatusUnauthorized, code, detail, details...)
}

// Forbidden creates a 403 Forbidden RFC 9457 Problem Details Response.
func Forbidden(code, detail string, details ...any) Response {
	return Error(http.StatusForbidden, code, detail, details...)
}

// NotFound creates a 404 Not Found RFC 9457 Problem Details Response.
func NotFound(code, detail string, details ...any) Response {
	return Error(http.StatusNotFound, code, detail, details...)
}

// InternalServerError creates a 500 Internal Server Error RFC 9457 Problem
// Details Response.
func InternalServerError(code, detail string, details ...any) Response {
	return Error(http.StatusInternalServerError, code, detail, details...)
}
