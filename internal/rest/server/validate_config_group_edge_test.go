package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateRegistrationsWithConfig_MixedGroupPriority ensures mixed groups fail before config mismatch.
func TestValidateRegistrationsWithConfig_MixedGroupPriority(t *testing.T) {
	mixedRegs := []*contracts.RRestAPIRegistration{
		makeRegWithGroup(contracts.APIGroupStandard, "/std"),
		makeRegWithGroup(contracts.APIGroupRaw, "/raw"),
	}

	t.Run("mixed blueprint with raw server config returns ErrMixedAPIGroups", func(t *testing.T) {
		cfg := &Config{APIGroup: "raw"}
		_, err := validateRegistrationsWithConfig(mixedRegs, cfg)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMixedAPIGroups)
		assert.NotErrorIs(t, err, ErrAPIGroupMismatch)
	})

	t.Run("mixed blueprint with standard server config returns ErrMixedAPIGroups", func(t *testing.T) {
		cfg := &Config{APIGroup: "standard"}
		_, err := validateRegistrationsWithConfig(mixedRegs, cfg)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMixedAPIGroups)
	})
}

// TestResolveConfigAPIGroup_StrictExactValues verifies rejection of casing and whitespace.
func TestResolveConfigAPIGroup_StrictExactValues(t *testing.T) {
	invalidValues := []string{
		"RAW",
		" raw ",
		"raw ",
		" raw",
		"STANDARD",
		"Standard",
		" standard ",
		"standard ",
		"random_mode",
	}

	for _, val := range invalidValues {
		t.Run(val, func(t *testing.T) {
			_, err := resolveConfigAPIGroup(val)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidAPIGroup)
			assert.Contains(t, err.Error(), val)
		})
	}
}

// TestValidateRegistrationsWithConfig_NilConfig checks nil config fallback to standard.
func TestValidateRegistrationsWithConfig_NilConfig(t *testing.T) {
	t.Run("nil config with standard registrations succeeds", func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{makeRegWithGroup(contracts.APIGroupStandard, "/v1")}
		routes, err := validateRegistrationsWithConfig(regs, nil)
		require.NoError(t, err)
		assert.NotEmpty(t, routes)
	})

	t.Run("nil config with raw registrations fails with mismatch", func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{makeRegWithGroup(contracts.APIGroupRaw, "/v1")}
		_, err := validateRegistrationsWithConfig(regs, nil)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrAPIGroupMismatch)
	})
}
