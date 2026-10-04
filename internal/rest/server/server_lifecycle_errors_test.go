package server

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v2/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const logLevelError = "ERROR"

func assertFailedStartLog(t *testing.T, data []byte) {
	t.Helper()
	entries := parseJSONLogs(data)
	failEntry := findLogEntry(entries, "REST server failed to start")
	require.NotNil(t, failEntry, "expected 'REST server failed to start' log entry")
	assert.Equal(t, logLevelError, failEntry["level"])
	assert.NotEmpty(t, failEntry["error"])
}

func TestServer_LifecycleLogging_PortInUse(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	l, port := occupyPort(t)
	defer func() {
		require.NoError(t, l.Close())
	}()

	srv, err := New(&Config{Host: testHost, Port: port, Mode: gin.TestMode})
	require.NoError(t, err)

	require.Error(t, srv.Start(context.Background()))
	assertFailedStartLog(t, buf.Bytes())
}

func testNilBlueprintLog(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	_, err := NewFromBluePrint(nil)
	require.ErrorIs(t, err, ErrNilBluePrint)

	entries := parseJSONLogs(buf.Bytes())
	nilEntry := findLogEntry(entries, "failed to initialize REST server: blueprint is nil")
	require.NotNil(t, nilEntry)
	assert.Equal(t, logLevelError, nilEntry["level"])
}

func assertConfigLoadErrorLog(t *testing.T, data []byte) {
	t.Helper()
	entries := parseJSONLogs(data)
	cfgEntry := findLogEntry(entries, "failed to load REST server config")
	require.NotNil(t, cfgEntry)
	assert.Equal(t, logLevelError, cfgEntry["level"])
	assert.NotEmpty(t, cfgEntry["error"])
}

func testConfigLoadErrorLog(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	t.Setenv("OHM9996_SERVER_PORT", "invalid-port")
	_, err := NewFromBluePrint(contracts.NewBluePrint())
	require.Error(t, err)
	assertConfigLoadErrorLog(t, buf.Bytes())
}

func makeReservedRouteBlueprint() *contracts.BluePrint {
	return contracts.NewBluePrint().WithAPIs(&contracts.RRestAPIRegistration{
		Apis: []*contracts.ExportableAPI{{
			Path:   ReservedHealthPath,
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK("override")
			},
		}},
	})
}

func assertRouteRegErrorLog(t *testing.T, data []byte) {
	t.Helper()
	entries := parseJSONLogs(data)
	routeEntry := findLogEntry(entries, "failed to register REST routes")
	require.NotNil(t, routeEntry)
	assert.Equal(t, logLevelError, routeEntry["level"])
	assert.NotEmpty(t, routeEntry["error"])
}

func testRouteRegErrorLog(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	_, err := NewFromBluePrint(makeReservedRouteBlueprint())
	require.ErrorIs(t, err, ErrReservedPath)
	assertRouteRegErrorLog(t, buf.Bytes())
}

func TestServer_LifecycleLogging_BlueprintErrors(t *testing.T) {
	t.Run("NilBlueprint", testNilBlueprintLog)
	t.Run("ConfigLoadError", testConfigLoadErrorLog)
	t.Run("RouteRegistrationError", testRouteRegErrorLog)
}
