package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type statusTestCase struct {
	attachError  error
	name         string
	expectedDesc string
	statusCode   int
	expectedCode codes.Code
}

func statusTestCases() []*statusTestCase {
	return []*statusTestCase{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			expectedCode: codes.Ok,
		},
		{
			name:         "404 Not Found",
			statusCode:   http.StatusNotFound,
			expectedCode: codes.Ok,
		},
		{
			name:         "500 Internal Error",
			statusCode:   http.StatusInternalServerError,
			expectedCode: codes.Error,
			expectedDesc: http.StatusText(http.StatusInternalServerError),
		},
		{
			name:         "503 Service Unavailable with Error",
			statusCode:   http.StatusServiceUnavailable,
			expectedCode: codes.Error,
			expectedDesc: "database connection failure",
			attachError:  errors.New("database connection failure"),
		},
	}
}

func setupStatusRoute(t *testing.T, tc *statusTestCase) *gin.Engine {
	t.Helper()
	router := setupTestEngine(Middleware(&config.Config{Enabled: true}))
	router.GET("/status", func(c *gin.Context) {
		if tc.attachError != nil {
			err := c.Error(tc.attachError)
			require.NotNil(t, err)
		}
		c.Status(tc.statusCode)
	})
	return router
}

func runSingleStatusTest(t *testing.T, tc *statusTestCase) {
	if tc == nil {
		return
	}
	exporter, cleanup := setupTestTracer(t)
	defer cleanup()

	router := setupStatusRoute(t, tc)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/status", nil))

	spans := getRecordedSpans(t, exporter)
	require.Len(t, spans, 1)
	verifySpanStatusCode(t, spans[0], tc)
}

func verifySpanStatusCode(t *testing.T, span sdktrace.ReadOnlySpan, tc *statusTestCase) {
	if tc == nil {
		return
	}
	assert.Equal(t, tc.expectedCode, span.Status().Code)
	if tc.expectedDesc != "" {
		assert.Contains(t, span.Status().Description, tc.expectedDesc)
	}
	val, ok := findAttribute(span, AttrHTTPStatusCode)
	require.True(t, ok)
	assert.Equal(t, int64(tc.statusCode), val.AsInt64())
}

func TestMiddleware_Attributes_StatusCodes(t *testing.T) {
	for _, tc := range statusTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			runSingleStatusTest(t, tc)
		})
	}
}
