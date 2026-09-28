package provider

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	sampleRateNegative  = -0.5
	sampleRateExcessive = 1.5
	sampleRateHalf      = 0.5
)

func TestBuildSampler(t *testing.T) {
	tests := []struct {
		name        string
		containsStr string
		rate        float64
	}{
		{
			name:        "zero rate never samples",
			containsStr: "AlwaysOffSampler",
			rate:        0.0,
		},
		{
			name:        "negative rate never samples",
			containsStr: "AlwaysOffSampler",
			rate:        sampleRateNegative,
		},
		{
			name:        "one rate always samples",
			containsStr: "AlwaysOnSampler",
			rate:        1.0,
		},
		{
			name:        "greater than one rate always samples",
			containsStr: "AlwaysOnSampler",
			rate:        sampleRateExcessive,
		},
		{
			name:        "ratio sampler description",
			containsStr: "TraceIDRatioBased",
			rate:        sampleRateHalf,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := buildSampler(tt.rate)
			assert.NotNil(t, sampler)
			desc := sampler.Description()
			assert.True(t, strings.Contains(desc, tt.containsStr), "expected %q in %q", tt.containsStr, desc)
			assert.True(t, strings.Contains(desc, "ParentBased"), "expected ParentBased in %q", desc)
		})
	}
}
