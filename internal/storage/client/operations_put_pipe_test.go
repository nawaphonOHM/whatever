package client_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func pipeWithPayload(t *testing.T) *os.File {
	t.Helper()
	pr, pw, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, pr.Close()) })
	_, err = pw.WriteString(streamPayload)
	require.NoError(t, err)
	require.NoError(t, pw.Close())
	return pr
}

// An *os.File pipe has a Seek method that always fails, so it must be treated as a plain stream.
func TestPutObject_PipeFileIsStreamed(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: putOKResponder}
	c := connectRecorded(t, endpointConfig(testHTTPEndpoint), rec)
	_, err := c.PutObject(t.Context(), contracts.PutObjectInput{Key: testObjectKey, Body: pipeWithPayload(t)})
	require.NoError(t, err)
	reqs := rec.all()
	require.Len(t, reqs, 1)
	assert.Equal(t, streamPayload, string(reqs[0].Body))
	assert.Equal(t, streamPayloadSize, reqs[0].ContentLength)
}
