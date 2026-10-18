// Package problem provides RFC 9457 Problem Details object encapsulation.
package problem

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// MediaTypeProblemJSON is canonical media type for RFC 9457 Problem Details.
const MediaTypeProblemJSON = contracts.MediaTypeProblemJSON

// DefaultProblemType is default URI reference identifying problem type.
const DefaultProblemType = "about:blank"

// PProblemDetails represents an RFC 9457 Problem Details object.
type PProblemDetails struct {
	headers  map[string]string
	Details  any    `json:"details,omitempty"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Code     string `json:"code,omitempty"`
	Status   int    `json:"status"`
}

// StatusCode returns the HTTP status code for the PProblemDetails.
func (p *PProblemDetails) StatusCode() int {
	if p == nil || p.Status == 0 {
		return http.StatusInternalServerError
	}
	return p.Status
}

// Write serializes the PProblemDetails object to the Gin context.
func (p *PProblemDetails) Write(c *gin.Context) {
	if c == nil || p == nil {
		return
	}
	p.ensureDefaults(c)
	contracts.WriteHeaders(c, p.headers)
	c.Header("Content-Type", MediaTypeProblemJSON)
	c.JSON(p.StatusCode(), p)
}

// WithHeader attaches a custom key-value header to the response.
func (p *PProblemDetails) WithHeader(key, value string) contracts.Response {
	if p == nil {
		return p
	}
	contracts.SetHeader(&p.headers, key, value)
	return p
}

// WithHeaderInt attaches a custom header whose integer value is converted to string.
func (p *PProblemDetails) WithHeaderInt(key string, value int) contracts.Response {
	if p == nil {
		return p
	}
	return p.WithHeader(key, strconv.Itoa(value))
}

// WithHeaderStringer attaches a custom header whose value is obtained from fmt.Stringer.
func (p *PProblemDetails) WithHeaderStringer(key string, value fmt.Stringer) contracts.Response {
	if p == nil {
		return p
	}
	return applyProblemHeaderStringer(p, key, value)
}

func applyProblemHeaderStringer(p *PProblemDetails, key string, value fmt.Stringer) contracts.Response {
	str, ok := contracts.SafeStringer(value)
	if !ok {
		return p
	}
	return p.WithHeader(key, str)
}

// New creates an RFC 9457 Problem Details pointer object.
func New(
	statusCode int,
	code, detail string,
	details ...any,
) *PProblemDetails {
	title := http.StatusText(statusCode)
	if title == "" {
		title = "Error"
	}
	prob := &PProblemDetails{
		Type:   DefaultProblemType,
		Title:  title,
		Status: statusCode,
		Detail: detail,
		Code:   code,
	}
	prob.setDetails(details...)
	return prob
}
