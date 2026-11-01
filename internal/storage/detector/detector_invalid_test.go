package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// invalidURLCases returns test cases for malformed URLs, bad ports, and unsupported schemes.
func invalidURLCases() []string {
	return []string{
		"http://bad host.com",
		"http://\x7fbad",
		"http://",
		"http:///s3path",
		"ftp://minio.local",
		"postgres://db.local",
		"http://minio:99999",
		"http://minio:0",
		"http://minio:-1",
		"https:/minio:9000",
		"http:/localhost:9000",
		"s3:/bucket/x",
		"x:/host",
		"/leading/slash",
		"//admin:pass@host",
	}
}

// whitespaceURLCases returns test cases with internal, leading, or trailing whitespace.
func whitespaceURLCases() []string {
	return []string{
		"http://localhost:9000/my path",
		" http://localhost:9000",
		"http://localhost:9000 ",
		"\thttp://localhost:9000",
	}
}

// TestDetectProvider_InvalidURL verifies rejection of malformed URLs and unsupported schemes.
func TestDetectProvider_InvalidURL(t *testing.T) {
	for _, bad := range invalidURLCases() {
		t.Run(bad, func(t *testing.T) {
			p, err := DetectProvider(bad, "")
			require.Error(t, err)
			assert.Empty(t, p)
		})
	}
}

// TestDetectProvider_Whitespace verifies rejection of URLs with whitespace.
func TestDetectProvider_Whitespace(t *testing.T) {
	for _, ws := range whitespaceURLCases() {
		t.Run(ws, func(t *testing.T) {
			p, err := DetectProvider(ws, "")
			require.Error(t, err)
			assert.Contains(t, err.Error(), "contains whitespace")
			assert.Empty(t, p)
		})
	}
}

// TestParseTarget_Empty verifies rejection of empty raw target string.
func TestParseTarget_Empty(t *testing.T) {
	_, err := ParseTarget("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty target")
}

// TestParseTarget_Valid verifies successful parsing for valid target URLs.
func TestParseTarget_Valid(t *testing.T) {
	u, err := ParseTarget("localhost:9000")
	require.NoError(t, err)
	assert.Equal(t, "localhost", u.Hostname())
	assert.Equal(t, "9000", u.Port())
	assert.Equal(t, "https", u.Scheme)
}

// TestParseTarget_RejectsNestedSchemeCredentialsInPath checks ambiguous "a/b://user:pw@host" rejection.
func TestParseTarget_RejectsNestedSchemeCredentialsInPath(t *testing.T) {
	for _, bad := range []string{
		"a/b://user:SECRETPW@host",
		"https://a/b://user:SECRETPW@host",
	} {
		_, err := ParseTarget(bad)
		require.Error(t, err, bad)
		assert.NotContains(t, err.Error(), "SECRETPW")
	}
}
