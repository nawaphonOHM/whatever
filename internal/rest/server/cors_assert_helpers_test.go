package server

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertCORSHeaders(t *testing.T, rec *httptest.ResponseRecorder, exp expectedHeaders) {
	t.Helper()
	assert.Equal(t, exp.origin, rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, originHeaderVal, rec.Header().Get("Vary"))
	assert.Equal(t, exp.methods, rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, exp.headers, rec.Header().Get("Access-Control-Allow-Headers"))
}

func assertNoCORSHeaders(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Headers"))
	assert.Empty(t, rec.Header().Get("Vary"))
}

func assertDisallowedCORS(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, originHeaderVal, rec.Header().Get("Vary"))
}
