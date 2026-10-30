package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func verifyOtelDisabledPrecedence(t *testing.T) {
	t.Helper()
	clearEnv(t)
	t.Setenv(testEnvOtelEnabled, testValTrue)
	t.Setenv(testEnvDisableTraceCorrelation, testValTrue)
	assert.True(t, IsOTelDisabled(nil))
	assert.True(t, IsOTelDisabled(&Config{DisableTraceCorrelation: true}))

	clearEnv(t)
	t.Setenv(testEnvOtelEnabled, testValFalse)
	assert.True(t, IsOTelDisabled(&Config{DisableTraceCorrelation: false}))
}

func verifyStartupAndInvalidPrecedence(t *testing.T) {
	t.Helper()
	setupNoopTracer(t)
	clearEnv(t)
	t.Setenv(testEnvOtelEnabled, testValTrue)
	assert.False(t, IsOTelDisabled(nil), "OTEL_ENABLED=true takes precedence over noop tracer")

	clearEnv(t)
	t.Setenv(testEnvOtelEnabled, "maybe")
	assert.True(t, IsOTelDisabled(nil), "invalid env falls back to noop tracer")

	setupActiveTracer(t)
	assert.False(t, IsOTelDisabled(nil), "invalid env falls back to active tracer")
}

func verifyProgrammaticOverride(t *testing.T) {
	t.Helper()
	setupActiveTracer(t)
	clearEnv(t)
	t.Setenv(testEnvDisableTraceCorrelation, testValTrue)
	assert.False(t, IsOTelDisabled(&Config{DisableTraceCorrelation: false}))

	clearEnv(t)
	t.Setenv(testEnvOtelEnabled, testValTrue)
	assert.True(t, IsOTelDisabled(&Config{DisableTraceCorrelation: true}))
}

func TestIsOTelDisabled_Precedence(t *testing.T) {
	setupActiveTracer(t)
	verifyOtelDisabledPrecedence(t)
	verifyStartupAndInvalidPrecedence(t)
	verifyProgrammaticOverride(t)
}

func runGCPStructuredLoggingCases(t *testing.T) {
	t.Helper()
	t.Run("nil", func(t *testing.T) {
		assert.False(t, IsGCPStructuredLoggingEnabled(nil))
	})
	t.Run("unconfigured", func(t *testing.T) {
		setupNoopTracer(t)
		clearEnv(t)
		cfg := &Config{CloudProvider: testValGCP}
		assert.True(t, IsGCPStructuredLoggingEnabled(cfg))
		assert.True(t, cfg.IsGCPStructuredLoggingEnabled())
	})
	t.Run("active_enabled", func(t *testing.T) {
		setupActiveTracer(t)
		clearEnv(t)
		cfg := &Config{CloudProvider: testValGCP}
		assert.False(t, IsGCPStructuredLoggingEnabled(cfg))
	})
	t.Run("active_disabled_env", func(t *testing.T) {
		setupActiveTracer(t)
		clearEnv(t)
		t.Setenv(testEnvOtelEnabled, testValFalse)
		cfg := &Config{CloudProvider: testValGCP}
		assert.True(t, IsGCPStructuredLoggingEnabled(cfg))
	})
	t.Run("active_disabled_cfg", func(t *testing.T) {
		setupActiveTracer(t)
		clearEnv(t)
		cfg := &Config{CloudProvider: testValGCP, DisableTraceCorrelation: true}
		assert.True(t, IsGCPStructuredLoggingEnabled(cfg))
	})
	t.Run("non_gcp", func(t *testing.T) {
		setupNoopTracer(t)
		clearEnv(t)
		cfg := &Config{CloudProvider: testValAWS}
		assert.False(t, IsGCPStructuredLoggingEnabled(cfg))
	})
}

func TestIsGCPStructuredLoggingEnabled(t *testing.T) {
	runGCPStructuredLoggingCases(t)
}
