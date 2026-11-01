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
	credentialCases := []string{
		"https://AKIAKEY:supersecret@minio.local:9000 /x",
		"https://AKIAKEY:supersecret@minio.local/%zz",
		"https://AKIAKEY:supersecret@",
		"https://AKIAKEY:supersecret@/s3path",
		"http://AKIAKEY:supersecret@bad host.com",
		"ftp://AKIAKEY:supersecret@minio.local",
		"http://AKIAKEY:supersecret@minio:99999",
		"http://AKIAKEY:supersecret@minio:-1",
		"AKIAKEY:supersecret@minio.local:9000 /x",
	}

	for _, bad := range credentialCases {
		t.Run("detect_endpoint_"+bad, func(t *testing.T) {
			_, err := Detect(bad, "")
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "supersecret")
			assert.NotContains(t, err.Error(), "AKIAKEY")

			_, err = DetectProvider(bad, "")
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "supersecret")
			assert.NotContains(t, err.Error(), "AKIAKEY")
		})

		t.Run("detect_bucketURL_"+bad, func(t *testing.T) {
			_, err := Detect("", bad)
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "supersecret")
			assert.NotContains(t, err.Error(), "AKIAKEY")

			_, err = DetectProvider("", bad)
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "supersecret")
			assert.NotContains(t, err.Error(), "AKIAKEY")
		})

		t.Run("parseTarget_"+bad, func(t *testing.T) {
			_, err := ParseTarget(bad)
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "supersecret")
			assert.NotContains(t, err.Error(), "AKIAKEY")
		})
	}
}
