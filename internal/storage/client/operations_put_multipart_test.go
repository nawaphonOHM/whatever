package client_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	partSize        = 8 * 1024 * 1024
	initiateXML     = `<InitiateMultipartUploadResult><UploadId>up-1</UploadId></InitiateMultipartUploadResult>`
	completeXMLBody = `<CompleteMultipartUploadResult><ETag>"final-etag"</ETag></CompleteMultipartUploadResult>`
	partNumberParam = "partNumber"

	multipartRequestCount = 4
	completeRequestIndex  = 3
)

func putZeros(t *testing.T, rec *recorder, size int64) (*contracts.PutObjectOutput, error) {
	t.Helper()
	c := connectRecorded(t, endpointConfig(testHTTPEndpoint), rec)
	return c.PutObject(t.Context(), contracts.PutObjectInput{
		Key: testObjectKey, ContentType: "application/octet-stream", Body: &zeroReader{left: size},
	})
}

func TestPutObject_StreamAtPartSizeIsSingleRequest(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: multipartResponder("0")}
	_, err := putZeros(t, rec, partSize)
	require.NoError(t, err)
	reqs := rec.all()
	require.Len(t, reqs, 1)
	assert.Equal(t, int64(partSize), reqs[0].BodyLen)
}

func TestPutObject_StreamAbovePartSizeUsesMultipart(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: multipartResponder("0")}
	out, err := putZeros(t, rec, partSize+1)
	require.NoError(t, err)
	assert.Equal(t, `"final-etag"`, out.ETag)
	assertMultipartRequests(t, rec.all())
}

func assertMultipartRequests(t *testing.T, reqs []recordedRequest) {
	t.Helper()
	require.Len(t, reqs, multipartRequestCount)
	assert.True(t, reqs[0].Query.Has("uploads"))
	assert.Equal(t, "application/octet-stream", reqs[0].Header.Get("Content-Type"))
	assert.Equal(t, int64(partSize), reqs[1].BodyLen)
	assert.Equal(t, int64(1), reqs[2].BodyLen)
	assert.Equal(t, "up-1", reqs[completeRequestIndex].Query.Get("uploadId"))
	assert.Contains(t, string(reqs[completeRequestIndex].Body), "<PartNumber>2</PartNumber>")
}

func TestPutObject_MultipartFailureAbortsUpload(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: multipartResponder("2")}
	out, err := putZeros(t, rec, partSize+1)
	require.Error(t, err)
	assert.Nil(t, out)
	reqs := rec.all()
	last := reqs[len(reqs)-1]
	assert.Equal(t, http.MethodDelete, last.Method)
	assert.Equal(t, "up-1", last.Query.Get("uploadId"))
}

func TestPutObject_BodyReadErrorSendsNothing(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: statusResponder(http.StatusOK)}
	c := connectRecorded(t, endpointConfig(testHTTPEndpoint), rec)
	boom := errors.New("boom")
	body := io.MultiReader(strings.NewReader("x"), &errReader{err: boom})
	_, err := c.PutObject(t.Context(), contracts.PutObjectInput{Key: testObjectKey, Body: body})
	require.ErrorIs(t, err, boom)
	assert.Empty(t, rec.all())
}
