package core_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	testHTTPStatus201Num = 201
	testHTTPStatus404Num = 404
	testMethodGETName    = "GET"
	testKeyPathName      = "path"
	testKeyStatusName    = "status"
)

func verifyMergeNonHTTP(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testMethodGETName, httpReq[testKeyReqMethod])
	assert.Equal(t, "/a", httpReq[testKeyReqURL])
	assert.Equal(t, "/tmp/file", entry[testKeyPathName])
	assert.Equal(t, float64(testHTTPStatus201Num), entry[testKeyStatusName])
}

func TestGCPHandler_MergeExplicitAndFlat_NonHTTPMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	explicitReq := &core.HttpRequest{RequestMethod: testMethodGETName, RequestURL: "/a"}
	logger.Info("job", "httpRequest", explicitReq,
		testKeyPathName, "/tmp/file", testKeyStatusName, testHTTPStatus201Num)

	verifyMergeNonHTTP(t, buf.Bytes())
}

func verifyMergeHTTP(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testMethodGETName, httpReq[testKeyReqMethod])
	assert.Equal(t, float64(testHTTPStatus404Num), httpReq[testKeyStatusName])
}

func TestGCPHandler_MergeExplicitAndFlat_HTTPMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	explicitReq := &core.HttpRequest{RequestMethod: testMethodGETName}
	logger.Info("HTTP Request", testKeyStatusName, testHTTPStatus404Num, "httpRequest", explicitReq)

	verifyMergeHTTP(t, buf.Bytes())
}
