package rest

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func testPublicStdEnvRawBlueprint(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, testGroupStd)
	bp := makeGroupBlueprint(APIGroupRaw, testCustomPath)

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrAPIGroupMismatch)
	assertPublicRouteRegError(t, buf.Bytes(), errMismatchSubstr)
}

func testPublicRawEnvStdBlueprint(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, testGroupRaw)
	bp := makeGroupBlueprint(APIGroupStandard, testCustomPath)

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrAPIGroupMismatch)
	assertPublicRouteRegError(t, buf.Bytes(), errMismatchSubstr)
}

func testPublicRawEnvDefaultBlueprint(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, testGroupRaw)
	bp := makeGroupBlueprint("", testCustomPath)

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrAPIGroupMismatch)
	assertPublicRouteRegError(t, buf.Bytes(), errMismatchSubstr)
}

func testPublicDefaultEnvRawBlueprint(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv(envAPIGroup, "")
	bp := makeGroupBlueprint(APIGroupRaw, testCustomPath)

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrAPIGroupMismatch)
	assertPublicRouteRegError(t, buf.Bytes(), errMismatchSubstr)
}

func TestStartREST_GroupMismatchErrors(t *testing.T) {
	t.Run("StdEnv_RawBlueprint", testPublicStdEnvRawBlueprint)
	t.Run("RawEnv_StdBlueprint", testPublicRawEnvStdBlueprint)
	t.Run("RawEnv_DefaultBlueprint", testPublicRawEnvDefaultBlueprint)
	t.Run("DefaultEnv_RawBlueprint", testPublicDefaultEnvRawBlueprint)
}
