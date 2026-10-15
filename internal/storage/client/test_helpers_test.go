package client_test

import (
	"bytes"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockTransport struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (m mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

func newMockHTTPClient(fn func(*http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{
		Transport: mockTransport{roundTrip: fn},
	}
}

func newHTTPResponse(statusCode int, body string, headers map[string]string) *http.Response {
	resp := &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
	for k, v := range headers {
		resp.Header.Set(k, v)
	}
	return resp
}

func buildRawClientForTest(httpClient *http.Client) *s3.Client {
	return s3.New(s3.Options{
		Region:           "us-east-1",
		RetryMaxAttempts: 1,
		HTTPClient:       httpClient,
	})
}

const (
	testAccessKey       = "AKIAIOSFODNN7EXAMPLE"
	testSecretKey       = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	testBucket          = "mybucket"
	testObjectKey       = "k"
	testVirtualPath     = "/k"
	testPathStylePath   = "/mybucket/k"
	testAWSRegionSE     = "ap-southeast-1"
	testAWSRegionalHost = "mybucket.s3.ap-southeast-1.amazonaws.com"
	testAWSVirtualHost  = "mybucket.s3.us-east-1.amazonaws.com"
	testMinioPort       = 9000
	testMinioHost       = "localhost:9000"
	sigV4Algorithm      = "AWS4-HMAC-SHA256"
)
