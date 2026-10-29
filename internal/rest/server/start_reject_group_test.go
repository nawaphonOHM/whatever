package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStartREST_RejectMixedGroups rejects blueprints with mixed groups.
func TestStartREST_RejectMixedGroups(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)

	regStd := &contracts.RRestAPIRegistration{
		Group:  contracts.APIGroupStandard,
		Prefix: "/std",
		Apis: []*contracts.ExportableAPI{{
			Path:   "/test",
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK("ok")
			},
		}},
	}
	regRaw := &contracts.RRestAPIRegistration{
		Group:  contracts.APIGroupRaw,
		Prefix: "/raw",
		Apis: []*contracts.ExportableAPI{{
			Path:   "/test",
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK("ok")
			},
		}},
	}

	bp := contracts.NewBluePrint().WithAPIs(regStd, regRaw)
	err := StartREST(bp)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrMixedAPIGroups)
}

// TestStartREST_RejectInvalidRegistrationGroup rejects invalid registration groups.
func TestStartREST_RejectInvalidRegistrationGroup(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)

	bp := contracts.NewBluePrint().WithAPIs(&contracts.RRestAPIRegistration{
		Group:  "invalid_mode",
		Prefix: "/invalid",
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
	assert.ErrorIs(t, err, ErrInvalidRegistrationGroup)
	assert.Contains(t, err.Error(), "invalid_mode")
}
