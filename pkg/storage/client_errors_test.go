package storage_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const (
	testNotFoundXML = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
    <Code>NoSuchKey</Code>
    <Message>The specified key does not exist.</Message>
    <Key>test-key</Key>
    <BucketName>test-bucket</BucketName>
</Error>`
)

func notFoundHTTPClient() *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     http.StatusText(http.StatusNotFound),
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader([]byte(testNotFoundXML))),
			}, nil
		}),
	}
}

func testStaticConfig() *storage.Config {
	cfg := storage.DefaultConfig()
	cfg.Bucket = "test-bucket"
	cfg.AccessKeyID = "AKIAIOSFODNN7EXAMPLE"
	cfg.SecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	return cfg
}

func TestPublicStorage_ErrorPropagation_NotFound(t *testing.T) {
	t.Parallel()

	c, err := storage.ConnectWithConfig(
		context.Background(),
		testStaticConfig(),
		storage.WithHTTPClient(notFoundHTTPClient()),
		storage.WithSkipProbe(true),
	)
	require.NoError(t, err)
	require.NotNil(t, c)

	out, err := c.GetObject(context.Background(), storage.GetObjectInput{
		Bucket: "test-bucket",
		Key:    "test-key",
	})
	assert.Nil(t, out)
	assert.ErrorIs(t, err, storage.ErrObjectNotFound)
}
