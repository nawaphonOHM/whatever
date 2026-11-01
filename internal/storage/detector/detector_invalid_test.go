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

// TestDetect_CredentialRedaction verifies that credentials in malformed URLs are redacted across detection functions.
func TestDetect_CredentialRedaction(t *testing.T) {
	credentialCases := []struct {
		name      string
		target    string
		forbidden []string
	}{
		{
			name:      "alphanumeric credentials with whitespace",
			target:    "https://AKIAKEY:supersecret@minio.local:9000 /x",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "alphanumeric credentials with bad escape",
			target:    "https://AKIAKEY:supersecret@minio.local/%zz",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "alphanumeric credentials with missing host",
			target:    "https://AKIAKEY:supersecret@",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "alphanumeric credentials with missing host and path",
			target:    "https://AKIAKEY:supersecret@/s3path",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "alphanumeric credentials with unsupported scheme",
			target:    "ftp://AKIAKEY:supersecret@minio.local",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "alphanumeric credentials with out-of-range port",
			target:    "http://AKIAKEY:supersecret@minio:99999",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "alphanumeric credentials with negative port",
			target:    "http://AKIAKEY:supersecret@minio:-1",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "schemeless credentials with whitespace",
			target:    "AKIAKEY:supersecret@minio.local:9000 /x",
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
		{
			name:      "AWS-style secret with slashes and whitespace",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@minio.local:9000 /x",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
		{
			name:      "AWS-style secret with slashes and bad escape",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@minio.local/%zz",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
		{
			name:      "AWS-style secret with slashes and missing host",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
		{
			name:      "AWS-style secret with slashes and missing host with path",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@/s3path",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
		{
			name:      "AWS-style secret with slashes and trailing at without host",
			target:    "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG@",
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG"},
		},
		{
			name:      "secret containing hash with trailing at",
			target:    "https://AKIAKEY:hash#secret@",
			forbidden: []string{"AKIAKEY", "hash#secret", "hash", "secret"},
		},
		{
			name:      "secret containing hash with whitespace",
			target:    "https://AKIAKEY:hash#secret@minio.local:9000 /x",
			forbidden: []string{"AKIAKEY", "hash#secret", "hash", "secret"},
		},
		{
			name:      "secret containing question mark with trailing at",
			target:    "https://AKIAKEY:quest?secret@",
			forbidden: []string{"AKIAKEY", "quest?secret", "quest", "secret"},
		},
		{
			name:      "secret containing question mark with bad escape",
			target:    "https://AKIAKEY:quest?secret@minio.local/%zz",
			forbidden: []string{"AKIAKEY", "quest?secret", "quest", "secret"},
		},
		{
			name:      "schemeless with credentials and query scheme delimiter with bad port",
			target:    "AKIAKEY:secret@minio:99999/p?x=://",
			forbidden: []string{"AKIAKEY", "secret"},
		},
	}

	for _, tc := range credentialCases {
		t.Run("detect_endpoint_"+tc.name, func(t *testing.T) {
			_, err := Detect(tc.target, "")
			require.Error(t, err)
			for _, secret := range tc.forbidden {
				assert.NotContains(t, err.Error(), secret)
			}

			_, err = DetectProvider(tc.target, "")
			require.Error(t, err)
			for _, secret := range tc.forbidden {
				assert.NotContains(t, err.Error(), secret)
			}
		})

		t.Run("detect_bucketURL_"+tc.name, func(t *testing.T) {
			_, err := Detect("", tc.target)
			require.Error(t, err)
			for _, secret := range tc.forbidden {
				assert.NotContains(t, err.Error(), secret)
			}

			_, err = DetectProvider("", tc.target)
			require.Error(t, err)
			for _, secret := range tc.forbidden {
				assert.NotContains(t, err.Error(), secret)
			}
		})

		t.Run("parseTarget_"+tc.name, func(t *testing.T) {
			_, err := ParseTarget(tc.target)
			require.Error(t, err)
			for _, secret := range tc.forbidden {
				assert.NotContains(t, err.Error(), secret)
			}
		})
	}
}
