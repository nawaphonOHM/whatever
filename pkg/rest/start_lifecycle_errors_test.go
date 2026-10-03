package rest

import (
	"net"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const logLevelError = "ERROR"

func assertPublicRouteRegError(t *testing.T, data []byte, substr string) {
	t.Helper()
	entries := parseJSONLogEntries(data)
	routeEntry := findLogByMsg(entries, "failed to register REST routes")
	require.NotNil(t, routeEntry)
	assert.Equal(t, logLevelError, routeEntry["level"])
	assert.Contains(t, routeEntry["error"], substr)
}

func makeSingleRouteBlueprint(path Pathz) *BluePrint {
	api := &ExportableAPI{Path: path, Method: GET, Handler: func(Context) Response { return OK("ok") }}
	return NewBluePrint().WithAPIs(&RRestAPIRegistration{Apis: []*ExportableAPI{api}})
}

func makeDuplicateRouteBlueprint(path Pathz) *BluePrint {
	api1 := &ExportableAPI{Path: path, Method: GET, Handler: func(Context) Response { return OK("1") }}
	api2 := &ExportableAPI{Path: path, Method: GET, Handler: func(Context) Response { return OK("2") }}
	return NewBluePrint().WithAPIs(&RRestAPIRegistration{Apis: []*ExportableAPI{api1, api2}})
}

func testPublicNilBlueprint(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	require.ErrorIs(t, StartREST(nil), ErrNilBluePrint)

	entries := parseJSONLogEntries(buf.Bytes())
	nilEntry := findLogByMsg(entries, "failed to initialize REST server: blueprint is nil")
	require.NotNil(t, nilEntry)
	assert.Equal(t, logLevelError, nilEntry["level"])
}

func testPublicConfigError(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envServerPort, "not-a-number")
	require.Error(t, StartREST(NewBluePrint()))

	entries := parseJSONLogEntries(buf.Bytes())
	cfgEntry := findLogByMsg(entries, "failed to load REST server config")
	require.NotNil(t, cfgEntry)
	assert.Equal(t, logLevelError, cfgEntry["level"])
	assert.NotEmpty(t, cfgEntry["error"])
}

func testPublicNilRegistration(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	require.ErrorIs(t, StartREST(NewBluePrint().WithAPIs(nil)), ErrNilRegistration)
	assertPublicRouteRegError(t, buf.Bytes(), "registration cannot be nil")
}

func testPublicReservedPath(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	require.ErrorIs(t, StartREST(makeSingleRouteBlueprint("/health")), ErrReservedPath)
	assertPublicRouteRegError(t, buf.Bytes(), "reserved")
}

func testPublicDuplicateRoute(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	require.ErrorIs(t, StartREST(makeDuplicateRouteBlueprint("/ping")), ErrDuplicateRoute)
	assertPublicRouteRegError(t, buf.Bytes(), "duplicate")
}

func TestStartREST_LifecycleLogging_Errors(t *testing.T) {
	t.Run("NilBlueprint", testPublicNilBlueprint)
	t.Run("ConfigLoadError", testPublicConfigError)
	t.Run("NilRegistration", testPublicNilRegistration)
	t.Run("ReservedPathConflict", testPublicReservedPath)
	t.Run("DuplicateRouteConflict", testPublicDuplicateRoute)
}

func assertPublicPortInUseLog(t *testing.T, data []byte) {
	t.Helper()
	entries := parseJSONLogEntries(data)
	failEntry := findLogByMsg(entries, "REST server failed to start")
	require.NotNil(t, failEntry)
	assert.Equal(t, logLevelError, failEntry["level"])
	assert.NotEmpty(t, failEntry["error"])
}

func occupyTCPPort(t *testing.T) (net.Listener, int) {
	t.Helper()
	l, err := net.Listen("tcp", testHost+":0")
	require.NoError(t, err)
	tcpAddr, ok := l.Addr().(*net.TCPAddr)
	require.True(t, ok)
	return l, tcpAddr.Port
}

func TestStartREST_LifecycleLogging_PortInUse(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	l, port := occupyTCPPort(t)
	defer func() { require.NoError(t, l.Close()) }()

	setStartEnv(t, port)
	require.Error(t, StartREST(NewBluePrint()))
	assertPublicPortInUseLog(t, buf.Bytes())
}
