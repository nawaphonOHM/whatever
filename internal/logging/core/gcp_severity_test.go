package core_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

func TestSeverityFromLevel(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		level    slog.Level
	}{
		{"trace", core.GCPSeverityDebug, config.SlogLevelTrace},
		{"debug", core.GCPSeverityDebug, config.SlogLevelDebug},
		{"info", core.GCPSeverityInfo, config.SlogLevelInfo},
		{"warn", core.GCPSeverityWarning, config.SlogLevelWarn},
		{"error", core.GCPSeverityError, config.SlogLevelError},
		{"fatal", core.GCPSeverityCritical, config.SlogLevelFatal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, core.SeverityFromLevel(tt.level))
		})
	}
}

func TestParseGCPSeveritySpecial(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"notice", core.GCPSeverityNotice},
		{"NOTICE", core.GCPSeverityNotice},
		{"alert", core.GCPSeverityAlert},
		{"emergency", core.GCPSeverityEmergency},
		{"emerg", core.GCPSeverityEmergency},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, core.ParseGCPSeverity(tt.input))
		})
	}
}

func TestParseGCPSeverityStandard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"debug", core.GCPSeverityDebug},
		{"info", core.GCPSeverityInfo},
		{"warn", core.GCPSeverityWarning},
		{"error", core.GCPSeverityError},
		{"fatal", core.GCPSeverityCritical},
		{"", core.GCPSeverityDefault},
		{"invalid", core.GCPSeverityDefault},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, core.ParseGCPSeverity(tt.input))
		})
	}
}
