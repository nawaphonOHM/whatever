package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

const (
	otelTestTraceID = "4bf92f3577b34da6a3ce929d0e0e4736"
	otelTestParent  = "00f067aa0ba902b7"
	otelTraceParent = "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
)

// itemsBluePrint creates a blueprint with a parameterized route.
func itemsBluePrint() *contracts.BluePrint {
	reg := &contracts.RRestAPIRegistration{
		Version: 1,
		Prefix:  "/items",
		Apis: []*contracts.ExportableAPI{{
			Path:   "/:id",
			Method: contracts.GET,
			Handler: func(c contracts.Context) contracts.Response {
				return testOK(c.Param("id"))
			},
		}},
	}
	return contracts.NewBluePrint().WithAPIs(reg)
}

// verifyTraceAttributes asserts span attributes for route and method.
func verifyTraceAttributes(t *testing.T, span sdktrace.ReadOnlySpan) {
	t.Helper()
	attrs := make(map[string]any)
	for _, a := range span.Attributes() {
		attrs[string(a.Key)] = a.Value.AsInterface()
	}
	assert.Equal(t, "/api/v1/items/:id", attrs["http.route"])
	assert.Equal(t, "GET", attrs["http.method"])
	assert.Equal(t, int64(http.StatusOK), attrs["http.status_code"])
}

// verifyCorrelatedLog validates trace_id and span_id in structured JSON access log.
func verifyCorrelatedLog(t *testing.T, data []byte, span sdktrace.ReadOnlySpan) {
	t.Helper()
	var entry map[string]any
	require.NoError(t, json.Unmarshal(data, &entry))
	assert.Equal(t, span.SpanContext().TraceID().String(), entry["trace_id"])
	assert.Equal(t, span.SpanContext().SpanID().String(), entry["span_id"])
}

// executeTracedRequest runs an HTTP GET request with traceparent and verifies response.
func executeTracedRequest(t *testing.T, srv *Server) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/items/42", nil)
	req.Header.Set("traceparent", otelTraceParent)
	w := httptest.NewRecorder()
	srv.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("traceparent"), otelTestTraceID)
}

// assertRecordedSpan validates span identifiers, route name, and attributes.
func assertRecordedSpan(t *testing.T, span sdktrace.ReadOnlySpan) {
	t.Helper()
	assert.Equal(t, otelTestTraceID, span.SpanContext().TraceID().String())
	assert.Equal(t, otelTestParent, span.Parent().SpanID().String())
	assert.Equal(t, "GET /api/v1/items/:id", span.Name())
	verifyTraceAttributes(t, span)
}

// TestServer_OTel_TracePropagationAndLogCorrelation verifies traceparent propagation and access log correlation.
func TestServer_OTel_TracePropagationAndLogCorrelation(t *testing.T) {
	setup := setupTraceCorrelation(t)
	defer setup.cleanup()
	defer setup.resetLogger()

	executeTracedRequest(t, setup.server)

	spans := setup.exporter.GetSpans().Snapshots()
	require.Len(t, spans, 1)
	assertRecordedSpan(t, spans[0])
	verifyCorrelatedLog(t, setup.logBuffer.Bytes(), spans[0])
}

type traceCorrelationSetup struct {
	exporter    *tracetest.InMemoryExporter
	cleanup     func()
	resetLogger func()
	logBuffer   *bytes.Buffer
	server      *Server
}

func setupTraceCorrelation(t *testing.T) *traceCorrelationSetup {
	t.Helper()
	exporter, cleanup := setupTestTracer(t)
	buf := new(bytes.Buffer)
	resetLogger := setTestLogger(buf)
	srv := buildTracedServer(t, DefaultConfig(), itemsBluePrint())
	return &traceCorrelationSetup{
		exporter: exporter, cleanup: cleanup, resetLogger: resetLogger,
		logBuffer: buf, server: srv,
	}
}
