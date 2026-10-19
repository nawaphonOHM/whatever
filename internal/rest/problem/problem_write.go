package problem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

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

func (p *PProblemDetails) ensureDefaults(c *gin.Context) {
	if p.Status == 0 {
		p.Status = p.StatusCode()
	}
	if p.Type == "" {
		p.Type = DefaultProblemType
	}
	p.ensureTitle()
	p.ensureInstance(c)
}

func (p *PProblemDetails) ensureTitle() {
	if p.Title != "" {
		return
	}
	p.Title = http.StatusText(p.StatusCode())
	if p.Title == "" {
		p.Title = "Error"
	}
}

func (p *PProblemDetails) ensureInstance(c *gin.Context) {
	if p.Instance != "" || c.Request == nil {
		return
	}
	p.extractURLPath(c)
}

func (p *PProblemDetails) extractURLPath(c *gin.Context) {
	if c.Request.URL != nil {
		p.Instance = c.Request.URL.Path
	}
}

func (p *PProblemDetails) setDetails(details ...any) {
	if len(details) == 1 {
		p.Details = details[0]
	} else if len(details) > 1 {
		p.Details = details
	}
}
