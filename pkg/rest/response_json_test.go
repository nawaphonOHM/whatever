package rest

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	resp := JSON(http.StatusAccepted, "accepted payload", "accepted message")
	assert.Equal(t, http.StatusAccepted, resp.StatusCode())
	r := setupTestRouter()
	r.GET("/test-json", func(c *gin.Context) { resp.Write(c) })
	w := executeGet(t, r, "/test-json")
	assert.Equal(t, http.StatusAccepted, w.Code)
}
