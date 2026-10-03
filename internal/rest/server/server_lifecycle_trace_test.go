package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func assertTraceIDs(t *testing.T, entry map[string]any, traceID, spanID string) {
	t.Helper()
	require.NotNil(t, entry)
	assert.Equal(t, traceID, entry["trace_id"])
	assert.Equal(t, spanID, entry["span_id"])
}

func assertLifecycleTraceLogs(t *testing.T, entries []map[string]any, traceID, spanID string) {
	t.Helper()
	assertTraceIDs(t, findLogEntry(entries, "starting REST server"), traceID, spanID)
	assertTraceIDs(t, findLogEntry(entries, "shutting down REST server gracefully"), traceID, spanID)
	assertTraceIDs(t, findLogEntry(entries, "REST server shut down successfully"), traceID, spanID)
}

func startTracedServer(t *testing.T, parentCtx context.Context, srv *Server) (context.CancelFunc, <-chan error) {
	t.Helper()
	runCtx, cancel := context.WithCancel(parentCtx)
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start(runCtx) }()
	time.Sleep(startWait)
	return cancel, errCh
}

func executeTracedLifecycle(t *testing.T, ctx context.Context, srv *Server, port int) {
	t.Helper()
	cancel, errCh := startTracedServer(t, ctx, srv)
	assertPingOK(t, port)
	cancel()
	require.NoError(t, waitErr(t, errCh))
}

func assertTraceLogOutput(t *testing.T, data []byte, span trace.Span) {
	t.Helper()
	entries := parseJSONLogs(data)
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	assertLifecycleTraceLogs(t, entries, traceID, spanID)
}

func TestServer_LifecycleLogging_TraceCorrelation(t *testing.T) {
	_, cleanup := setupTestTracer(t)
	defer cleanup()
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	srv, port := newBoundServer(t)
	ctx, span := otel.Tracer("test-lifecycle").Start(context.Background(), "lifecycle-span")
	defer span.End()

	executeTracedLifecycle(t, ctx, srv, port)
	assertTraceLogOutput(t, buf.Bytes(), span)
}
