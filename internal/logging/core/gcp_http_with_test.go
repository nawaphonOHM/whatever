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
	testMethodGETStr = "GET"
	testPathXStr     = "/x"
	testKeyMethodStr = "method"
	testKeyPathStr   = "path"
	testKeyStatusStr = "status"
)

func verifyUngroupedWithHTTP(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testMethodGETStr, httpReq[testKeyReqMethod])
	assert.Equal(t, testPathXStr, httpReq[testKeyReqURL])
	assert.Equal(t, float64(testHTTPStatus200Num), httpReq[testKeyStatusStr])
	assert.Nil(t, entry[testKeyMethodStr])
	assert.Nil(t, entry[testKeyPathStr])
}

func TestGCPHandler_UngroupedWithAttrs_IncludedInHTTP(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	derived := logger.With(testKeyMethodStr, testMethodGETStr, testKeyPathStr, testPathXStr)
	derived.Info("HTTP Request", testKeyStatusStr, testHTTPStatus200Num)

	verifyUngroupedWithHTTP(t, buf.Bytes())
}

func TestGCPHandler_UngroupedWithAttrs_NonHTTP(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	derived := logger.With(testKeyMethodStr, testMethodGETStr, testKeyPathStr, testPathXStr)
	derived.Info("job finished")

	entry := parseLogMap(t, buf.Bytes())
	assert.Nil(t, entry[testKeyHTTPRequest])
	assert.Equal(t, testMethodGETStr, entry[testKeyMethodStr])
	assert.Equal(t, testPathXStr, entry[testKeyPathStr])
}

func verifyGroupedWithHTTP(t *testing.T, data []byte) {
	entry := parseLogMap(t, data)
	assert.Nil(t, entry[testKeyHTTPRequest])

	gMap, ok := entry[testGroupName].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testMethodGETStr, gMap[testKeyMethodStr])
	assert.Equal(t, testPathXStr, gMap[testKeyPathStr])
	assert.Equal(t, float64(testHTTPStatus200Num), gMap[testKeyStatusStr])
}

func TestGCPHandler_GroupedWithAttrs_NotExtractedToHTTP(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	derived := logger.WithGroup(testGroupName).With(testKeyMethodStr, testMethodGETStr, testKeyPathStr, testPathXStr)
	derived.Info("HTTP Request", testKeyStatusStr, testHTTPStatus200Num)

	verifyGroupedWithHTTP(t, buf.Bytes())
}

func verifyDerivedBaseLine(t *testing.T, line []byte) {
	entry := parseLogMap(t, line)
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, testMethodGETStr, httpReq[testKeyReqMethod])
	assert.Equal(t, testPathXStr, httpReq[testKeyReqURL])
	assert.Nil(t, entry[testKeyMethodStr])
}

func verifyDerivedGroupLine(t *testing.T, line []byte) {
	entry := parseLogMap(t, line)
	assert.Nil(t, entry[testKeyHTTPRequest])
	assert.Equal(t, testMethodGETStr, entry[testKeyMethodStr])
	assert.Equal(t, testPathXStr, entry[testKeyPathStr])
	gMap, ok := entry[testGroupName].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(testHTTPStatus200Num), gMap[testKeyStatusStr])
}

const (
	testHTTPMsg      = "HTTP Request"
	testGroupName    = "g"
	testDerivedLines = 3
)

func emitDerivedLines(buf *bytes.Buffer) [][]byte {
	logger := slog.New(core.NewGCPHandler(buf, nil))
	base := logger.With(testKeyMethodStr, testMethodGETStr, testKeyPathStr, testPathXStr)
	grouped := base.WithGroup(testGroupName)
	base.Info(testHTTPMsg, testKeyStatusStr, testHTTPStatus200Num)
	grouped.Info(testHTTPMsg, testKeyStatusStr, testHTTPStatus200Num)
	base.Info(testHTTPMsg, testKeyStatusStr, testHTTPStatus200Num)
	return bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
}

func TestGCPHandler_WithThenWithGroup_NoDuplicationOrLoss(t *testing.T) {
	var buf bytes.Buffer
	lines := emitDerivedLines(&buf)
	require.Len(t, lines, testDerivedLines)
	verifyDerivedBaseLine(t, lines[0])
	verifyDerivedGroupLine(t, lines[1])
	verifyDerivedBaseLine(t, lines[2])
}
