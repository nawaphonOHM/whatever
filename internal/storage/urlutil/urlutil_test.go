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

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string returns empty",
			input:    "",
			expected: "",
		},
		{
			name:     "standard absolute URL with userinfo and query",
			input:    "http://user:secret@localhost:9000/path?token=xyz#frag",
			expected: "http://localhost:9000/path",
		},
		{
			name:     "standard schemeless URL with userinfo and query",
			input:    "user:pass@minio.corp:9000/data?key=val",
			expected: "minio.corp:9000/data",
		},
		{
			name:     "schemeless with credentials and scheme in query parameter",
			input:    "user:pass@host:9000/x?r=a://b",
			expected: "host:9000/x",
		},
		{
			name:     "schemeless with credentials and scheme in fragment",
			input:    "user:pass@host:9000/x#frag://b",
			expected: "host:9000/x",
		},
		{
			name:     "absolute url with credentials and scheme in query parameter",
			input:    "https://user:pass@host:9000/x?r=a://b",
			expected: "https://host:9000/x",
		},
		{
			name:     "schemeless redirect with nested credentials in query",
			input:    "user:pass@localhost:9000/redirect?url=http://user2:pass2@other.com",
			expected: "localhost:9000/redirect",
		},
		{
			name:     "clean absolute URL unchanged",
			input:    "https://storage.googleapis.com",
			expected: "https://storage.googleapis.com",
		},
		{
			name:     "clean schemeless URL unchanged",
			input:    "localhost:9000",
			expected: "localhost:9000",
		},
		{
			name:     "custom scheme with credentials",
			input:    "custom-scheme://user:pass@service.internal:8080/data?key=val",
			expected: "custom-scheme://service.internal:8080/data",
		},
		{
			name:     "protocol-relative with credentials fails closed to unparseable",
			input:    "//admin:pass@minio.internal:9000/x",
			expected: "<unparseable>",
		},
		{
			name:     "empty-authority quadruple slash with credentials fails closed",
			input:    "https:////user:pass@host",
			expected: "<unparseable>",
		},
		{
			name:     "malformed path pseudo-scheme with credentials fails closed",
			input:    "a/b://user:pass@host",
			expected: "<unparseable>",
		},
		{
			name:     "leading slash with credentials fails closed",
			input:    "/user:pass@host/path",
			expected: "<unparseable>",
		},
		{
			name:     "single slash https with credentials fails closed",
			input:    "https:/user:pass@host",
			expected: "<unparseable>",
		},
		{
			name:     "single slash http with credentials and port fails closed",
			input:    "http:/admin:pass@minio:9000",
			expected: "<unparseable>",
		},
		{
			name:     "single slash s3 with credentials and path fails closed",
			input:    "s3:/user:pass@host/x",
			expected: "<unparseable>",
		},
		{
			name:     "single slash custom scheme with credentials fails closed",
			input:    "x:/user:pass@host",
			expected: "<unparseable>",
		},
		{
			name:     "single slash https without credentials fails closed",
			input:    "https:/host",
			expected: "<unparseable>",
		},
		{
			name:     "single slash http without credentials fails closed",
			input:    "http:/localhost:9000",
			expected: "<unparseable>",
		},
		{
			name:     "scheme without slash with credentials redacts credentials",
			input:    "https:user:pass@host",
			expected: "host",
		},
		{
			name:     "userinfo with multiple @ redacts credentials",
			input:    "user:pass@host@evil",
			expected: "evil",
		},
		{
			name:     "leading colon with credentials redacts credentials",
			input:    ":user:pass@host",
			expected: "host",
		},
		{
			name:     "colon without username with credentials",
			input:    "https://:pass@host",
			expected: "https://host",
		},
		{
			name:     "invalid URL with whitespace returns placeholder",
			input:    "https://invalid url with spaces.com",
			expected: "<unparseable>",
		},
		{
			name:     "invalid URL with tab character returns placeholder",
			input:    "https://invalid\turl.com",
			expected: "<unparseable>",
		},
		{
			name:     "malformed percent encoding returns placeholder",
			input:    "https://%zz@invalid.com",
			expected: "<unparseable>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SanitizeURL(tt.input))
		})
	}
}
