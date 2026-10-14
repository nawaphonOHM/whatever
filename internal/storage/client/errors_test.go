package client_test

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

type dummySmithyError struct {
	code string
}

func (e dummySmithyError) Error() string { return "smithy error: " + e.code }

func (e dummySmithyError) ErrorCode() string { return e.code }

func (e dummySmithyError) ErrorMessage() string { return "message: " + e.code }

func (dummySmithyError) ErrorFault() smithy.ErrorFault { return smithy.FaultClient }

func TestClient_ErrorMapping_NoSuchKey(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		_ = req
		return nil, &types.NoSuchKey{Message: nil}
	})
	c := client.NewClient(client.ClientParams{
		Raw:           buildRawClientForTest(httpClient),
		DefaultBucket: "my-bucket",
	})
	_, err := c.GetObject(t.Context(), contracts.GetObjectInput{Key: "k"})
	assert.ErrorIs(t, err, contracts.ErrObjectNotFound)
}

func TestClient_ErrorMapping_NoSuchBucket(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		_ = req
		return nil, &types.NoSuchBucket{Message: nil}
	})
	c := client.NewClient(client.ClientParams{
		Raw:           buildRawClientForTest(httpClient),
		DefaultBucket: "my-bucket",
	})
	err := c.Ping(t.Context())
	assert.ErrorIs(t, err, contracts.ErrBucketNotFound)
}

func TestClient_ErrorMapping_ResponseError404(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		_ = req
		return newHTTPResponse(http.StatusNotFound, "plain 404 text not found", nil), nil
	})
	c := client.NewClient(client.ClientParams{
		Raw:           buildRawClientForTest(httpClient),
		DefaultBucket: "my-bucket",
	})
	_, err := c.GetObject(t.Context(), contracts.GetObjectInput{Key: "k"})
	assert.ErrorIs(t, err, contracts.ErrObjectNotFound)
}

func TestClient_ErrorMapping_DirectMapError(t *testing.T) {
	t.Parallel()

	t.Run("smithy error with NoSuchKey", func(t *testing.T) {
		t.Parallel()
		mapped := client.MapErrorForTest(dummySmithyError{code: "NoSuchKey"})
		assert.ErrorIs(t, mapped, contracts.ErrObjectNotFound)
	})

	t.Run("smithy error with NoSuchBucket", func(t *testing.T) {
		t.Parallel()
		mapped := client.MapErrorForTest(dummySmithyError{code: "NoSuchBucket"})
		assert.ErrorIs(t, mapped, contracts.ErrBucketNotFound)
	})

	t.Run("nil error returns nil", func(t *testing.T) {
		t.Parallel()
		assert.Nil(t, client.MapErrorForTest(nil))
	})

	t.Run("sentinel error preserved", func(t *testing.T) {
		t.Parallel()
		assert.ErrorIs(t, client.MapErrorForTest(contracts.ErrNilClient), contracts.ErrNilClient)
	})
}
