package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateRegistrations_ReservedPath rejects /health and /ready.
func TestValidateRegistrations_ReservedPath(t *testing.T) {
	for _, reserved := range []string{ReservedHealthPath, ReservedReadyPath} {
		t.Run(reserved, func(t *testing.T) {
			// Arrange
			regs := []*contracts.RRestAPIRegistration{{
				Apis: []*contracts.ExportableAPI{{
					Path:    contracts.Pathz(reserved),
					Method:  contracts.GET,
					Handler: dummyHandler,
				}},
			}}

			// Act
			_, err := validateRegistrations(regs)

			// Assert
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrReservedPath)
		})
	}
}

// TestValidateRegistrations_Duplicate rejects colliding routes.
func TestValidateRegistrations_Duplicate(t *testing.T) {
	t.Run("identical registration", func(t *testing.T) {
		// Arrange
		regs := []*contracts.RRestAPIRegistration{
			{
				Prefix: "/items",
				Apis: []*contracts.ExportableAPI{{
					Path:    "",
					Method:  contracts.GET,
					Handler: dummyHandler,
				}},
			},
			{
				Prefix: "/items",
				Apis: []*contracts.ExportableAPI{{
					Path:    "",
					Method:  contracts.GET,
					Handler: dummyHandler,
				}},
			},
		}

		// Act
		_, err := validateRegistrations(regs)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateRoute)
	})

	t.Run("canonical path collision with different formatting", func(t *testing.T) {
		// Arrange
		regs := []*contracts.RRestAPIRegistration{
			{
				Version: 1,
				Prefix:  "/items",
				Apis: []*contracts.ExportableAPI{{
					Path:    "/list",
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

		// Act
		_, err := validateRegistrations(regs)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateRoute)
	})
}
