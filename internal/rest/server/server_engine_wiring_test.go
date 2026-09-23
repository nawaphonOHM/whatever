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

const (
	customReadHeaderSec = 8
	customMaxHeader     = 2097152
	customMaxBody       = 67108864
)

func customWiringConfig() *Config {
	cfg := DefaultConfig()
	cfg.Host = testHost
	cfg.Port = customPort
	cfg.Mode = gin.TestMode
	applyCustomResourceConfig(cfg)
	applyCustomEngineConfig(cfg)
	return cfg
}

func applyCustomResourceConfig(cfg *Config) {
	cfg.ReadHeaderTimeout = customReadHeaderSec * time.Second
	cfg.MaxHeaderBytes = customMaxHeader
	cfg.MaxBodySize = customMaxBody
}

func applyCustomEngineConfig(cfg *Config) {
	cfg.ForwardedByClientIP = false
	cfg.RemoteIPHeaders = []string{"X-Custom-IP"}
	cfg.TrustedProxies = []string{"192.168.1.0/24"}
	cfg.RedirectTrailingSlash = false
	cfg.RedirectFixedPath = true
	cfg.HandleMethodNotAllowed = false
	cfg.UseRawPath = true
	cfg.UnescapePathValues = false
	cfg.RemoveExtraSlash = true
}

func TestNew_WithCustomWiring(t *testing.T) {
	cfg := customWiringConfig()
	srv, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)
	assert.Equal(t, fmt.Sprintf("%s:%d", testHost, customPort), srv.httpServer.Addr)
	assertCustomEngineWiring(t, srv.Engine)
	assertCustomServerWiring(t, srv.httpServer)
}

func assertCustomEngineWiring(t *testing.T, engine *gin.Engine) {
	assert.False(t, engine.RedirectTrailingSlash)
	assert.True(t, engine.RedirectFixedPath)
	assert.False(t, engine.HandleMethodNotAllowed)
	assert.True(t, engine.UseRawPath)
	assert.False(t, engine.UnescapePathValues)
	assert.True(t, engine.RemoveExtraSlash)
	assert.False(t, engine.ForwardedByClientIP)
	assert.Equal(t, []string{"X-Custom-IP"}, engine.RemoteIPHeaders)
	assert.Equal(t, int64(customMaxBody), engine.MaxMultipartMemory)
}

func assertCustomServerWiring(t *testing.T, hs *http.Server) {
	assert.Equal(t, customReadHeaderSec*time.Second, hs.ReadHeaderTimeout)
	assert.Equal(t, customMaxHeader, hs.MaxHeaderBytes)
}
