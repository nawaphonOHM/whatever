package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseJSONLogEntries(data []byte) []map[string]any {
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

func setupTestLogger(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()
	buf := new(bytes.Buffer)
	prev := logging.Default()
	logging.SetDefault(logging.NewJSON(buf, logging.LevelDebug))
	return buf, func() {
		logging.SetDefault(prev)
	}
}

func findLogByMsg(entries []map[string]any, msg string) map[string]any {
	for _, entry := range entries {
		if entry["msg"] == msg {
			return entry
		}
	}
	return nil
}

func TestStartREST_LifecycleLogging_Success(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	port := getFreePort(t)
	setStartEnv(t, port)
	errCh := startRESTAsync(newTestE2EBlueprint())

	assertStatusOK(t, fmt.Sprintf("http://%s:%d/api/v1/ping", testHost, port))

	signalAndWait(t, errCh)

	entries := parseJSONLogEntries(buf.Bytes())

	startEntry := findLogByMsg(entries, "starting REST server")
	require.NotNil(t, startEntry, "expected 'starting REST server' log entry")
	assert.Equal(t, "INFO", startEntry["level"])
	assert.Equal(t, "application", startEntry["service"])
	assert.Equal(t, "1.0.0", startEntry["version"])
	assert.Equal(t, testHost, startEntry["host"])
	assert.Equal(t, float64(port), startEntry["port"])
	assert.Equal(t, gin.TestMode, startEntry["mode"])
	assert.Equal(t, "10s", startEntry["shutdown_timeout"])

	shutdownSignalEntry := findLogByMsg(entries, "shutting down REST server gracefully")
	require.NotNil(t, shutdownSignalEntry, "expected 'shutting down REST server gracefully' log entry")
	assert.Equal(t, "INFO", shutdownSignalEntry["level"])
	assert.Equal(t, "10s", shutdownSignalEntry["shutdown_timeout"])

	shutdownSuccessEntry := findLogByMsg(entries, "REST server shut down successfully")
	require.NotNil(t, shutdownSuccessEntry, "expected 'REST server shut down successfully' log entry")
	assert.Equal(t, "INFO", shutdownSuccessEntry["level"])
}

func TestStartREST_LifecycleLogging_Errors(t *testing.T) {
	t.Run("NilBlueprint", func(t *testing.T) {
		buf, reset := setupTestLogger(t)
		defer reset()

		err := StartREST(nil)
		require.ErrorIs(t, err, ErrNilBluePrint)

		entries := parseJSONLogEntries(buf.Bytes())
		nilEntry := findLogByMsg(entries, "failed to initialize REST server: blueprint is nil")
		require.NotNil(t, nilEntry, "expected nil blueprint log entry")
		assert.Equal(t, "ERROR", nilEntry["level"])
	})

	t.Run("ConfigLoadError", func(t *testing.T) {
		buf, reset := setupTestLogger(t)
		defer reset()

		t.Setenv(envServerPort, "not-a-number")
		err := StartREST(NewBluePrint())
		require.Error(t, err)

		entries := parseJSONLogEntries(buf.Bytes())
		cfgEntry := findLogByMsg(entries, "failed to load REST server config")
		require.NotNil(t, cfgEntry, "expected config load error log entry")
		assert.Equal(t, "ERROR", cfgEntry["level"])
		assert.NotEmpty(t, cfgEntry["error"])
	})

	t.Run("NilRegistration", func(t *testing.T) {
		buf, reset := setupTestLogger(t)
		defer reset()

		t.Setenv(envGinMode, gin.TestMode)
		bp := NewBluePrint().WithAPIs(nil)
		err := StartREST(bp)
		require.ErrorIs(t, err, ErrNilRegistration)

		entries := parseJSONLogEntries(buf.Bytes())
		routeEntry := findLogByMsg(entries, "failed to register REST routes")
		require.NotNil(t, routeEntry, "expected route registration error log entry")
		assert.Equal(t, "ERROR", routeEntry["level"])
		assert.Contains(t, routeEntry["error"], "registration cannot be nil")
	})

	t.Run("ReservedPathConflict", func(t *testing.T) {
		buf, reset := setupTestLogger(t)
		defer reset()

		t.Setenv(envGinMode, gin.TestMode)
		api := &ExportableAPI{
			Path:    "/health",
			Method:  GET,
			Handler: func(Context) Response { return OK("override") },
		}
		reg := &RRestAPIRegistration{
			Apis: []*ExportableAPI{api},
		}
		bp := NewBluePrint().WithAPIs(reg)

		err := StartREST(bp)
		require.ErrorIs(t, err, ErrReservedPath)

		entries := parseJSONLogEntries(buf.Bytes())
		routeEntry := findLogByMsg(entries, "failed to register REST routes")
		require.NotNil(t, routeEntry, "expected route registration error log entry")
		assert.Equal(t, "ERROR", routeEntry["level"])
		assert.Contains(t, routeEntry["error"], "reserved")
	})

	t.Run("DuplicateRouteConflict", func(t *testing.T) {
		buf, reset := setupTestLogger(t)
		defer reset()

		t.Setenv(envGinMode, gin.TestMode)
		api1 := &ExportableAPI{
			Path:    "/ping",
			Method:  GET,
			Handler: func(Context) Response { return OK("pong") },
		}
		api2 := &ExportableAPI{
			Path:    "/ping",
			Method:  GET,
			Handler: func(Context) Response { return OK("pong") },
		}
		reg := &RRestAPIRegistration{
			Apis: []*ExportableAPI{api1, api2},
		}
		bp := NewBluePrint().WithAPIs(reg)

		err := StartREST(bp)
		require.ErrorIs(t, err, ErrDuplicateRoute)

		entries := parseJSONLogEntries(buf.Bytes())
		routeEntry := findLogByMsg(entries, "failed to register REST routes")
		require.NotNil(t, routeEntry, "expected route registration error log entry")
		assert.Equal(t, "ERROR", routeEntry["level"])
		assert.Contains(t, routeEntry["error"], "duplicate")
	})
}

func TestStartREST_LifecycleLogging_PortInUse(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	l, err := net.Listen("tcp", testHost+":0")
	require.NoError(t, err)
	defer func() { assert.NoError(t, l.Close()) }()

	tcpAddr, ok := l.Addr().(*net.TCPAddr)
	require.True(t, ok)

	setStartEnv(t, tcpAddr.Port)
	err = StartREST(NewBluePrint())
	require.Error(t, err)

	entries := parseJSONLogEntries(buf.Bytes())
	failEntry := findLogByMsg(entries, "REST server failed to start")
	require.NotNil(t, failEntry, "expected 'REST server failed to start' log entry")
	assert.Equal(t, "ERROR", failEntry["level"])
	assert.NotEmpty(t, failEntry["error"])
}
