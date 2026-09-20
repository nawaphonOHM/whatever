// Package health provides health check handlers for framework liveness
// and readiness probes.
package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Status represents the response payload for health probes.
type Status struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// responseEnvelope represents the JSON response envelope for health probes.
type responseEnvelope struct {
	Timestamp time.Time `json:"timestamp"`
	Data      Status    `json:"data"`
	Success   bool      `json:"success"`
}

// Handler manages internal liveness (/health) and readiness (/ready)
// health checks.
type Handler struct {
	version string
}

// New creates an internal health Handler with the given version metadata.
func New(version string) *Handler {
	return &Handler{
		version: version,
	}
}

func (h *Handler) writeStatus(c *gin.Context, status string) {
	c.JSON(http.StatusOK, responseEnvelope{
		Data: Status{
			Status:    status,
			Timestamp: time.Now().UTC(),
			Version:   h.version,
		},
		Timestamp: time.Now().UTC(),
		Success:   true,
	})
}

// Health handles GET /health liveness probe requests.
func (h *Handler) Health(c *gin.Context) {
	h.writeStatus(c, "up")
}

// Ready handles GET /ready readiness probe requests.
func (h *Handler) Ready(c *gin.Context) {
	h.writeStatus(c, "ready")
}
