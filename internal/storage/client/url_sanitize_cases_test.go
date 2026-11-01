package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeLoggedURL_MalformedAndEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "protocol-relative with credentials fails closed to unparseable",
			input:    "//admin:pass@minio.internal:9000/x",
			expected: testClientUnparseable,
		},
		{
			name:     "empty-authority quadruple slash with credentials fails closed",
			input:    "https:////user:pass@host",
			expected: testClientUnparseable,
		},
		{
			name:     "malformed path pseudo-scheme with credentials fails closed",
			input:    "a/b://user:pass@host",
			expected: testClientUnparseable,
		},
		{
			name:     "leading slash with credentials fails closed",
			input:    "/user:pass@host/path",
			expected: testClientUnparseable,
		},
		{
			name:     "single slash https with credentials fails closed",
			input:    "https:/user:pass@host",
			expected: testClientUnparseable,
		},
		{
			name:     "single slash http with credentials and port fails closed",
			input:    "http:/admin:pass@minio:9000",
			expected: testClientUnparseable,
		},
		{
			name:     "single slash s3 with credentials and path fails closed",
			input:    "s3:/user:pass@host/x",
			expected: testClientUnparseable,
		},
		{
			name:     "single slash custom scheme with credentials fails closed",
			input:    "x:/user:pass@host",
			expected: testClientUnparseable,
		},
		{
			name:     "single slash https without credentials fails closed",
			input:    "https:/host",
			expected: testClientUnparseable,
		},
		{
			name:     "single slash http without credentials fails closed",
			input:    "http:/localhost:9000",
			expected: testClientUnparseable,
		},
		{
			name:     "schemeless with only user credentials fails closed",
			input:    "user:pass@",
			expected: testClientUnparseable,
		},
		{
			name:     "invalid URL with whitespace returns placeholder",
			input:    "https://invalid url with spaces.com",
			expected: testClientUnparseable,
		},
		{
			name:     "invalid URL with tab character returns placeholder",
			input:    "https://invalid\turl.com",
			expected: testClientUnparseable,
		},
		{
			name:     "invalid URL with newline returns placeholder",
			input:    "https://invalid\nurl.com",
			expected: testClientUnparseable,
		},
		{
			name:     "malformed percent encoding returns placeholder",
			input:    "https://%zz@invalid.com",
			expected: testClientUnparseable,
		},
		{
			name:     "malformed schemeless percent encoding returns placeholder",
			input:    "%zz@invalid.com",
			expected: testClientUnparseable,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, sanitizeLoggedURL(tc.input))
		})
	}
}
