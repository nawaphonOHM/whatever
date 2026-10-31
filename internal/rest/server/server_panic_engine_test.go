package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	enginePanicPath = "/engine-panic"
	enginePanicMsg  = "engine route exploded"
	enginePanicReq  = "engine-panic-req-id"
	enginePanicLog  = "panic recovered during request processing"
)

func newEnginePanicServer(t *testing.T) *Server {
	t.Helper()
	srv, err := New(&Config{Host: testHost, Port: getFreePort(t), Mode: gin.TestMode})
	require.NoError(t, err)
	srv.SetupDefaultMiddlewares()
	srv.Engine.GET(enginePanicPath, func(*gin.Context) {
		panic(enginePanicMsg)
	})
	return srv
}

func executeEnginePanic(t *testing.T, srv *Server) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, enginePanicPath, nil)
	req.Header.Set(middleware.HeaderXRequestID, enginePanicReq)
	srv.Engine.ServeHTTP(w, req)
	return w
}

func assertEnginePanicFields(t *testing.T, entry map[string]any) {
	t.Helper()
	assert.Equal(t, enginePanicMsg, entry["error"])
	assert.NotEmpty(t, entry["stack_trace"])
	assert.Equal(t, enginePanicReq, entry["request_id"])
	assert.Equal(t, enginePanicPath, entry["path"])
	assert.Equal(t, http.MethodGet, entry["method"])
}

func TestServerEngine_PanicJSONLogEntry(t *testing.T) {
	var buf bytes.Buffer
	reset := core.CaptureLogs(&buf)
	srv := newEnginePanicServer(t)
	w := executeEnginePanic(t, srv)
	reset()

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	entry := findLogEntry(parseJSONLogs(t, buf.Bytes()), enginePanicLog)
	require.NotNil(t, entry, "expected panic log entry")
	assert.Equal(t, "ERROR", entry["level"])
	assertEnginePanicFields(t, entry)
}

func TestServerEngine_PanicGCPLogEntry(t *testing.T) {
	var buf bytes.Buffer
	reset := core.SetTestLogger(core.NewWithHandler(core.NewGCPHandler(&buf, nil)))
	srv := newEnginePanicServer(t)
	w := executeEnginePanic(t, srv)
	reset()

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	entry := findGCPLogEntry(parseJSONLogs(t, buf.Bytes()), enginePanicLog)
	require.NotNil(t, entry, "expected GCP panic log entry")
	assert.Equal(t, "ERROR", entry["severity"])
	assertEnginePanicFields(t, entry)
}
