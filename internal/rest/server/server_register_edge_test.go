package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer_ConfigImmutability tests that server Config is not mutated by registration operations.
func TestServer_ConfigImmutability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("server config AppVersion unchanged after RegisterRoutesWithVersion", func(t *testing.T) {
		srv, err := New(&Config{AppVersion: testVersion1})
		require.NoError(t, err)

		err = srv.RegisterRoutesWithVersion(rawRegs(), "9.9.9")
		require.NoError(t, err)
		assert.Equal(t, testVersion1, srv.Config.AppVersion)
	})

	t.Run("server config unchanged after RegisterRoutes with empty slice", func(t *testing.T) {
		srv, err := New(&Config{AppVersion: testVersion1})
		require.NoError(t, err)

		err = srv.RegisterRoutes(nil)
		require.NoError(t, err)
		assert.Equal(t, testVersion1, srv.Config.AppVersion)
	})
}
