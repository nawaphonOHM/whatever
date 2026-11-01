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
			name:     "regression: schemeless with credentials and scheme in query parameter",
			input:    "user:pass@host:9000/x?r=a://b",
			expected: "host:9000/x",
		},
		{
			name:     "regression: schemeless with credentials and scheme in fragment",
			input:    "user:pass@host:9000/x#frag://b",
			expected: "host:9000/x",
		},
		{
			name:     "regression: absolute url with credentials and scheme in query parameter",
			input:    "https://user:pass@host:9000/x?r=a://b",
			expected: "https://host:9000/x",
		},
		{
			name:     "schemeless redirect with nested credentials in query",
			input:    "user:pass@localhost:9000/redirect?url=http://user2:pass2@other.com",
			expected: "localhost:9000/redirect",
		},
		{
			name:     "absolute redirect with nested credentials in query",
			input:    "https://user:pass@localhost:9000/redirect?url=http://user2:pass2@other.com",
			expected: "https://localhost:9000/redirect",
		},
		{
			name:     "ipv6 schemeless with credentials and scheme in query",
			input:    "user:pass@[::1]:9000/x?r=a://b",
			expected: "[::1]:9000/x",
		},
		{
			name:     "ipv6 absolute with credentials and scheme in query",
			input:    "http://user:pass@[::1]:9000/x?r=a://b",
			expected: "http://[::1]:9000/x",
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
			name:     "s3 scheme with credentials and query",
			input:    "s3://user:secret@my-bucket/path/to/key?versionId=123#frag",
			expected: "s3://my-bucket/path/to/key",
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
			name:     "schemeless with only user credentials fails closed",
			input:    "user:pass@",
			expected: "<unparseable>",
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
			name:     "invalid URL with newline returns placeholder",
			input:    "https://invalid\nurl.com",
			expected: "<unparseable>",
		},
		{
			name:     "malformed percent encoding returns placeholder",
			input:    "https://%zz@invalid.com",
			expected: "<unparseable>",
		},
		{
			name:     "malformed schemeless percent encoding returns placeholder",
			input:    "%zz@invalid.com",
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
