package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStartREST_ReservedPathRejected blocks /health overrides.
func TestStartREST_ReservedPathRejected(t *testing.T) {
	// Arrange
	t.Setenv(envGinMode, gin.TestMode)
	bp := contracts.NewBluePrint()
	bluePrint := bp.WithAPIs(&contracts.RRestAPIRegistration{
		Apis: []*contracts.ExportableAPI{{
			Path:   ReservedHealthPath,
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK("override")
			},
		}},
	})

	// Act
	err := StartREST(bluePrint)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrReservedPath)
}

// TestStartREST_DuplicateRejected blocks colliding routes.
func TestStartREST_DuplicateRejected(t *testing.T) {
	// Arrange
	t.Setenv(envGinMode, gin.TestMode)
	handler := func(contracts.Context) contracts.Response {
		return testNoContent()
	}
	bluePrint := contracts.NewBluePrint().WithAPIs(
		&contracts.RRestAPIRegistration{
			Version: 1,
			Prefix:  "/x",
			Apis: []*contracts.ExportableAPI{{
				Path: "", Method: contracts.GET, Handler: handler,
			}},
		},
		&contracts.RRestAPIRegistration{
			Version: 1,
			Prefix:  "/x",
			Apis: []*contracts.ExportableAPI{{
				Path: "", Method: contracts.GET, Handler: handler,
			}},
		},
	)

	// Act
	err := StartREST(bluePrint)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrDuplicateRoute)
}

// TestStartREST_NilRegistrationRejected blocks nil groups.
func TestStartREST_NilRegistrationRejected(t *testing.T) {
	// Arrange
	t.Setenv(envGinMode, gin.TestMode)

	// Act
	err := StartREST(contracts.NewBluePrint().WithAPIs(nil))

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilRegistration)
}

// TestStartREST_NilBluePrintRejected rejects a nil blueprint.
func TestStartREST_NilBluePrintRejected(t *testing.T) {
	assert.ErrorIs(t, StartREST(nil), ErrNilBluePrint)
}
