package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runDisabledMiddlewareTest(t *testing.T, cfg *config.Config) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	router := setupTestEngine(Middleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, getRecordedSpans(t, exporter))
}

func TestMiddleware_Disabled(t *testing.T) {
	testCases := []struct {
		cfg  *config.Config
		name string
	}{
		{
			cfg:  nil,
			name: "nil config",
		},
		{
			cfg:  &config.Config{Enabled: false},
			name: "explicitly disabled config",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runDisabledMiddlewareTest(t, tc.cfg)
		})
	}
}

func TestMiddleware_TraceMiddlewareAlias(t *testing.T) {
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	cfg := &config.Config{Enabled: true}
	router := setupTestEngine(TraceMiddleware(cfg))
	router.GET("/alias", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/alias", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	spans := getRecordedSpans(t, exporter)
	require.Len(t, spans, 1)
}
