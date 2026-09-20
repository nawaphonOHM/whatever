package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateRegistrations_ReservedPath rejects /health.
func TestValidateRegistrations_ReservedPath(t *testing.T) {
	// Arrange
	regs := []*contracts.RRestAPIRegistration{{
		Apis: []*contracts.ExportableAPI{{
			Path:    ReservedHealthPath,
			Method:  contracts.GET,
			Handler: dummyHandler,
		}},
	}}

	// Act
	_, err := validateRegistrations(regs)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrReservedPath)
}

// TestValidateRegistrations_Duplicate rejects colliding routes.
func TestValidateRegistrations_Duplicate(t *testing.T) {
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
}
