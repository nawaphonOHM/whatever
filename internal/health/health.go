package health

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
)

// Status represents the response payload for health probes.
type Status struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// Handler manages internal liveness (/health) and readiness (/ready) health checks.
type Handler struct {
	version string
}

// New creates an internal health Handler with the given version metadata.
func New(version string) *Handler {
	return &Handler{
		version: version,
	}
}

// Health handles GET /health liveness probe requests.
func (h *Handler) Health(c *gin.Context) {
	resp := response.OK(Status{
		Status:    "up",
		Timestamp: time.Now().UTC(),
		Version:   h.version,
	})
	resp.Write(c)
}

// Ready handles GET /ready readiness probe requests.
func (h *Handler) Ready(c *gin.Context) {
	resp := response.OK(Status{
		Status:    "ready",
		Timestamp: time.Now().UTC(),
		Version:   h.version,
	})
	resp.Write(c)
}
