package client_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	listPageFormat = `<ListBucketResult><IsTruncated>%t</IsTruncated>%s%s</ListBucketResult>`
	listTokenXML   = `<NextContinuationToken>tok</NextContinuationToken>`

	maxRepeatedTokenRequests = 3
)

func listXML(truncated bool, token, entries string) string {
	return fmt.Sprintf(listPageFormat, truncated, token, entries)
}

func wrapEach(format string, items []string) string {
	parts := make([]string, len(items))
	for i, item := range items {
		parts[i] = fmt.Sprintf(format, item)
	}
	return strings.Join(parts, "")
}

func objectsXML(keys ...string) string {
	return wrapEach("<Contents><Key>%s</Key><Size>1</Size></Contents>", keys)
}

func prefixesXML(prefixes ...string) string {
	return wrapEach("<CommonPrefixes><Prefix>%s</Prefix></CommonPrefixes>", prefixes)
}

func newListClient(rec *recorder) *client.Client {
	return client.NewClient(client.ClientParams{
		Raw:           buildRawClientForTest(rec.httpClient()),
		DefaultBucket: "my-bucket",
	})
}

func TestListObjects_ErrorOnLaterPageDiscardsPartialResult(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: func(req recordedRequest) *http.Response {
		if req.Query.Get("continuation-token") == "" {
			return newHTTPResponse(http.StatusOK, listXML(true, listTokenXML, objectsXML("a", "b")), nil)
		}
		return newHTTPResponse(http.StatusInternalServerError, "", nil)
	}}
	out, err := newListClient(rec).ListObjects(t.Context(), contracts.ListObjectsInput{})
	require.Error(t, err)
	assert.Nil(t, out)
	assert.Len(t, rec.all(), 2)
}

func TestListObjects_RepeatedTokenDoesNotLoopForever(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: func(recordedRequest) *http.Response {
		return newHTTPResponse(http.StatusOK, listXML(true, listTokenXML, objectsXML("a")), nil)
	}}
	out, err := newListClient(rec).ListObjects(t.Context(), contracts.ListObjectsInput{})
	require.NoError(t, err)
	assert.LessOrEqual(t, len(rec.all()), maxRepeatedTokenRequests)
	assert.NotEmpty(t, out.Objects)
}

func TestListObjects_MaxKeysCountsPrefixesAndStopsPaging(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: func(recordedRequest) *http.Response {
		entries := objectsXML("a") + prefixesXML("p1/", "p2/")
		return newHTTPResponse(http.StatusOK, listXML(true, listTokenXML, entries), nil)
	}}
	input := contracts.ListObjectsInput{MaxKeys: 2, Delimiter: "/"}
	out, err := newListClient(rec).ListObjects(t.Context(), input)
	require.NoError(t, err)
	assert.Len(t, rec.all(), 1)
	assert.Equal(t, "2", rec.all()[0].Query.Get("max-keys"))
	assert.Len(t, out.Objects, 1)
	assert.Equal(t, []string{"p1/"}, out.Prefixes)
}
