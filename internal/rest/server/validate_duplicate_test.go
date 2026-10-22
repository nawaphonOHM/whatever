package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateRegistrations_Duplicate rejects colliding routes.
func TestValidateRegistrations_Duplicate(t *testing.T) {
	t.Run("identical registration", func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{
			singleAPI(testItemsPrefix, "", contracts.APIGroupStandard),
			singleAPI(testItemsPrefix, "", contracts.APIGroupStandard),
		}

		_, err := validateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateRoute)
	})

	t.Run("canonical path collision with different formatting", func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{
			{
				Version: 1,
				Prefix:  testItemsPrefix,
				Apis: []*contracts.ExportableAPI{{
					Path:    testListPath,
					Method:  contracts.GET,
					Handler: dummyHandler,
				}},
			},
			{
				Version: 1,
				Prefix:  "/api/v1/items",
				Apis: []*contracts.ExportableAPI{{
					Path:    "///list///",
					Method:  contracts.GET,
					Handler: dummyHandler,
				}},
			},
		}

		_, err := validateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateRoute)
	})

	t.Run("raw mode duplicate routes across registrations", func(t *testing.T) {
		raw := contracts.APIGroupRaw
		regs := []*contracts.RRestAPIRegistration{
			{
				Group:  raw,
				Prefix: "/custom",
				Apis: []*contracts.ExportableAPI{{
					Path:    "/hook",
					Method:  contracts.POST,
					Handler: dummyHandler,
				}},
			},
			{
				Group:  raw,
				Prefix: "///custom///",
				Apis: []*contracts.ExportableAPI{{
					Path:    "///hook///",
					Method:  contracts.POST,
					Handler: dummyHandler,
				}},
			},
		}

		_, err := validateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateRoute)
	})
}
