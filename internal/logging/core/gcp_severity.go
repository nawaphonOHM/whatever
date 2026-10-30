package core

import (
	"log/slog"
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
)

// Google Cloud Logging standard severity level constants.
const (
	GCPSeverityDefault   = "DEFAULT"
	GCPSeverityDebug     = "DEBUG"
	GCPSeverityInfo      = "INFO"
	GCPSeverityNotice    = "NOTICE"
	GCPSeverityWarning   = "WARNING"
	GCPSeverityError     = "ERROR"
	GCPSeverityCritical  = "CRITICAL"
	GCPSeverityAlert     = "ALERT"
	GCPSeverityEmergency = "EMERGENCY"
)

func mapUpperSeverity(l slog.Level) string {
	if l < config.SlogLevelError {
		return GCPSeverityWarning
	}
	if l < config.SlogLevelFatal {
		return GCPSeverityError
	}
	return GCPSeverityCritical
}

// SeverityFromLevel maps a slog.Level to its standard GCP severity string.
func SeverityFromLevel(l slog.Level) string {
	if l < config.SlogLevelInfo {
		return GCPSeverityDebug
	}
	if l < config.SlogLevelWarn {
		return GCPSeverityInfo
	}
	return mapUpperSeverity(l)
}

func matchAlertEmergency(s string) (string, bool) {
	switch s {
	case "alert":
		return GCPSeverityAlert, true
	case "emergency", "emerg":
		return GCPSeverityEmergency, true
	}
	return "", false
}

func matchSpecialSeverity(s string) (string, bool) {
	if s == "notice" {
		return GCPSeverityNotice, true
	}
	return matchAlertEmergency(s)
}

func parseStandardSeverity(norm string) string {
	lvl, err := config.ParseSlogLevel(norm)
	if err != nil {
		return GCPSeverityDefault
	}
	return SeverityFromLevel(lvl)
}

// ParseGCPSeverity maps a string representation to its standard GCP severity string.
func ParseGCPSeverity(s string) string {
	norm := strings.ToLower(strings.TrimSpace(s))
	if norm == "" {
		return GCPSeverityDefault
	}
	if sev, ok := matchSpecialSeverity(norm); ok {
		return sev
	}
	return parseStandardSeverity(norm)
}
