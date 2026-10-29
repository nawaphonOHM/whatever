package rest

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

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

func testPublicNilRegistrationGroup(t *testing.T) {
	buf, reset := setupTestLogger(t)
	defer reset()

	t.Setenv(envGinMode, gin.TestMode)
	bp := NewBluePrint().WithAPIs(nil)

	require.ErrorIs(t, startRESTWithTimeout(t, bp), ErrNilRegistration)
	assertPublicRouteRegError(t, buf.Bytes(), "registration cannot be nil")
}

func TestStartREST_GroupLifecycleErrors(t *testing.T) {
	t.Run("MixedAPIGroups", testPublicMixedAPIGroups)
	t.Run("InvalidRegistrationGroup", testPublicInvalidRegistrationGroup)
	t.Run("NilRegistrationGroup", testPublicNilRegistrationGroup)
}
