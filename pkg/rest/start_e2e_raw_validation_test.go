package rest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPublicRawReservedPath(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, testGroupRaw)
	bp := makeGroupBlueprint(APIGroupRaw, "/health")

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrReservedPath)
	assertPublicRouteRegError(t, buf.Bytes(), "reserved")
}

func testPublicRawCatchAllReserved(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, testGroupRaw)
	bp := makeGroupBlueprint(APIGroupRaw, "/*any")

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrReservedPath)
	assertPublicRouteRegError(t, buf.Bytes(), "reserved")
}

func testPublicRawDuplicateRoute(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, testGroupRaw)
	api1 := &ExportableAPI{Path: "/ping", Method: GET, Handler: func(Context) Response { return OK("1") }}
	api2 := &ExportableAPI{Path: "/ping", Method: GET, Handler: func(Context) Response { return OK("2") }}
	bp := NewBluePrint().WithAPIs(&RRestAPIRegistration{Group: APIGroupRaw, Apis: []*ExportableAPI{api1, api2}})

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrDuplicateRoute)
	assertPublicRouteRegError(t, buf.Bytes(), "duplicate")
}

func testPublicRawPostProbeReserved(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, testGroupRaw)
	apiHealthPost := &ExportableAPI{
		Path:    "/health",
		Method:  POST,
		Handler: func(Context) Response { return OK("health-post") },
	}
	bp := NewBluePrint().WithAPIs(&RRestAPIRegistration{
		Group: APIGroupRaw,
		Apis:  []*ExportableAPI{apiHealthPost},
	})

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrReservedPath)
	assertPublicRouteRegError(t, buf.Bytes(), "reserved")
}

func assertAllowedPostCatchAll(t *testing.T, port int) {
	t.Helper()
	assertStatusOK(t, fmt.Sprintf("http://%s:%d/health", testHost, port))
	assertStatusOK(t, fmt.Sprintf("http://%s:%d/ready", testHost, port))
	status, body := httpPostBody(t, fmt.Sprintf("http://%s:%d/anything/here", testHost, port), "{}")
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "catch-all-post")
}

func testPublicRawAllowedPostCatchAll(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envAPIGroup, testGroupRaw)

	apiCatchAllPost := &ExportableAPI{
		Path:    "/*any",
		Method:  POST,
		Handler: func(Context) Response { return OK("catch-all-post") },
	}
	bp := NewBluePrint().WithAPIs(&RRestAPIRegistration{
		Group: APIGroupRaw,
		Apis:  []*ExportableAPI{apiCatchAllPost},
	})

	errCh := startRESTAsync(bp)
	assertAllowedPostCatchAll(t, port)
	signalAndWait(t, errCh)
}

func TestStartREST_RawValidationAndSafeguards(t *testing.T) {
	t.Run("ReservedPath_Direct", testPublicRawReservedPath)
	t.Run("ReservedPath_CatchAll", testPublicRawCatchAllReserved)
	t.Run("ReservedPath_PostHealth", testPublicRawPostProbeReserved)
	t.Run("DuplicateRoute", testPublicRawDuplicateRoute)
	t.Run("AllowedPostCatchAll", testPublicRawAllowedPostCatchAll)
}
