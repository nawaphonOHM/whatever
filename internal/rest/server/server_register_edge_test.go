package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer_RegisterRoutes_MutatedInvalidGroup tests that invalid server group returns ErrInvalidAPIGroup.
func TestServer_RegisterRoutes_MutatedInvalidGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("mutated srv.Config.APIGroup fails RegisterRoutes", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: standardGroupMode})
		require.NoError(t, err)
		srv.Config.APIGroup = "invalid_mode"

		err = srv.RegisterRoutes(stdRegs())
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidAPIGroup)
	})

	t.Run("mutated srv.Config.APIGroup fails RegisterRoutesWithVersion", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: standardGroupMode})
		require.NoError(t, err)
		srv.Config.APIGroup = "invalid_mode"

		err = srv.RegisterRoutesWithVersion(stdRegs(), testVersion1)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidAPIGroup)
	})

	t.Run("mutated srv.Config.APIGroup fails RegisterRoutesWithConfig nil", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: standardGroupMode})
		require.NoError(t, err)
		srv.Config.APIGroup = "invalid_mode"

		err = srv.RegisterRoutesWithConfig(stdRegs(), nil)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidAPIGroup)
	})

	t.Run("manually constructed server with invalid APIGroup fails", func(t *testing.T) {
		srv := &Server{
			Engine: gin.New(),
			Config: &Config{APIGroup: "bad_mode"},
		}
		err := srv.RegisterRoutes(stdRegs())
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidAPIGroup)
	})
}

// TestServer_ConfigImmutability tests that server Config is not mutated by registration operations.
func TestServer_ConfigImmutability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("server config unchanged after rejected override", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: rawGroupMode, AppVersion: testVersion1})
		require.NoError(t, err)

		err = srv.RegisterRoutesWithConfig(stdRegs(), &Config{APIGroup: standardGroupMode})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrAPIGroupMismatch)
		assert.Equal(t, rawGroupMode, srv.Config.APIGroup)
		assert.Equal(t, testVersion1, srv.Config.AppVersion)
	})

	t.Run("server config AppVersion unchanged after RegisterRoutesWithVersion", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: rawGroupMode, AppVersion: testVersion1})
		require.NoError(t, err)

		err = srv.RegisterRoutesWithVersion(rawRegs(), "9.9.9")
		require.NoError(t, err)
		assert.Equal(t, testVersion1, srv.Config.AppVersion)
		assert.Equal(t, rawGroupMode, srv.Config.APIGroup)
	})
}
