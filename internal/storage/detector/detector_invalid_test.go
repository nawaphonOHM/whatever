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

func TestDetector_ErrorsDoNotLeakCredentials(t *testing.T) {
	secretPart := "SENSITIVE_SECRET_123"
	tokenPart := "SENSITIVE_TOKEN_456"

	malformedCases := []struct {
		name   string
		target string
	}{
		{
			name:   "whitespace with credentials and token",
			target: "https://admin:" + secretPart + "@minio .internal:9000/data?token=" + tokenPart,
		},
		{
			name:   "tab with credentials",
			target: "admin:" + secretPart + "@minio\t.internal:9000/data",
		},
		{
			name:   "schemeless with only credentials",
			target: "user:" + secretPart + "@",
		},
		{
			name:   "malformed percent encoding with credentials",
			target: "https://%zz:" + secretPart + "@minio.internal:9000",
		},
		{
			name:   "single slash https with credentials",
			target: "https:/admin:" + secretPart + "@minio.internal:9000",
		},
		{
			name:   "single slash http with credentials and port",
			target: "http:/admin:" + secretPart + "@minio:9000",
		},
		{
			name:   "single slash s3 with credentials and path",
			target: "s3:/user:" + secretPart + "@host/x",
		},
		{
			name:   "single slash custom scheme with credentials",
			target: "x:/user:" + secretPart + "@host",
		},
		{
			name:   "secret in port position with absolute url and credentials",
			target: "https://admin:" + secretPart + "@minio.internal:" + secretPart,
		},
		{
			name:   "secret in port position with absolute url without credentials",
			target: "https://host:" + secretPart,
		},
		{
			name:   "secret in port position with schemeless endpoint",
			target: "minio.internal:" + secretPart,
		},
	}

	for _, tc := range malformedCases {
		t.Run(tc.name, func(t *testing.T) {
			_, parseErr := ParseTarget(tc.target)
			require.Error(t, parseErr)
			assert.NotContains(t, parseErr.Error(), secretPart)
			assert.NotContains(t, parseErr.Error(), tokenPart)

			_, detectErr := Detect(tc.target, "")
			require.Error(t, detectErr)
			assert.NotContains(t, detectErr.Error(), secretPart)
			assert.NotContains(t, detectErr.Error(), tokenPart)

			_, provErr := DetectProvider(tc.target, "")
			require.Error(t, provErr)
			assert.NotContains(t, provErr.Error(), secretPart)
			assert.NotContains(t, provErr.Error(), tokenPart)
		})
	}
}

// TestParseTarget_PreservesLegitimateEndpoints locks in that hardening does not reject valid IPv6 literals
// or "@" / "%40" characters in URL paths.
func TestParseTarget_PreservesLegitimateEndpoints(t *testing.T) {
	tests := []struct {
		target   string
		wantHost string
		wantPort string
	}{
		{"http://[::1]:9000", "::1", "9000"},
		{"[::1]:9000", "::1", "9000"},
		{"https://[::1]/x", "::1", ""},
		{"http://[::ffff:127.0.0.1]:9000", "::ffff:127.0.0.1", "9000"},
		{"[2001:db8::1]", "2001:db8::1", ""},
		{"https://minio.local/tenant@team", "minio.local", ""},
		{"https://proxy.example.com/s3/a@b/", "proxy.example.com", ""},
		{"https://bucket.s3.amazonaws.com/key@v1", "bucket.s3.amazonaws.com", ""},
		{"https://minio.local/a%40b", "minio.local", ""},
		{"admin:pw@minio.local:9000/tenant@team", "minio.local", "9000"},
	}
	for _, tc := range tests {
		t.Run(tc.target, func(t *testing.T) {
			u, err := ParseTarget(tc.target)
			require.NoError(t, err)
			assert.Equal(t, tc.wantHost, u.Hostname())
			assert.Equal(t, tc.wantPort, u.Port())

			d, err := Detect(tc.target, "")
			require.NoError(t, err)
			assert.Equal(t, tc.wantHost, d.Host)
		})
	}
}

// TestParseTarget_RejectsNestedSchemeCredentialsInPath verifies the ambiguous "a/b://user:pw@host" shape stays rejected.
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
