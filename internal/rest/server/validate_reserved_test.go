package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/pkg/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateRegistrations_ReservedPath rejects /health.
func TestValidateRegistrations_ReservedPath(t *testing.T) {
	// Arrange
	regs := []*rest.RestAPIRegistration{{
		Apis: []*rest.ExportableAPI{{
			Path:    ReservedHealthPath,
			Method:  rest.GET,
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
	regs := []*rest.RestAPIRegistration{
		{
			Prefix: "/items",
			Apis: []*rest.ExportableAPI{{
				Path:    "",
				Method:  rest.GET,
				Handler: dummyHandler,
			}},
		},
		{
			Prefix: "/items",
			Apis: []*rest.ExportableAPI{{
				Path:    "",
				Method:  rest.GET,
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
