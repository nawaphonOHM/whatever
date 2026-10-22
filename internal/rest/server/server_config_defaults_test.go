package server

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testDefaultTimeout       = 10 * time.Second
	testDefaultIdleTimeout   = 60 * time.Second
	testDefaultHeaderTimeout = 5 * time.Second
)

func TestConfig_DefaultValues(t *testing.T) {
	clearServerConfigEnv(t)
	cfg := loadServerConfig(t)
	assertBaseDefaults(t, cfg)
	assertProxyDefaults(t, cfg)
	assertHTTPDefaults(t, cfg)
	assertResourceDefaults(t, cfg)
	assertTelemetryDefaults(t, cfg)
}

func assertBaseDefaults(t *testing.T, cfg *Config) {
	assert.Empty(t, cfg.Host)
	assert.Equal(t, defaultServerPort, cfg.Port)
	assert.Equal(t, gin.ReleaseMode, cfg.Mode)
	assert.Empty(t, cfg.AppVersion)
	assert.Equal(t, "application", cfg.DisplayName)
	assert.Equal(t, "standard", cfg.APIGroup)
	assert.Equal(t, testDefaultTimeout, cfg.ReadTimeout)
	assert.Equal(t, testDefaultTimeout, cfg.WriteTimeout)
	assert.Equal(t, testDefaultIdleTimeout, cfg.IdleTimeout)
	assert.Equal(t, testDefaultTimeout, cfg.ShutdownTimeout)
}

func assertProxyDefaults(t *testing.T, cfg *Config) {
	assert.Nil(t, cfg.TrustedProxies)
	assert.Equal(t, []string{"X-Forwarded-For", "X-Real-IP"}, cfg.RemoteIPHeaders)
	assert.True(t, cfg.ForwardedByClientIP)
}

func assertHTTPDefaults(t *testing.T, cfg *Config) {
	assert.True(t, cfg.RedirectTrailingSlash)
	assert.False(t, cfg.RedirectFixedPath)
	assert.True(t, cfg.HandleMethodNotAllowed)
	assert.False(t, cfg.UseRawPath)
	assert.True(t, cfg.UnescapePathValues)
	assert.False(t, cfg.RemoveExtraSlash)
	assert.Empty(t, cfg.CORSPath)
}

func assertResourceDefaults(t *testing.T, cfg *Config) {
	assert.Equal(t, int64(defaultMaxBodySize), cfg.MaxBodySize)
	assert.Equal(t, testDefaultHeaderTimeout, cfg.ReadHeaderTimeout)
	assert.Equal(t, defaultMaxHeaderBytes, cfg.MaxHeaderBytes)
}

func assertOTelDefaults(t *testing.T, cfg *Config) {
	assert.NotNil(t, cfg.OTel)
	assert.True(t, cfg.OTel.Enabled)
	assert.Equal(t, "whatever-service", cfg.OTel.ServiceName)
	assert.Equal(t, "localhost:4317", cfg.OTel.Endpoint)
	assert.Equal(t, "grpc", cfg.OTel.Protocol)
	assert.True(t, cfg.OTel.Insecure)
	assert.Equal(t, 1.0, cfg.OTel.SampleRate)
	assert.Nil(t, cfg.OTel.SkipPaths)
}

func assertTelemetryDefaults(t *testing.T, cfg *Config) {
	assert.Equal(t, "/debug/pprof", cfg.ProfilePath)
	assert.True(t, cfg.EnableAccessLog)
	assert.False(t, cfg.EnableMetrics)
	assert.False(t, cfg.EnableProfiling)
	assertOTelDefaults(t, cfg)
}
