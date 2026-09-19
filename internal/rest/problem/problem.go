// Package problem provides RFC 9457 Problem Details object encapsulation.
package problem

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MediaTypeProblemJSON is canonical media type for RFC 9457 Problem Details.
const MediaTypeProblemJSON = "application/problem+json"

// DefaultProblemType is default URI reference identifying problem type.
const DefaultProblemType = "about:blank"

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

// Write serializes the ProblemDetails object to the Gin context.
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
		p.Type = DefaultProblemType
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

// New creates an RFC 9457 Problem Details pointer object.
func New(
	statusCode int,
	code, detail string,
	details ...any,
) *ProblemDetails {
	title := http.StatusText(statusCode)
	if title == "" {
		title = "Error"
	}
	prob := &ProblemDetails{
		Type:   DefaultProblemType,
		Title:  title,
		Status: statusCode,
		Detail: detail,
		Code:   code,
	}
	prob.setDetails(details...)
	return prob
}
