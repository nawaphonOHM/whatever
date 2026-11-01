package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const endpointPort9000 = "9000"

// TestParseTarget_PreservesLegitimateEndpoints locks in that hardening does not reject valid IPv6 literals
// or "@" / "%40" characters in URL paths.
func TestParseTarget_PreservesLegitimateEndpoints(t *testing.T) {
	tests := []struct {
		target   string
		wantHost string
		wantPort string
	}{
		{"http://[::1]:9000", "::1", endpointPort9000},
		{"[::1]:9000", "::1", endpointPort9000},
		{"https://[::1]/x", "::1", ""},
		{"http://[::ffff:127.0.0.1]:9000", "::ffff:127.0.0.1", endpointPort9000},
		{"[2001:db8::1]", "2001:db8::1", ""},
		{"https://minio.local/tenant@team", "minio.local", ""},
		{"https://proxy.example.com/s3/a@b/", "proxy.example.com", ""},
		{"https://bucket.s3.amazonaws.com/key@v1", "bucket.s3.amazonaws.com", ""},
		{"https://minio.local/a%40b", "minio.local", ""},
		{"admin:pw@minio.local:9000/tenant@team", "minio.local", endpointPort9000},
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
