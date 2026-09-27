package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testCustomSkipPath = "/custom"
	testOtherPath      = "/other"
)

// TestConfig_ShouldSkip tests path matching against configured SkipPaths.
func TestConfig_ShouldSkip(t *testing.T) {
	var nilCfg *Config
	assert.False(t, nilCfg.ShouldSkip(testHealthPath))

	cfg := &Config{
		SkipPaths: []string{testHealthPath, testReadyPath, testCustomSkipPath, ""},
	}
	assert.True(t, cfg.ShouldSkip(testHealthPath))
	assert.True(t, cfg.ShouldSkip(testReadyPath))
	assert.True(t, cfg.ShouldSkip(testCustomSkipPath))
	assert.False(t, cfg.ShouldSkip(testOtherPath))
	assert.False(t, cfg.ShouldSkip(""))
}

// TestConfig_IsGRPC tests gRPC protocol detection.
func TestConfig_IsGRPC(t *testing.T) {
	var nilCfg *Config
	assert.False(t, nilCfg.IsGRPC())

	cfg := &Config{Protocol: ProtocolGRPC}
	assert.True(t, cfg.IsGRPC())

	upperCfg := &Config{Protocol: "GRPC"}
	assert.True(t, upperCfg.IsGRPC())

	httpCfg := &Config{Protocol: ProtocolHTTP}
	assert.False(t, httpCfg.IsGRPC())
}

// TestConfig_IsHTTP tests HTTP protocol variant detection.
func TestConfig_IsHTTP(t *testing.T) {
	var nilCfg *Config
	assert.False(t, nilCfg.IsHTTP())

	for _, proto := range []string{ProtocolHTTP, ProtocolHTTPProtobuf, ProtocolHTTPJSON} {
		cfg := &Config{Protocol: proto}
		assert.True(t, cfg.IsHTTP())
	}

	grpcCfg := &Config{Protocol: ProtocolGRPC}
	assert.False(t, grpcCfg.IsHTTP())
}
