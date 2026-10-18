package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func newTestPresignClient() *client.Client {
	raw := s3.New(s3.Options{
		Region: "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider(
			"AKIAIOSFODNN7EXAMPLE",
			"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			"",
		),
	})
	presign := s3.NewPresignClient(raw)
	return client.NewClient(client.ClientParams{
		Raw:           raw,
		Presign:       presign,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
		PresignExpiry: testPresignExpiry15m,
	})
}

func TestClient_PresignGetObject(t *testing.T) {
	t.Parallel()

	c := newTestPresignClient()
	out, err := c.PresignGetObject(context.Background(), contracts.PresignGetObjectInput{
		Key:       "download.pdf",
		ExpiresIn: testPresignExpiry5m,
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, http.MethodGet, out.Method)
	assert.Contains(t, out.URL, "my-bucket")
	assert.Contains(t, out.URL, "download.pdf")
	assert.Contains(t, out.URL, "X-Amz-Signature")
	assert.Equal(t, testPresignExpiry5m, out.ExpiresIn)
}

func TestClient_PresignPutObject(t *testing.T) {
	t.Parallel()

	c := newTestPresignClient()
	out, err := c.PresignPutObject(context.Background(), contracts.PresignPutObjectInput{
		Key:         "upload.png",
		ContentType: "image/png",
		ExpiresIn:   testPresignExpiry20m,
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, http.MethodPut, out.Method)
	assert.Contains(t, out.URL, "upload.png")
	assert.Contains(t, out.URL, "X-Amz-Signature")
	assert.Equal(t, testPresignExpiry20m, out.ExpiresIn)
}

func TestClient_Presign_NilClient(t *testing.T) {
	t.Parallel()

	var c *client.Client
	getOut, getErr := c.PresignGetObject(context.Background(), contracts.PresignGetObjectInput{Key: "k"})
	assert.Nil(t, getOut)
	assert.ErrorIs(t, getErr, contracts.ErrNilClient)

	putOut, putErr := c.PresignPutObject(context.Background(), contracts.PresignPutObjectInput{Key: "k"})
	assert.Nil(t, putOut)
	assert.ErrorIs(t, putErr, contracts.ErrNilClient)
}
