package client_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func isInitiate(rec recordedRequest) bool {
	return rec.Method == http.MethodPost && rec.Query.Has("uploads")
}

func isComplete(rec recordedRequest) bool {
	return rec.Method == http.MethodPost && !rec.Query.Has("uploads")
}

func failingControl(shouldFail func(recordedRequest) bool) func(recordedRequest) *http.Response {
	ok := multipartResponder("0")
	return func(rec recordedRequest) *http.Response {
		if shouldFail(rec) {
			return xmlError(http.StatusForbidden, "AccessDenied")(rec)
		}
		return ok(rec)
	}
}

func TestPutObject_MultipartInitiateFailureNeedsNoAbort(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: failingControl(isInitiate)}
	_, err := putZeros(t, rec, partSize+1)
	require.ErrorIs(t, err, contracts.ErrAccessDenied)
	assert.Len(t, rec.all(), 1)
}

func TestPutObject_MultipartCompleteFailureAborts(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: failingControl(isComplete)}
	_, err := putZeros(t, rec, partSize+1)
	require.ErrorIs(t, err, contracts.ErrAccessDenied)
	reqs := rec.all()
	assert.Equal(t, http.MethodDelete, reqs[len(reqs)-1].Method)
}
