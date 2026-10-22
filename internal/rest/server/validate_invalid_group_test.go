package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// invalidGroupCase defines a test case for invalid APIGroup rejection.
type invalidGroupCase struct {
	name      string
	group     contracts.APIGroup
	expectErr string
}

func invalidGroupCases() []*invalidGroupCase {
	return []*invalidGroupCase{
		{
			name:      "cased Raw",
			group:     "Raw",
			expectErr: `registration at index 0 has invalid group "Raw"`,
		},
		{
			name:      "cased STANDARD",
			group:     "STANDARD",
			expectErr: `registration at index 0 has invalid group "STANDARD"`,
		},
		{
			name:      "bogus group value",
			group:     "bogus",
			expectErr: `registration at index 0 has invalid group "bogus"`,
		},
		{
			name:      "unknown group value",
			group:     "custom_group",
			expectErr: `registration at index 0 has invalid group "custom_group"`,
		},
	}
}

// TestValidateRegistrations_InvalidRegistrationGroup rejects unknown or incorrectly cased registration groups.
func TestValidateRegistrations_InvalidRegistrationGroup(t *testing.T) {
	for _, tt := range invalidGroupCases() {
		t.Run(tt.name, func(t *testing.T) {
			regs := []*contracts.RRestAPIRegistration{{
				Group: tt.group,
				Apis: []*contracts.ExportableAPI{{
					Path:    "/test",
					Method:  contracts.GET,
					Handler: dummyHandler,
				}},
			}}

			_, err := validateRegistrations(regs)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidRegistrationGroup)
			assert.Contains(t, err.Error(), tt.expectErr)
		})
	}
}

// TestValidateRegistrations_InvalidRegistrationGroupAtIndex verifies diagnostic index
// for multi-registration blueprints.
func TestValidateRegistrations_InvalidRegistrationGroupAtIndex(t *testing.T) {
	regs := []*contracts.RRestAPIRegistration{
		{
			Group: contracts.APIGroupStandard,
			Apis: []*contracts.ExportableAPI{{
				Path:    "/valid1",
				Method:  contracts.GET,
				Handler: dummyHandler,
			}},
		},
		{
			Group: "bogus",
			Apis: []*contracts.ExportableAPI{{
				Path:    "/invalid",
				Method:  contracts.GET,
				Handler: dummyHandler,
			}},
		},
	}

	_, err := validateRegistrations(regs)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidRegistrationGroup)
	assert.Contains(t, err.Error(), `registration at index 1 has invalid group "bogus"`)
}
