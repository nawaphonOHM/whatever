package client_test

import (
	"context"
	"net/http"
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

func setConnectEnv(t *testing.T) {
	t.Helper()
	t.Setenv("OHM9996_STORAGE_BUCKET", "env-bucket")
	t.Setenv("OHM9996_STORAGE_REGION", "us-east-1")
	t.Setenv("OHM9996_STORAGE_ENDPOINT", "https://s3.amazonaws.com")
	setAWSEnvCredentials(t)
}

func TestConnect_Environment(t *testing.T) {
	setConnectEnv(t)

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		_ = req
		return newHTTPResponse(http.StatusOK, "", nil), nil
	})

	c, err := client.Connect(context.Background(), client.WithHTTPClient(httpClient))
	require.NoError(t, err)
	require.NotNil(t, c)
	assert.Equal(t, "env-bucket", c.DefaultBucket())
}

func TestSetExitFunc(t *testing.T) {
	exitCode := 0
	prev := client.SetExitFunc(func(code int) {
		exitCode = code
	})
	defer func() {
		_ = client.SetExitFunc(prev)
	}()

	assert.NotNil(t, prev)
	assert.Equal(t, 0, exitCode)
}
