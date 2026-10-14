package client_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	streamPayload     = "hello stream"
	streamPayloadSize = int64(len(streamPayload))
)

func putOKResponder(rec recordedRequest) *http.Response {
	_ = rec
	return newHTTPResponse(http.StatusOK, "", map[string]string{"ETag": `"put-etag"`})
}

func nonSeekableBody() io.Reader {
	return io.MultiReader(strings.NewReader("hello "), strings.NewReader("stream"))
}

func putStream(t *testing.T, endpoint string) (recordedRequest, *contracts.PutObjectOutput) {
	t.Helper()
	rec := &recorder{respond: putOKResponder}
	c := connectRecorded(t, endpointConfig(endpoint), rec)
	out, err := c.PutObject(t.Context(), putInput(nonSeekableBody()))
	require.NoError(t, err)
	require.Len(t, rec.all(), 1)
	return rec.all()[0], out
}

func TestPutObject_NonSeekableBody_SignedOverHTTP(t *testing.T) {
	t.Parallel()

	req, out := putStream(t, testHTTPEndpoint)
	assert.Equal(t, "http", req.Scheme)
	assert.Equal(t, http.MethodPut, req.Method)
	assert.Equal(t, testPathStylePath, req.Path)
	assert.Contains(t, req.Header.Get("Authorization"), sigV4Algorithm)
	assert.Equal(t, sha256Hex([]byte(streamPayload)), req.Header.Get("X-Amz-Content-Sha256"))
	assert.Equal(t, streamPayloadSize, req.ContentLength)
	assert.Equal(t, streamPayload, string(req.Body))
	assert.Equal(t, `"put-etag"`, out.ETag)
}

func TestPutObject_NonSeekableBody_SignedOverHTTPS(t *testing.T) {
	t.Parallel()

	req, _ := putStream(t, testHTTPSEndpoint)
	assert.Equal(t, "https", req.Scheme)
	assert.Equal(t, testPathStylePath, req.Path)
	assert.Contains(t, req.Header.Get("Authorization"), sigV4Algorithm)
	assert.Equal(t, streamPayloadSize, req.ContentLength)
	assert.Equal(t, streamPayload, string(req.Body))
	assert.NotContains(t, req.Header.Get("Content-Encoding"), "aws-chunked")
}

func TestPutObject_NonSeekableBody_AWSHTTPS(t *testing.T) {
	t.Parallel()

	req, _ := putStream(t, "")
	assert.Equal(t, "https", req.Scheme)
	assert.Equal(t, testAWSVirtualHost, req.Host)
	assert.Contains(t, req.Header.Get("Authorization"), sigV4Algorithm)
	assert.Contains(t, string(req.Body), streamPayload)
}

func putInput(body io.Reader) contracts.PutObjectInput {
	return contracts.PutObjectInput{Key: testObjectKey, Body: body}
}

func putTwice(t *testing.T, bodies ...io.Reader) []recordedRequest {
	t.Helper()
	rec := &recorder{respond: putOKResponder}
	c := connectRecorded(t, endpointConfig(testHTTPEndpoint), rec)
	for _, body := range bodies {
		_, err := c.PutObject(t.Context(), putInput(body))
		require.NoError(t, err)
	}
	return rec.all()
}

func TestPutObject_SeekableAndNilBody(t *testing.T) {
	t.Parallel()

	reqs := putTwice(t, strings.NewReader("seekable"), nil)
	require.Len(t, reqs, 2)
	assert.Equal(t, "seekable", string(reqs[0].Body))
	assert.Equal(t, int64(0), reqs[1].ContentLength)
}
