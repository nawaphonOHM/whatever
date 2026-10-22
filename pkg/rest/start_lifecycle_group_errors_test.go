package rest

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertPublicServerCreateError(t *testing.T, data []byte, substr string) {
	t.Helper()
	entries := parseJSONLogEntries(data)
	srvEntry := findLogByMsg(entries, "failed to create REST server")
	require.NotNil(t, srvEntry)
	assert.Equal(t, logLevelError, srvEntry[logKeyLevel])
	assert.Contains(t, srvEntry["error"], substr)
}

func makeGroupBlueprint(group APIGroup, path Pathz) *BluePrint {
	api := &ExportableAPI{Path: path, Method: GET, Handler: func(Context) Response { return OK("ok") }}
	return NewBluePrint().WithAPIs(&RRestAPIRegistration{Group: group, Apis: []*ExportableAPI{api}})
}

func testPublicMixedAPIGroups(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	regStd := &RRestAPIRegistration{Group: APIGroupStandard, Apis: []*ExportableAPI{
		{Path: "/std", Method: GET, Handler: func(Context) Response { return OK("std") }},
	}}
	regRaw := &RRestAPIRegistration{Group: APIGroupRaw, Apis: []*ExportableAPI{
		{Path: "/raw", Method: GET, Handler: func(Context) Response { return OK("raw") }},
	}}
	bp := NewBluePrint().WithAPIs(regStd, regRaw)

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrMixedAPIGroups)
	assertPublicRouteRegError(t, buf.Bytes(), "cannot mix standard and raw")
}

func testPublicInvalidRegistrationGroup(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	bp := makeGroupBlueprint(APIGroup("invalid_group"), "/test")

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrInvalidRegistrationGroup)
	assertPublicRouteRegError(t, buf.Bytes(), "invalid API registration group")
}

func testPublicInvalidEnvGroup(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, "invalid_mode")

	require.ErrorIs(t, startRESTWithTimeout(t, NewBluePrint()), ErrInvalidAPIGroup)
	assertPublicServerCreateError(t, buf.Bytes(), "invalid OHM9996_SERVER_API_GROUP configuration")
}

func testPublicInvalidEnvGroupWithRoutes(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, "invalid_mode")
	bp := makeGroupBlueprint(APIGroupStandard, "/users")

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrInvalidAPIGroup)
	assertPublicServerCreateError(t, buf.Bytes(), "invalid OHM9996_SERVER_API_GROUP configuration")
}

func TestStartREST_GroupLifecycleErrors(t *testing.T) {
	t.Run("MixedAPIGroups", testPublicMixedAPIGroups)
	t.Run("InvalidRegistrationGroup", testPublicInvalidRegistrationGroup)
	t.Run("InvalidEnvGroup_EmptyBP", testPublicInvalidEnvGroup)
	t.Run("InvalidEnvGroup_WithRoutes", testPublicInvalidEnvGroupWithRoutes)
}
