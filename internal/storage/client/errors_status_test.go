package client_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

var (
	errObject = contracts.ErrObjectNotFound
	errBucket = contracts.ErrBucketNotFound
	errDenied = contracts.ErrAccessDenied
)

func newStatusClient(responder func(recordedRequest) *http.Response) *client.Client {
	rec := &recorder{respond: responder}
	return client.NewClient(client.ClientParams{
		Raw:           buildRawClientForTest(rec.httpClient()),
		DefaultBucket: "my-bucket",
	})
}

func xmlError(status int, code string) func(recordedRequest) *http.Response {
	return func(recordedRequest) *http.Response {
		return newHTTPResponse(status, "<Error><Code>"+code+"</Code><Message>m</Message></Error>", nil)
	}
}

// headRoute answers object HEADs with objectStatus and bucket HEADs (path "/") with bucketStatus.
func headRoute(objectStatus, bucketStatus int) func(recordedRequest) *http.Response {
	return func(rec recordedRequest) *http.Response {
		if rec.Path == "/" {
			return newHTTPResponse(bucketStatus, "", nil)
		}
		return newHTTPResponse(objectStatus, "", nil)
	}
}

func TestErrorMapping_ObjectVersusBucketVersusAccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		respond func(recordedRequest) *http.Response
		call    opCall
		want    error
		name    string
	}{
		{xmlError(http.StatusNotFound, "NoSuchKey"), callGet, errObject, "get no such key"},
		{xmlError(http.StatusNotFound, "NoSuchBucket"), callGet, errBucket, "get no such bucket"},
		{xmlError(http.StatusForbidden, "AccessDenied"), callGet, errDenied, "get access denied"},
		{xmlError(http.StatusForbidden, "AccessDenied"), callDelete, errDenied, "delete denied"},
		{xmlError(http.StatusNotFound, "NoSuchBucket"), callDelete, errBucket, "delete bucket"},
		{xmlError(http.StatusNotFound, "NoSuchBucket"), callList, errBucket, "list bucket"},
		{xmlError(http.StatusForbidden, "AccessDenied"), callList, errDenied, "list denied"},
		{headRoute(http.StatusNotFound, http.StatusOK), callHead, errObject, "head object"},
		{headRoute(http.StatusNotFound, http.StatusNotFound), callHead, errBucket, "head bucket"},
		{headRoute(http.StatusNotFound, http.StatusForbidden), callHead, errObject, "head 404 then 403"},
		{headRoute(http.StatusForbidden, http.StatusOK), callHead, errDenied, "head denied"},
		{headRoute(http.StatusForbidden, http.StatusOK), callExists, errDenied, "exists denied"},
		{headRoute(http.StatusNotFound, http.StatusNotFound), callExists, errBucket, "exists bucket"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.call(t.Context(), newStatusClient(tc.respond))
			require.ErrorIs(t, err, tc.want)
			assertNotOtherSentinels(t, err, tc.want)
		})
	}
}

func assertNotOtherSentinels(t *testing.T, err, want error) {
	t.Helper()
	for _, s := range []error{contracts.ErrObjectNotFound, contracts.ErrBucketNotFound, contracts.ErrAccessDenied} {
		if s != want {
			assert.NotErrorIs(t, err, s)
		}
	}
}

func TestExists_FalseOnlyWhenObjectIsMissing(t *testing.T) {
	t.Parallel()

	c := newStatusClient(headRoute(http.StatusNotFound, http.StatusOK))
	exists, err := c.Exists(t.Context(), "", testObjectKey)
	require.NoError(t, err)
	assert.False(t, exists)
}
