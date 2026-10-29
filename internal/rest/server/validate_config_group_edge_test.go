package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testEmptyPrefix = "/empty"

// TestResolveAPIGroup_StrictExactValues verifies rejection of casing and whitespace in registration groups.
func TestResolveAPIGroup_StrictExactValues(t *testing.T) {
	invalidValues := []contracts.APIGroup{
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
		t.Run(string(val), func(t *testing.T) {
			_, err := resolveAPIGroup(val, 0)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidRegistrationGroup)
			assert.Contains(t, err.Error(), string(val))
		})
	}
}

// TestValidateRegistrationsWithResult_DefaultAndExplicitCombinations tests combinations of empty and explicit groups.
func TestValidateRegistrationsWithResult_DefaultAndExplicitCombinations(t *testing.T) {
	t.Run("default empty then explicit standard succeeds", func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{
			makeRegWithGroup("", testEmptyPrefix),
			makeRegWithGroup(contracts.APIGroupStandard, "/std"),
		}
		res, err := validateRegistrationsWithResult(regs)
		require.NoError(t, err)
		assert.Equal(t, contracts.APIGroupStandard, res.group)
		assert.Len(t, res.routes, 2)
	})

	t.Run("explicit standard then default empty succeeds", func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{
			makeRegWithGroup(contracts.APIGroupStandard, "/std"),
			makeRegWithGroup("", testEmptyPrefix),
		}
		res, err := validateRegistrationsWithResult(regs)
		require.NoError(t, err)
		assert.Equal(t, contracts.APIGroupStandard, res.group)
		assert.Len(t, res.routes, 2)
	})

	t.Run("raw then default empty fails with mixed error", func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{
			makeRegWithGroup(contracts.APIGroupRaw, "/raw"),
			makeRegWithGroup("", testEmptyPrefix),
		}
		res, err := validateRegistrationsWithResult(regs)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrMixedAPIGroups)
	})

	t.Run("default empty then raw fails with mixed error", func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{
			makeRegWithGroup("", testEmptyPrefix),
			makeRegWithGroup(contracts.APIGroupRaw, "/raw"),
		}
		res, err := validateRegistrationsWithResult(regs)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrMixedAPIGroups)
	})
}
