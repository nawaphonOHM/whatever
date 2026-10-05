package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildClientOptions_Firestore tests that Firestore configurations generate
// appropriate driver ClientOptions with TLS, load balancing, and SCRAM-SHA-256.
func TestBuildClientOptions_Firestore(t *testing.T) {
	t.Run("BuildClientOptions applies firestore defaults and settings", func(t *testing.T) {
		cfg := createFirestoreTestConfig()
		clientOpts := BuildClientOptions(cfg)
		require.NotNil(t, clientOpts)

		expectedURI := "mongodb://unreach_user:unreach_pass@test-project.firestore.goog:443/" +
			"?uuidRepresentation=unspecified&tls=true&loadBalanced=true&retryWrites=false" +
			"&authMechanism=SCRAM-SHA-256"
		assert.Equal(t, expectedURI, clientOpts.GetURI())

		require.NotNil(t, clientOpts.LoadBalanced)
		assert.True(t, *clientOpts.LoadBalanced)

		require.NotNil(t, clientOpts.RetryWrites)
		assert.False(t, *clientOpts.RetryWrites)

		assert.NotNil(t, clientOpts.TLSConfig)

		require.NotNil(t, clientOpts.Auth)
		assert.Equal(t, "SCRAM-SHA-256", clientOpts.Auth.AuthMechanism)
		assert.Equal(t, "unreach_user", clientOpts.Auth.Username)
	})

	t.Run("BuildClientOptionsWithTLS enforces TLS=true even when enableTLS is false for firestore", func(t *testing.T) {
		cfg := createFirestoreTestConfig()
		clientOpts := BuildClientOptionsWithTLS(cfg, false)
		require.NotNil(t, clientOpts)

		expectedURI := "mongodb://unreach_user:unreach_pass@test-project.firestore.goog:443/" +
			"?uuidRepresentation=unspecified&tls=true&loadBalanced=true&retryWrites=false" +
			"&authMechanism=SCRAM-SHA-256"
		assert.Equal(t, expectedURI, clientOpts.GetURI())

		require.NotNil(t, clientOpts.LoadBalanced)
		assert.True(t, *clientOpts.LoadBalanced)

		require.NotNil(t, clientOpts.RetryWrites)
		assert.False(t, *clientOpts.RetryWrites)

		assert.NotNil(t, clientOpts.TLSConfig)
	})
}
