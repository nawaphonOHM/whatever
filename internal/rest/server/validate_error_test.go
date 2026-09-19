package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/pkg/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dummyHandler is a no-op handler used by validation tests.
func dummyHandler(*rest.Context) rest.Response {
	return rest.NoContent()
}

// TestValidateRegistrations_NilRegistration rejects nil groups.
func TestValidateRegistrations_NilRegistration(t *testing.T) {
	// Act
	_, err := validateRegistrations([]*rest.RestAPIRegistration{nil})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilRegistration)
}

// TestValidateRegistrations_NilAPI rejects nil API entries.
func TestValidateRegistrations_NilAPI(t *testing.T) {
	// Arrange
	regs := []*rest.RestAPIRegistration{{
		Prefix: "/x",
		Apis:   []*rest.ExportableAPI{nil},
	}}

	// Act
	_, err := validateRegistrations(regs)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilAPI)
}

// TestValidateRegistrations_NilHandler rejects nil handlers.
func TestValidateRegistrations_NilHandler(t *testing.T) {
	// Arrange
	regs := []*rest.RestAPIRegistration{{
		Prefix: "/x",
		Apis: []*rest.ExportableAPI{{
			Path:    "/y",
			Method:  rest.GET,
			Handler: nil,
		}},
	}}

	// Act
	_, err := validateRegistrations(regs)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilHandler)
}

// TestValidateRegistrations_InvalidMethod rejects bad methods.
func TestValidateRegistrations_InvalidMethod(t *testing.T) {
	// Arrange
	regs := []*rest.RestAPIRegistration{{
		Prefix: "/x",
		Apis: []*rest.ExportableAPI{{
			Path:    "/y",
			Method:  rest.HTTPMethod(invalidMethodCode),
			Handler: dummyHandler,
		}},
	}}

	// Act
	_, err := validateRegistrations(regs)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidMethod)
}
