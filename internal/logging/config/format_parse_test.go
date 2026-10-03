package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFormat_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Format
	}{
		{"text lower", "text", FormatText},
		{"text upper", "TEXT", FormatText},
		{"text spaced", "  text  ", FormatText},
		{"json lower", testJSON, FormatJSON},
		{"json upper", "JSON", FormatJSON},
		{"json spaced", "  json  ", FormatJSON},
		{"empty defaults to text", "", FormatText},
		{"whitespace defaults to text", "   ", FormatText},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmtEnum, err := ParseFormat(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, fmtEnum)
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	invalidInputs := []string{
		"yaml",
		"xml",
		"csv",
		"html",
	}

	for _, in := range invalidInputs {
		t.Run(in, func(t *testing.T) {
			fmtEnum, err := ParseFormat(in)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidFormat)
			assert.Empty(t, fmtEnum)
		})
	}
}
