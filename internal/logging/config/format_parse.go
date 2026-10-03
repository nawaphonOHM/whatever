package config

import (
	"fmt"
	"strings"
)

func isTextFormat(norm string) bool {
	return norm == "" || norm == string(FormatText)
}

func matchFormat(norm string) (Format, bool) {
	if isTextFormat(norm) {
		return FormatText, true
	}
	if norm == string(FormatJSON) {
		return FormatJSON, true
	}
	return "", false
}

// ParseFormat parses a case-insensitive string into a Format.
func ParseFormat(s string) (Format, error) {
	norm := strings.ToLower(strings.TrimSpace(s))
	if fmtEnum, ok := matchFormat(norm); ok {
		return fmtEnum, nil
	}
	return "", fmt.Errorf("%w: %q (supported: %s, %s)", ErrInvalidFormat, s, FormatText, FormatJSON)
}
