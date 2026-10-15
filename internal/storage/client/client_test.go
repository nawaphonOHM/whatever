package client_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	testPresignExpiry    = 10 * time.Minute
	updatedPresignExpiry = 30 * time.Minute
	nilPresignExpiry     = 5 * time.Minute
)

func assertClientAccessors(t *testing.T, c *client.Client, raw *s3.Client, presign *s3.PresignClient) {
	t.Helper()
	assert.Equal(t, contracts.ProviderAWS, c.Provider())
	assert.Equal(t, "my-test-bucket", c.DefaultBucket())
	assert.Same(t, raw, c.RawClient())
	assert.Same(t, presign, c.PresignClient())
}

func TestNewClient_Accessors(t *testing.T) {
	t.Parallel()

	raw := s3.New(s3.Options{Region: "us-east-1"})
	presign := s3.NewPresignClient(raw)
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Presign:       presign,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-test-bucket",
		PresignExpiry: testPresignExpiry,
	})

	require.NotNil(t, c)
	assertClientAccessors(t, c, raw, presign)
	c.SetPresignExpiry(updatedPresignExpiry)
	c.SetPresignExpiry(-1)
}

func TestClient_NilAccessors(t *testing.T) {
	t.Parallel()

	var c *client.Client
	assert.Empty(t, c.Provider())
	assert.Empty(t, c.DefaultBucket())
	assert.Nil(t, c.RawClient())
	assert.Nil(t, c.PresignClient())
	c.SetPresignExpiry(nilPresignExpiry)
}
