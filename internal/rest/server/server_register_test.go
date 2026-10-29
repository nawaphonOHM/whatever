package server

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer_RegisterRoutes_RawSuccess tests route registration on a server.
func TestServer_RegisterRoutes_RawSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Server.RegisterRoutes with raw config", func(t *testing.T) {
		srv, err := New(&Config{AppVersion: "1.0.0"})
		require.NoError(t, err)
		require.NoError(t, srv.RegisterRoutes(rawRegs()))

		code, _ := getJSON(t, srv.Engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Server.RegisterRoutesWithVersion with raw config", func(t *testing.T) {
		srv, err := New(&Config{})
		require.NoError(t, err)
		require.NoError(t, srv.RegisterRoutesWithVersion(rawRegs(), "3.0.0"))

		code, _ := getJSON(t, srv.Engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Server.RegisterRoutesWithConfig with nil config on raw server", func(t *testing.T) {
		srv, err := New(&Config{})
		require.NoError(t, err)
		require.NoError(t, srv.RegisterRoutesWithConfig(rawRegs(), nil))

		code, _ := getJSON(t, srv.Engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Server.RegisterRoutesWithConfig custom config", func(t *testing.T) {
		srv, err := New(&Config{})
		require.NoError(t, err)
		require.NoError(t, srv.RegisterRoutesWithConfig(rawRegs(), &Config{AppVersion: "2.0.0"}))

		code, _ := getJSON(t, srv.Engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
	})
}

// TestServer_RegisterRoutes_RejectInvalidRegistrations tests rejection of invalid or mixed blueprints.
func TestServer_RegisterRoutes_RejectInvalidRegistrations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("server rejects mixed group blueprint", func(t *testing.T) {
		srv, err := New(&Config{})
		require.NoError(t, err)
		mixed := append(rawRegs(), stdRegs()...)
		err = srv.RegisterRoutes(mixed)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMixedAPIGroups)
	})

	t.Run("server rejects invalid registration group", func(t *testing.T) {
		srv, err := New(&Config{})
		require.NoError(t, err)
		invalidReg := []*contracts.RRestAPIRegistration{{
			Group:  "invalid_mode",
			Prefix: "/invalid",
			Apis: []*contracts.ExportableAPI{{
				Path:    "/ping",
				Method:  contracts.GET,
				Handler: func(contracts.Context) contracts.Response { return testOK("ok") },
			}},
		}}
		err = srv.RegisterRoutes(invalidReg)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidRegistrationGroup)
	})

	t.Run("server rejects nil registration in slice", func(t *testing.T) {
		srv, err := New(&Config{})
		require.NoError(t, err)
		err = srv.RegisterRoutes([]*contracts.RRestAPIRegistration{nil})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilRegistration)
	})
}
