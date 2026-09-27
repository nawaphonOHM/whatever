package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const totalConcurrentRequests = 20

func setupConcurrentRouter() *gin.Engine {
	router := setupTestEngine(Middleware(&config.Config{Enabled: true}))
	router.GET("/concurrent/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Header("X-Trace-ID", GetTraceID(c))
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
	return router
}

func sendSingleConcurrentRequest(t *testing.T, router *gin.Engine, index int) {
	path := fmt.Sprintf("/concurrent/%d", index)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("X-Trace-ID"))
}

func TestMiddleware_Concurrency(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	router := setupConcurrentRouter()
	var wg sync.WaitGroup

	for i := range totalConcurrentRequests {
		wg.Go(func() {
			sendSingleConcurrentRequest(t, router, i)
		})
	}
	wg.Wait()

	spans := getRecordedSpans(t, exporter)
	require.Len(t, spans, totalConcurrentRequests)
}
