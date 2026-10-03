package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/nawaphonOHM/whatever/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func parseJSONLogs(data []byte) []map[string]any {
	var entries []map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	for decoder.More() {
		var entry map[string]any
		if err := decoder.Decode(&entry); err == nil {
			entries = append(entries, entry)
		}
	}
	return entries
}

func setupLifecycleLogger(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()
	buf := new(bytes.Buffer)
	prev := logging.Default()
	logging.SetDefault(logging.NewJSON(buf, logging.LevelDebug))
	return buf, func() {
		logging.SetDefault(prev)
	}
}

func findLogEntry(entries []map[string]any, msg string) map[string]any {
	for _, entry := range entries {
		if entry["msg"] == msg {
			return entry
		}
	}
	return nil
}

func TestServer_LifecycleLogging_StartupAndShutdown(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	srv, port := newBoundServer(t)
	cancel, errCh := startWithCancel(srv)

	assertPingOK(t, port)

	cancel()
	require.NoError(t, waitErr(t, errCh))

	entries := parseJSONLogs(buf.Bytes())

	startEntry := findLogEntry(entries, "starting REST server")
	require.NotNil(t, startEntry, "expected 'starting REST server' log entry")
	assert.Equal(t, "INFO", startEntry["level"])
	assert.Equal(t, srv.Config.DisplayName, startEntry["service"])
	assert.Equal(t, testHost, startEntry["host"])
	assert.Equal(t, float64(port), startEntry["port"])
	assert.Equal(t, gin.TestMode, startEntry["mode"])
	assert.Equal(t, srv.Config.ShutdownTimeout.String(), startEntry["shutdown_timeout"])
	assert.Equal(t, srv.Config.EnableAccessLog, startEntry["access_log"])
	assert.Equal(t, srv.Config.EnableMetrics, startEntry["metrics"])

	shutdownSignalEntry := findLogEntry(entries, "shutting down REST server gracefully")
	require.NotNil(t, shutdownSignalEntry, "expected 'shutting down REST server gracefully' log entry")
	assert.Equal(t, "INFO", shutdownSignalEntry["level"])
	assert.Equal(t, srv.Config.ShutdownTimeout.String(), shutdownSignalEntry["shutdown_timeout"])

	shutdownSuccessEntry := findLogEntry(entries, "REST server shut down successfully")
	require.NotNil(t, shutdownSuccessEntry, "expected 'REST server shut down successfully' log entry")
	assert.Equal(t, "INFO", shutdownSuccessEntry["level"])
}

func TestServer_LifecycleLogging_PortInUse(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	l, port := occupyPort(t)
	defer func() { assert.NoError(t, l.Close()) }()

	srv, err := New(&Config{
		Host: testHost,
		Port: port,
		Mode: gin.TestMode,
	})
	require.NoError(t, err)

	err = srv.Start(context.Background())
	require.Error(t, err)

	entries := parseJSONLogs(buf.Bytes())
	failEntry := findLogEntry(entries, "REST server failed to start")
	require.NotNil(t, failEntry, "expected 'REST server failed to start' log entry")
	assert.Equal(t, "ERROR", failEntry["level"])
	assert.NotEmpty(t, failEntry["error"])
}

func TestServer_LifecycleLogging_GracefulShutdownTimeout(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	port := getFreePort(t)
	srv, err := New(&Config{
		Host:          testHost,
		Port:          port,
		Mode:          gin.TestMode,
		TimeoutFields: TimeoutFields{ShutdownTimeout: 20 * time.Millisecond},
	})
	require.NoError(t, err)

	started := make(chan struct{})
	blockChan := make(chan struct{})
	srv.Engine.GET("/slow", func(c *gin.Context) {
		close(started)
		<-blockChan
		c.String(http.StatusOK, "slow response")
	})

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start(ctx) }()
	time.Sleep(startWait)

	// Fire slow request and wait until it is being handled
	go func() {
		resp, _ := http.Get(fmt.Sprintf("http://%s:%d/slow", testHost, port))
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	select {
	case <-started:
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for slow request to start")
	}

	cancel()
	err = waitErr(t, errCh)
	close(blockChan)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "graceful shutdown failed")

	entries := parseJSONLogs(buf.Bytes())
	shutdownFailEntry := findLogEntry(entries, "REST server graceful shutdown failed")
	require.NotNil(t, shutdownFailEntry, "expected 'REST server graceful shutdown failed' log entry")
	assert.Equal(t, "ERROR", shutdownFailEntry["level"])
	assert.NotEmpty(t, shutdownFailEntry["error"])
}

