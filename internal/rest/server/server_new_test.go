package server

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew_WithNilConfig uses package defaults.
func TestNew_WithNilConfig(t *testing.T) {
	srv, err := New(nil)
	require.NoError(t, err)
	require.NotNil(t, srv)
	require.NotNil(t, srv.Engine)
	require.NotNil(t, srv.Config)
	assert.Equal(t, defaultServerPort, srv.Config.Port)
	assert.Equal(t, gin.ReleaseMode, srv.Config.Mode)
	assert.Equal(t, fmt.Sprintf(":%d", defaultServerPort), srv.httpServer.Addr)
	assertDefaultEngineWiring(t, srv.Engine)
	assertDefaultServerWiring(t, srv.httpServer)
}

func assertDefaultEngineWiring(t *testing.T, engine *gin.Engine) {
	assert.True(t, engine.RedirectTrailingSlash)
	assert.False(t, engine.RedirectFixedPath)
	assert.True(t, engine.HandleMethodNotAllowed)
	assert.False(t, engine.UseRawPath)
	assert.True(t, engine.UnescapePathValues)
	assert.False(t, engine.RemoveExtraSlash)
	assert.True(t, engine.ForwardedByClientIP)
	assert.Equal(t, []string{"X-Forwarded-For", "X-Real-IP"}, engine.RemoteIPHeaders)
	assert.Equal(t, int64(defaultMaxBodySize), engine.MaxMultipartMemory)
}

func assertDefaultServerWiring(t *testing.T, hs *http.Server) {
	assert.Equal(t, defaultReadHeaderSec*time.Second, hs.ReadHeaderTimeout)
	assert.Equal(t, defaultMaxHeaderBytes, hs.MaxHeaderBytes)
	assert.Equal(t, defaultTimeoutSec*time.Second, hs.ReadTimeout)
	assert.Equal(t, defaultTimeoutSec*time.Second, hs.WriteTimeout)
	assert.Equal(t, defaultIdleTimeoutSec*time.Second, hs.IdleTimeout)
}

// TestNew_InvalidTrustedProxies fails when proxies format is invalid.
func TestNew_InvalidTrustedProxies(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TrustedProxies = []string{"invalid-cidr-or-ip"}
	srv, err := New(cfg)
	require.Error(t, err)
	assert.Nil(t, srv)
	assert.Contains(t, err.Error(), "failed to set trusted proxies")
}

// TestDefaultConfig_Values checks recommended defaults.
func TestDefaultConfig_Values(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	assert.Equal(t, defaultServerPort, cfg.Port)
	assert.Equal(t, gin.ReleaseMode, cfg.Mode)
	assert.Equal(t, defaultTimeoutSec*time.Second, cfg.ReadTimeout)
	assert.Equal(t, defaultTimeoutSec*time.Second, cfg.WriteTimeout)
	assert.Equal(t, defaultIdleTimeoutSec*time.Second, cfg.IdleTimeout)
	assert.Equal(t, defaultTimeoutSec*time.Second, cfg.ShutdownTimeout)
}
