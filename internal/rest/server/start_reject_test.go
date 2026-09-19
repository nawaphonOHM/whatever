package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStartREST_ReservedPathRejected blocks /health overrides.
func TestStartREST_ReservedPathRejected(t *testing.T) {
	// Arrange
	t.Setenv(envGinMode, gin.TestMode)
	regs := []*rest.RestAPIRegistration{{
		Apis: []*rest.ExportableAPI{{
			Path:   ReservedHealthPath,
			Method: rest.GET,
			Handler: func(*rest.Context) rest.Response {
				return rest.OK("override")
			},
		}},
	}}

	// Act
	err := StartREST(regs)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrReservedPath)
}

// TestStartREST_DuplicateRejected blocks colliding routes.
func TestStartREST_DuplicateRejected(t *testing.T) {
	// Arrange
	t.Setenv(envGinMode, gin.TestMode)
	handler := func(*rest.Context) rest.Response {
		return rest.NoContent()
	}
	regs := []*rest.RestAPIRegistration{
		{
			Version: 1,
			Prefix:  "/x",
			Apis: []*rest.ExportableAPI{{
				Path: "", Method: rest.GET, Handler: handler,
			}},
		},
		{
			Version: 1,
			Prefix:  "/x",
			Apis: []*rest.ExportableAPI{{
				Path: "", Method: rest.GET, Handler: handler,
			}},
		},
	}

	// Act
	err := StartREST(regs)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrDuplicateRoute)
}

// TestStartREST_NilRegistrationRejected blocks nil groups.
func TestStartREST_NilRegistrationRejected(t *testing.T) {
	// Arrange
	t.Setenv(envGinMode, gin.TestMode)

	// Act
	err := StartREST([]*rest.RestAPIRegistration{nil})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilRegistration)
}
