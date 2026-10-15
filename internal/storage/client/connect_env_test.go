package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
)

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
