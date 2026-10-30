package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestIsOTelDisabled_UnconfiguredTracer(t *testing.T) {
	setupNoopTracer(t)
	clearEnv(t)

	assert.True(t, IsOTelDisabled(nil))
	assert.True(t, IsOTelDisabled(&Config{DisableTraceCorrelation: false}))
	assert.True(t, IsOTelDisabled(&Config{DisableTraceCorrelation: true}))
	assert.True(t, (&Config{}).IsOTelDisabled())
}

func TestIsOTelDisabled_ActiveTracer(t *testing.T) {
	setupActiveTracer(t)
	clearEnv(t)

	assert.False(t, IsOTelDisabled(nil))
	assert.False(t, IsOTelDisabled(&Config{DisableTraceCorrelation: false}))
	assert.False(t, (&Config{DisableTraceCorrelation: false}).IsOTelDisabled())
	assert.True(t, IsOTelDisabled(&Config{DisableTraceCorrelation: true}))
}

func TestIsOTelDisabled_Samplers(t *testing.T) {
	samplers := []sdktrace.Sampler{
		sdktrace.NeverSample(),
		sdktrace.TraceIDRatioBased(testSampleRateLow),
		sdktrace.ParentBased(sdktrace.TraceIDRatioBased(testSampleRateRatio)),
	}
	for _, sampler := range samplers {
		setupActiveTracerWithSampler(t, sampler)
		clearEnv(t)
		assert.False(t, IsOTelDisabled(nil), "active tracer with %s should be enabled", sampler.Description())
	}
}

func checkOtelDisabledForEnvList(t *testing.T, envKey string, values []string, expectDisabled bool) {
	t.Helper()
	for _, val := range values {
		t.Run(envKey+"_"+val, func(t *testing.T) {
			clearEnv(t)
			t.Setenv(envKey, val)
			assert.Equal(t, expectDisabled, IsOTelDisabled(nil))
			assert.Equal(t, expectDisabled, IsOTelDisabled(&Config{DisableTraceCorrelation: false}))
		})
	}
}

func TestIsOTelDisabled_EnvOtelEnabled(t *testing.T) {
	setupActiveTracer(t)
	falsy := []string{"false", "False", "FALSE", "0", "f", "F", " false "}
	checkOtelDisabledForEnvList(t, testEnvOtelEnabled, falsy, true)

	truthy := []string{"true", "True", "TRUE", "1", "t", "T", " true "}
	checkOtelDisabledForEnvList(t, testEnvOtelEnabled, truthy, false)
}

func TestIsOTelDisabled_EnvDisableTraceCorrelation(t *testing.T) {
	setupActiveTracer(t)
	truthy := []string{"true", "True", "TRUE", "1", "t", "T", " true "}
	for _, val := range truthy {
		t.Run("truthy_"+val, func(t *testing.T) {
			clearEnv(t)
			t.Setenv(testEnvDisableTraceCorrelation, val)
			assert.True(t, IsOTelDisabled(nil))
		})
	}

	falsy := []string{"false", "False", "FALSE", "0", "f", "F", " false "}
	for _, val := range falsy {
		t.Run("falsy_"+val, func(t *testing.T) {
			clearEnv(t)
			t.Setenv(testEnvDisableTraceCorrelation, val)
			assert.False(t, IsOTelDisabled(nil))
		})
	}
}
