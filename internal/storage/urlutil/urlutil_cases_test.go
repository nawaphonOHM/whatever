package urlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeURL_MalformedAndInvalid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "leading slash with credentials fails closed",
			input:    "/user:pass@host/path",
			expected: testUnparseable,
		},
		{
			name:     "single slash https with credentials fails closed",
			input:    "https:/user:pass@host",
			expected: testUnparseable,
		},
		{
			name:     "single slash http with credentials and port fails closed",
			input:    "http:/admin:pass@minio:9000",
			expected: testUnparseable,
		},
		{
			name:     "single slash s3 with credentials and path fails closed",
			input:    "s3:/user:pass@host/x",
			expected: testUnparseable,
		},
		{
			name:     "single slash custom scheme with credentials fails closed",
			input:    "x:/user:pass@host",
			expected: testUnparseable,
		},
		{
			name:     "single slash https without credentials fails closed",
			input:    "https:/host",
			expected: testUnparseable,
		},
		{
			name:     "single slash http without credentials fails closed",
			input:    "http:/localhost:9000",
			expected: testUnparseable,
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
			expected: testUnparseable,
		},
		{
			name:     "invalid URL with tab character returns placeholder",
			input:    "https://invalid\turl.com",
			expected: testUnparseable,
		},
		{
			name:     "malformed percent encoding returns placeholder",
			input:    "https://%zz@invalid.com",
			expected: testUnparseable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SanitizeURL(tt.input))
		})
	}
}
