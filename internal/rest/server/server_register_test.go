package server

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer_RegisterRoutes_RawSuccess tests route registration on a raw server.
func TestServer_RegisterRoutes_RawSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Server.RegisterRoutes with raw config", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: rawGroupMode, AppVersion: "1.0.0"})
		require.NoError(t, err)
		require.NoError(t, srv.RegisterRoutes(rawRegs()))

		code, _ := getJSON(t, srv.Engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Server.RegisterRoutesWithVersion with raw config", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: rawGroupMode})
		require.NoError(t, err)
		require.NoError(t, srv.RegisterRoutesWithVersion(rawRegs(), "3.0.0"))

		code, _ := getJSON(t, srv.Engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Server.RegisterRoutesWithConfig with nil config on raw server", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: rawGroupMode})
		require.NoError(t, err)
		require.NoError(t, srv.RegisterRoutesWithConfig(rawRegs(), nil))

		code, _ := getJSON(t, srv.Engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Server.RegisterRoutesWithConfig matching raw override", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: rawGroupMode})
		require.NoError(t, err)
		require.NoError(t, srv.RegisterRoutesWithConfig(rawRegs(), &Config{APIGroup: rawGroupMode}))

		code, _ := getJSON(t, srv.Engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Server.RegisterRoutesWithConfig empty override inherits raw group", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: rawGroupMode})
		require.NoError(t, err)
		require.NoError(t, srv.RegisterRoutesWithConfig(rawRegs(), &Config{AppVersion: "2.0.0"}))

		code, _ := getJSON(t, srv.Engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
	})
}

// TestServer_RegisterRoutes_RejectConflictingOverride rejects changing APIGroup at runtime.
func TestServer_RegisterRoutes_RejectConflictingOverride(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("raw server rejects standard override", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: rawGroupMode})
		require.NoError(t, err)
		err = srv.RegisterRoutesWithConfig(stdRegs(), &Config{APIGroup: "standard"})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrAPIGroupMismatch)
		assert.Contains(t, err.Error(), `config override group "standard" conflicts with server group "raw"`)
	})

	t.Run("standard server rejects raw override", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: "standard"})
		require.NoError(t, err)
		err = srv.RegisterRoutesWithConfig(rawRegs(), &Config{APIGroup: rawGroupMode})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrAPIGroupMismatch)
		assert.Contains(t, err.Error(), `config override group "raw" conflicts with server group "standard"`)
	})

	t.Run("server rejects invalid override group", func(t *testing.T) {
		srv, err := New(&Config{APIGroup: rawGroupMode})
		require.NoError(t, err)
		err = srv.RegisterRoutesWithConfig(rawRegs(), &Config{APIGroup: "invalid_override"})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidAPIGroup)
	})
}
