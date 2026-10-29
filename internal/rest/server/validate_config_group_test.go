package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRegWithGroup(group contracts.APIGroup, prefix string) *contracts.RRestAPIRegistration {
	return &contracts.RRestAPIRegistration{
		Group:  group,
		Prefix: contracts.Pathz(prefix),
		Apis: []*contracts.ExportableAPI{{
			Path:   "/ping",
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK("pong")
			},
		}},
	}
}

// TestValidateRegistrationsWithResult_Success tests valid standard, raw, and default group registrations.
func TestValidateRegistrationsWithResult_Success(t *testing.T) {
	tests := []struct {
		name          string
		expectedPath  string
		expectedGroup contracts.APIGroup
		regs          []*contracts.RRestAPIRegistration
	}{
		{
			name: "standard explicit", expectedPath: "/api/std/ping", expectedGroup: contracts.APIGroupStandard,
			regs: []*contracts.RRestAPIRegistration{makeRegWithGroup(contracts.APIGroupStandard, "/std")},
		},
		{
			name: "default empty", expectedPath: "/api/default/ping", expectedGroup: contracts.APIGroupStandard,
			regs: []*contracts.RRestAPIRegistration{makeRegWithGroup("", "/default")},
		},
		{
			name: "raw group", expectedPath: "/raw/ping", expectedGroup: contracts.APIGroupRaw,
			regs: []*contracts.RRestAPIRegistration{makeRegWithGroup(contracts.APIGroupRaw, "/raw")},
		},
		{
			name: "empty list", expectedGroup: contracts.APIGroupStandard,
			regs: make([]*contracts.RRestAPIRegistration, 0),
		},
		{
			name: "nil slice", expectedGroup: contracts.APIGroupStandard,
			regs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := validateRegistrationsWithResult(tt.regs)
			require.NoError(t, err)
			require.NotNil(t, res)
			assert.Equal(t, tt.expectedGroup, res.group)
			if len(tt.regs) > 0 {
				require.NotEmpty(t, res.routes)
				assert.Equal(t, tt.expectedPath, res.routes[0].Path)
			}
		})
	}
}

// TestValidateRegistrationsWithResult_Errors tests errors returned by validateRegistrationsWithResult.
func TestValidateRegistrationsWithResult_Errors(t *testing.T) {
	t.Run("mixed group returns ErrMixedAPIGroups", func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{
			makeRegWithGroup(contracts.APIGroupStandard, "/std"),
			makeRegWithGroup(contracts.APIGroupRaw, "/raw"),
		}
		res, err := validateRegistrationsWithResult(regs)
		require.ErrorIs(t, err, ErrMixedAPIGroups)
		assert.Nil(t, res)
	})

	t.Run("invalid group returns ErrInvalidRegistrationGroup", func(t *testing.T) {
		res, err := validateRegistrationsWithResult([]*contracts.RRestAPIRegistration{
			makeRegWithGroup("invalid_group", "/invalid"),
		})
		require.ErrorIs(t, err, ErrInvalidRegistrationGroup)
		assert.Nil(t, res)
	})

	t.Run("nil registration in slice returns ErrNilRegistration", func(t *testing.T) {
		res, err := validateRegistrationsWithResult([]*contracts.RRestAPIRegistration{nil})
		require.ErrorIs(t, err, ErrNilRegistration)
		assert.Nil(t, res)
	})
}
