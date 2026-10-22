package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStartREST_InvalidAPIGroupEnv rejects unknown OHM9996_SERVER_API_GROUP.
func TestStartREST_InvalidAPIGroupEnv(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv("OHM9996_SERVER_API_GROUP", "invalid_group_mode")

	bp := contracts.NewBluePrint()
	err := StartREST(bp)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidAPIGroup)
	assert.Contains(t, err.Error(), "invalid_group_mode")
}

// TestStartREST_MismatchStandardEnvRawBlueprint rejects raw blueprint under standard env.
func TestStartREST_MismatchStandardEnvRawBlueprint(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv("OHM9996_SERVER_API_GROUP", "standard")

	bp := contracts.NewBluePrint().WithAPIs(&contracts.RRestAPIRegistration{
		Group:  contracts.APIGroupRaw,
		Prefix: "/raw",
		Apis: []*contracts.ExportableAPI{{
			Path:   "/test",
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK("ok")
			},
		}},
	})

	err := StartREST(bp)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAPIGroupMismatch)
	assert.Contains(t, err.Error(), "standard")
	assert.Contains(t, err.Error(), "raw")
}

// TestStartREST_MismatchRawEnvStandardBlueprint rejects standard blueprint under raw env.
func TestStartREST_MismatchRawEnvStandardBlueprint(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	t.Setenv("OHM9996_SERVER_API_GROUP", "raw")

	bp := contracts.NewBluePrint().WithAPIs(&contracts.RRestAPIRegistration{
		Group:  contracts.APIGroupStandard,
		Prefix: "/std",
		Apis: []*contracts.ExportableAPI{{
			Path:   "/test",
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK("ok")
			},
		}},
	})

	err := StartREST(bp)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAPIGroupMismatch)
	assert.Contains(t, err.Error(), "raw")
	assert.Contains(t, err.Error(), "standard")
}
