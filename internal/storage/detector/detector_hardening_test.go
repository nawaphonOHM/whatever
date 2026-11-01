package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
