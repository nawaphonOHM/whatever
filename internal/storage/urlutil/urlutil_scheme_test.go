package urlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidScheme(t *testing.T) {
	tests := []struct {
		scheme string
		valid  bool
	}{
		{"http", true},
		{"https", true},
		{"s3", true},
		{"custom-scheme", true},
		{"scheme+v1.alpha-2", true},
		{"HTTP", true},
		{"", false},
		{"1http", false},
		{"+http", false},
		{"-http", false},
		{".http", false},
		{"http:", false},
		{"a/b", false},
		{"user:pass", false},
		{"user@host", false},
	}

	for _, tt := range tests {
		t.Run(tt.scheme, func(t *testing.T) {
			assert.Equal(t, tt.valid, IsValidScheme(tt.scheme))
		})
	}
}

func TestHasSchemePrefix(t *testing.T) {
	tests := []struct {
		raw   string
		valid bool
	}{
		{"http://localhost:9000", true},
		{"https://storage.googleapis.com", true},
		{"s3://bucket/key", true},
		{"custom+v1://host", true},
		{"localhost:9000", false},
		{"user:pass@host:9000/x?r=a://b", false},
		{"//host:9000/x", false},
		{"a/b://host:9000", false},
		{"://host", false},
		{"https:/host", false},
		{"http:/admin:pass@host:9000", false},
		{"s3:/bucket/x", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			assert.Equal(t, tt.valid, HasSchemePrefix(tt.raw))
		})
	}
}

func TestHasSchemeSingleSlash(t *testing.T) {
	tests := []struct {
		raw      string
		expected bool
	}{
		{"https:/host", true},
		{"http:/admin:pass@minio:9000", true},
		{"s3:/user:pass@host/x", true},
		{"x:/user:pass@host", true},
		{"https://storage.googleapis.com", false},
		{"http://localhost:9000", false},
		{"s3://bucket/key", false},
		{"localhost:9000", false},
		{"user:pass@host:9000/x?r=a://b", false},
		{"//host:9000", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			assert.Equal(t, tt.expected, HasSchemeSingleSlash(tt.raw))
		})
	}
}
