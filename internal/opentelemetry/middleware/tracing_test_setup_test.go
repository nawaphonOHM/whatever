package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func registerTestGlobals(tp *sdktrace.TracerProvider) {
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

func createTestCleanup(
	t *testing.T,
	tp *sdktrace.TracerProvider,
	prevTP trace.TracerProvider,
	prevProp propagation.TextMapPropagator,
) func() {
	return func() {
		err := tp.Shutdown(t.Context())
		assert.NoError(t, err)
		otel.SetTracerProvider(prevTP)
		otel.SetTextMapPropagator(prevProp)
	}
}

// setupTestTracer initializes an in-memory exporter and returns a reset function.
func setupTestTracer(t *testing.T) (*tracetest.InMemoryExporter, func()) {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	prevTP := otel.GetTracerProvider()
	prevProp := otel.GetTextMapPropagator()

	registerTestGlobals(tp)
	return exporter, createTestCleanup(t, tp, prevTP, prevProp)
}

// setupTestEngine creates a Gin engine with recovery and tracing middleware.
func setupTestEngine(mw gin.HandlerFunc) *gin.Engine {
	engine := gin.New()
	if mw != nil {
		engine.Use(mw)
	}
	return engine
}

// getRecordedSpans extracts read-only spans from the in-memory exporter.
func getRecordedSpans(t *testing.T, exporter *tracetest.InMemoryExporter) []sdktrace.ReadOnlySpan {
	t.Helper()
	require.NotNil(t, exporter)
	return exporter.GetSpans().Snapshots()
}
