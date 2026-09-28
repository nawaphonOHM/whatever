package server

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
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

func registerOTelTestGlobals(tp *sdktrace.TracerProvider) {
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

func createOTelTestCleanup(
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

	registerOTelTestGlobals(tp)
	return exporter, createOTelTestCleanup(t, tp, prevTP, prevProp)
}

// setTestLogger configures default slog logger to write to buffer.
func setTestLogger(buf *bytes.Buffer) func() {
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(buf, nil)))
	return func() {
		slog.SetDefault(prev)
	}
}

// buildTracedServer creates a configured Server with blueprint routes and middleware.
func buildTracedServer(
	t *testing.T,
	cfg *Config,
	bp *contracts.BluePrint,
) *Server {
	t.Helper()
	srv, err := New(cfg)
	require.NoError(t, err)
	srv.SetupMiddlewares(buildCORSConfig(bp.Meta()))
	require.NoError(t, srv.RegisterRoutesWithVersion(bp.Apis(), cfg.AppVersion))
	return srv
}
