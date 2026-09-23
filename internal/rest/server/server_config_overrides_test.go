package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_EnvironmentOverrides(t *testing.T) {
	clearServerConfigEnv(t)
	setServerConfigOverrides(t)
	cfg := loadServerConfig(t)
	assertOverrideBase(t, cfg)
	assertOverrideProxy(t, cfg)
	assertOverrideHTTP(t, cfg)
	assertOverrideResource(t, cfg)
	assertOverrideTelemetry(t, cfg)
}

func assertOverrideBase(t *testing.T, cfg *Config) {
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, overridePort, cfg.Port)
	assert.Equal(t, "test", cfg.Mode)
	assert.Equal(t, "2.0.0", cfg.AppVersion)
	assert.Equal(t, "custom-app", cfg.DisplayName)
	assert.Equal(t, 15*time.Second, cfg.ReadTimeout)
	assert.Equal(t, 15*time.Second, cfg.WriteTimeout)
	assert.Equal(t, 90*time.Second, cfg.IdleTimeout)
	assert.Equal(t, 20*time.Second, cfg.ShutdownTimeout)
}

func assertOverrideProxy(t *testing.T, cfg *Config) {
	assert.Equal(t, []string{"10.0.0.0/8", "192.168.1.1"}, cfg.TrustedProxies)
	assert.Equal(t, []string{"X-Custom-IP", "X-Real-IP"}, cfg.RemoteIPHeaders)
	assert.False(t, cfg.ForwardedByClientIP)
}

func assertOverrideHTTP(t *testing.T, cfg *Config) {
	assert.False(t, cfg.RedirectTrailingSlash)
	assert.True(t, cfg.RedirectFixedPath)
	assert.False(t, cfg.HandleMethodNotAllowed)
	assert.True(t, cfg.UseRawPath)
	assert.False(t, cfg.UnescapePathValues)
	assert.True(t, cfg.RemoveExtraSlash)
}

func assertOverrideResource(t *testing.T, cfg *Config) {
	assert.Equal(t, int64(overrideBodySize), cfg.MaxBodySize)
	assert.Equal(t, 8*time.Second, cfg.ReadHeaderTimeout)
	assert.Equal(t, overrideHeaders, cfg.MaxHeaderBytes)
}

func assertOverrideTelemetry(t *testing.T, cfg *Config) {
	assert.Equal(t, "/custom/pprof", cfg.ProfilePath)
	assert.False(t, cfg.EnableAccessLog)
	assert.True(t, cfg.EnableMetrics)
	assert.True(t, cfg.EnableProfiling)
}
