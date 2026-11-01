package urlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testLocalhost9000 = "localhost:9000"
	testGCSURL        = "https://storage.googleapis.com"
	testUnparseable   = "<unparseable>"
)

func TestSanitizeURL_Standard(t *testing.T) {
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
			input:    testGCSURL,
			expected: testGCSURL,
		},
		{
			name:     "clean schemeless URL unchanged",
			input:    testLocalhost9000,
			expected: testLocalhost9000,
		},
		{
			name:     "custom scheme with credentials",
			input:    "custom-scheme://user:pass@service.internal:8080/data?key=val",
			expected: "custom-scheme://service.internal:8080/data",
		},
		{
			name:     "protocol-relative with credentials fails closed to unparseable",
			input:    "//admin:pass@minio.internal:9000/x",
			expected: testUnparseable,
		},
		{
			name:     "empty-authority quadruple slash with credentials fails closed",
			input:    "https:////user:pass@host",
			expected: testUnparseable,
		},
		{
			name:     "malformed path pseudo-scheme with credentials fails closed",
			input:    "a/b://user:pass@host",
			expected: testUnparseable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SanitizeURL(tt.input))
		})
	}
}
