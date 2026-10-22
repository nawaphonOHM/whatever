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

// TestValidateRegistrationsWithConfig_Success tests valid standard, raw, and default group configs.
func TestValidateRegistrationsWithConfig_Success(t *testing.T) {
	tests := []struct {
		name     string
		cfgGroup string
		regGroup contracts.APIGroup
	}{
		{"standard config with standard reg", "standard", contracts.APIGroupStandard},
		{"standard config with default empty reg", "standard", ""},
		{"empty config (default standard) with standard reg", "", contracts.APIGroupStandard},
		{"raw config with raw reg", "raw", contracts.APIGroupRaw},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{APIGroup: tt.cfgGroup}
			regs := []*contracts.RRestAPIRegistration{makeRegWithGroup(tt.regGroup, "/v1")}
			routes, err := validateRegistrationsWithConfig(regs, cfg)
			require.NoError(t, err)
			assert.NotEmpty(t, routes)
		})
	}
}

// TestValidateRegistrationsWithConfig_Mismatch rejects blueprint group differing from server config.
func TestValidateRegistrationsWithConfig_Mismatch(t *testing.T) {
	tests := []struct {
		name     string
		cfgGroup string
		regGroup contracts.APIGroup
	}{
		{"standard config with raw reg", "standard", contracts.APIGroupRaw},
		{"raw config with standard reg", "raw", contracts.APIGroupStandard},
		{"raw config with default empty reg", "raw", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{APIGroup: tt.cfgGroup}
			regs := []*contracts.RRestAPIRegistration{makeRegWithGroup(tt.regGroup, "/v1")}
			routes, err := validateRegistrationsWithConfig(regs, cfg)
			require.Error(t, err)
			assert.Nil(t, routes)
			assert.ErrorIs(t, err, ErrAPIGroupMismatch)
		})
	}
}

// TestValidateRegistrationsWithConfig_InvalidConfig rejects invalid config group even with no regs.
func TestValidateRegistrationsWithConfig_InvalidConfig(t *testing.T) {
	cfg := &Config{APIGroup: "invalid_mode"}

	// with registrations
	regs := []*contracts.RRestAPIRegistration{makeRegWithGroup(contracts.APIGroupStandard, "/v1")}
	_, err := validateRegistrationsWithConfig(regs, cfg)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidAPIGroup)
	assert.Contains(t, err.Error(), "invalid_mode")

	// with empty registrations
	_, err = validateRegistrationsWithConfig(nil, cfg)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidAPIGroup)
}
