package client

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeLoggedURL_RedactsCredentialsAndQueries(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "url with user credentials and query",
			input:    "https://user:secretpass@minio.internal:9000/bucket?token=12345#ref",
			expected: "https://minio.internal:9000/bucket",
		},
		{
			name:     "schemeless with user credentials and query",
			input:    "admin:secret@minio.corp:9000/data?key=value",
			expected: "minio.corp:9000/data",
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
			name:     "invalid URL with whitespace returns placeholder",
			input:    "https://invalid url with spaces.com",
			expected: "<unparseable>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, sanitizeLoggedURL(tc.input))
		})
	}
}

func TestFormatLoggedEndpoint(t *testing.T) {
	assert.Empty(t, formatLoggedEndpoint(nil))
	ep := aws.String("https://user:pass@s3.amazonaws.com?secret=val")
	assert.Equal(t, "https://s3.amazonaws.com", formatLoggedEndpoint(ep))
}
