package provider

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nawaphonOHM/whatever/internal/opentelemetry/config"
)

// buildSampler returns an appropriate Sampler based on the configured sample rate.
func buildSampler(rate float64) sdktrace.Sampler {
	if rate <= config.MinSampleRate {
		return sdktrace.ParentBased(sdktrace.NeverSample())
	}
	if rate >= config.MaxSampleRate {
		return sdktrace.ParentBased(sdktrace.AlwaysSample())
	}
	return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(rate))
}
