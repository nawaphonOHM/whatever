package client_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	crcPart1 = "crc32-part-1"
	crcPart2 = "crc32-part-2"
)

func awsPartReply(partNumber string) *http.Response {
	if partNumber == "1" {
		return newHTTPResponse(http.StatusOK, "", map[string]string{
			"ETag": `"etag-1"`, "x-amz-checksum-crc32": crcPart1,
		})
	}
	return newHTTPResponse(http.StatusOK, "", map[string]string{
		"ETag": `"etag-2"`, "x-amz-checksum-crc32": crcPart2,
	})
}

func awsControlReply(rec recordedRequest) *http.Response {
	if rec.Method != http.MethodPost {
		return newHTTPResponse(http.StatusNoContent, "", nil)
	}
	if rec.Query.Has("uploads") {
		return newHTTPResponse(http.StatusOK, initiateXML, nil)
	}
	return newHTTPResponse(http.StatusOK, completeXMLBody, nil)
}

func awsChecksumResponder() func(recordedRequest) *http.Response {
	return func(rec recordedRequest) *http.Response {
		if rec.Method == http.MethodPut {
			return awsPartReply(rec.Query.Get(partNumberParam))
		}
		return awsControlReply(rec)
	}
}

func nonAWSChecksumResponder() func(recordedRequest) *http.Response {
	return func(rec recordedRequest) *http.Response {
		if rec.Method == http.MethodPut {
			return newHTTPResponse(http.StatusOK, "", map[string]string{
				"ETag": `"etag-` + rec.Query.Get(partNumberParam) + `"`,
			})
		}
		return awsControlReply(rec)
	}
}

func connectAWSDefaultRecorded(t *testing.T, rec *recorder) *client.Client {
	t.Helper()
	setAWSEnvCredentials(t)
	cfg := config.DefaultConfig()
	cfg.Bucket = testBucket
	cfg.Region = "us-east-1"
	c, err := client.ConnectWithConfig(t.Context(), cfg, connectOptions(rec.httpClient())...)
	require.NoError(t, err)
	return c
}

func connectWithConfigRecorded(t *testing.T, cfg *config.Config, rec *recorder) *client.Client {
	t.Helper()
	c, err := client.ConnectWithConfig(t.Context(), cfg, connectOptions(rec.httpClient())...)
	require.NoError(t, err)
	return c
}

func assertCompleteChecksums(t *testing.T, req recordedRequest) {
	t.Helper()
	body := string(req.Body)
	assert.Contains(t, body, "<ChecksumCRC32>"+crcPart1+"</ChecksumCRC32>")
	assert.Contains(t, body, "<ChecksumCRC32>"+crcPart2+"</ChecksumCRC32>")
	assert.Contains(t, body, "<PartNumber>1</PartNumber>")
	assert.Contains(t, body, "<PartNumber>2</PartNumber>")
}

func assertNoMultipartChecksums(t *testing.T, reqs []recordedRequest) {
	t.Helper()
	require.Len(t, reqs, multipartRequestCount)
	assert.Empty(t, reqs[0].Header.Get("x-amz-checksum-algorithm"))
	assert.Empty(t, reqs[1].Header.Get("x-amz-checksum-algorithm"))
	assert.Empty(t, reqs[2].Header.Get("x-amz-checksum-algorithm"))
	body := string(reqs[completeRequestIndex].Body)
	assert.NotContains(t, body, "Checksum")
	assert.Contains(t, body, "<PartNumber>1</PartNumber>")
	assert.Contains(t, body, "<PartNumber>2</PartNumber>")
}

func TestPutObject_AWSDefaultChainMultipartChecksums(t *testing.T) {
	rec := &recorder{respond: awsChecksumResponder()}
	c := connectAWSDefaultRecorded(t, rec)
	out, err := c.PutObject(t.Context(), contracts.PutObjectInput{
		Key: testObjectKey, Body: &zeroReader{left: partSize + 1},
	})
	require.NoError(t, err)
	assert.Equal(t, `"final-etag"`, out.ETag)
	reqs := rec.all()
	require.Len(t, reqs, multipartRequestCount)
	assert.Equal(t, "CRC32", reqs[0].Header.Get("x-amz-checksum-algorithm"))
	assertCompleteChecksums(t, reqs[completeRequestIndex])
}

func TestPutObject_AWSStaticCredentials_MultipartChecksums(t *testing.T) {
	rec := &recorder{respond: awsChecksumResponder()}
	cfg := testConfigWithKeys()
	c := connectWithConfigRecorded(t, cfg, rec)
	out, err := c.PutObject(t.Context(), contracts.PutObjectInput{
		Key: testObjectKey, Body: &zeroReader{left: partSize + 1},
	})
	require.NoError(t, err)
	assert.Equal(t, `"final-etag"`, out.ETag)
	reqs := rec.all()
	require.Len(t, reqs, multipartRequestCount)
	assert.Equal(t, "CRC32", reqs[0].Header.Get("x-amz-checksum-algorithm"))
	assertCompleteChecksums(t, reqs[completeRequestIndex])
}

func TestPutObject_GCS_MultipartNoChecksums(t *testing.T) {
	rec := &recorder{respond: nonAWSChecksumResponder()}
	cfg := testConfigWithKeys()
	cfg.Endpoint = gcsServiceEndpoint
	c := connectWithConfigRecorded(t, cfg, rec)
	out, err := c.PutObject(t.Context(), contracts.PutObjectInput{
		Key: testObjectKey, Body: &zeroReader{left: partSize + 1},
	})
	require.NoError(t, err)
	assert.Equal(t, `"final-etag"`, out.ETag)
	assertNoMultipartChecksums(t, rec.all())
}

func TestPutObject_S3Compatible_MultipartNoChecksums(t *testing.T) {
	rec := &recorder{respond: nonAWSChecksumResponder()}
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://minio.internal.net"
	c := connectWithConfigRecorded(t, cfg, rec)
	out, err := c.PutObject(t.Context(), contracts.PutObjectInput{
		Key: testObjectKey, Body: &zeroReader{left: partSize + 1},
	})
	require.NoError(t, err)
	assert.Equal(t, `"final-etag"`, out.ETag)
	assertNoMultipartChecksums(t, rec.all())
}
