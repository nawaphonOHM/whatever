package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStartREST_ReservedPathRejected blocks /health overrides.
func TestStartREST_ReservedPathRejected(t *testing.T) {
	// Arrange
	t.Setenv(envGinMode, gin.TestMode)
	regs := []*RestAPIRegistration{{
		Apis: []*ExportableAPI{{
			Path:   ReservedHealthPath,
			Method: GET,
			Handler: func(*Context) response.Response {
				return response.OK("override")
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
	handler := func(*Context) response.Response {
		return response.NoContent()
	}
	regs := []*RestAPIRegistration{
		{
			Version: 1,
			Prefix:  "/x",
			Apis: []*ExportableAPI{{
				Path: "", Method: GET, Handler: handler,
			}},
		},
		{
			Version: 1,
			Prefix:  "/x",
			Apis: []*ExportableAPI{{
				Path: "", Method: GET, Handler: handler,
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
	err := StartREST([]*RestAPIRegistration{nil})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilRegistration)
}