func TestServer_LifecycleLogging_BlueprintErrors(t *testing.T) {
	t.Run("NilBlueprint", func(t *testing.T) {
		buf, reset := setupLifecycleLogger(t)
		defer reset()

		_, err := NewFromBluePrint(nil)
		require.ErrorIs(t, err, ErrNilBluePrint)

		entries := parseJSONLogs(buf.Bytes())
		nilEntry := findLogEntry(entries, "failed to initialize REST server: blueprint is nil")
		require.NotNil(t, nilEntry, "expected nil blueprint log entry")
		assert.Equal(t, "ERROR", nilEntry["level"])
	})

	t.Run("ConfigLoadError", func(t *testing.T) {
		buf, reset := setupLifecycleLogger(t)
		defer reset()

		t.Setenv("OHM9996_SERVER_PORT", "invalid-port")
		bp := contracts.NewBluePrint()
		_, err := NewFromBluePrint(bp)
		require.Error(t, err)

		entries := parseJSONLogs(buf.Bytes())
		cfgEntry := findLogEntry(entries, "failed to load REST server config")
		require.NotNil(t, cfgEntry, "expected config load error log entry")
		assert.Equal(t, "ERROR", cfgEntry["level"])
		assert.NotEmpty(t, cfgEntry["error"])
	})

	t.Run("RouteRegistrationError", func(t *testing.T) {
		buf, reset := setupLifecycleLogger(t)
		defer reset()

		t.Setenv(envGinMode, gin.TestMode)
		bp := contracts.NewBluePrint().WithAPIs(&contracts.RRestAPIRegistration{
			Apis: []*contracts.ExportableAPI{{
				Path:   ReservedHealthPath,
				Method: contracts.GET,
				Handler: func(contracts.Context) contracts.Response {
					return testOK("override")
				},
			}},
		})
		_, err := NewFromBluePrint(bp)
		require.ErrorIs(t, err, ErrReservedPath)

		entries := parseJSONLogs(buf.Bytes())
		routeEntry := findLogEntry(entries, "failed to register REST routes")
		require.NotNil(t, routeEntry, "expected route registration error log entry")
		assert.Equal(t, "ERROR", routeEntry["level"])
		assert.NotEmpty(t, routeEntry["error"])
	})
}

func TestServer_LifecycleLogging_TraceCorrelation(t *testing.T) {
	_, cleanup := setupTestTracer(t)
	defer cleanup()

	buf, reset := setupLifecycleLogger(t)
	defer reset()

	srv, port := newBoundServer(t)

	tracer := otel.Tracer("test-lifecycle")
	ctx, span := tracer.Start(context.Background(), "lifecycle-span")
	defer span.End()

	ctx, cancel := context.WithCancel(ctx)
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start(ctx) }()
	time.Sleep(startWait)

	assertPingOK(t, port)

	cancel()
	require.NoError(t, waitErr(t, errCh))

	entries := parseJSONLogs(buf.Bytes())
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	startEntry := findLogEntry(entries, "starting REST server")
	require.NotNil(t, startEntry)
	assert.Equal(t, traceID, startEntry["trace_id"])
	assert.Equal(t, spanID, startEntry["span_id"])

	shutdownSignalEntry := findLogEntry(entries, "shutting down REST server gracefully")
	require.NotNil(t, shutdownSignalEntry)
	assert.Equal(t, traceID, shutdownSignalEntry["trace_id"])
	assert.Equal(t, spanID, shutdownSignalEntry["span_id"])

	shutdownSuccessEntry := findLogEntry(entries, "REST server shut down successfully")
	require.NotNil(t, shutdownSuccessEntry)
	assert.Equal(t, traceID, shutdownSuccessEntry["trace_id"])
	assert.Equal(t, spanID, shutdownSuccessEntry["span_id"])
}
