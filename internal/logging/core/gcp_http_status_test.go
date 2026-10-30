package core_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	testStatusUint200   = uint(200)
	testStatusFloat200  = float64(200)
	testStatusString404 = "404"
	testStatusFloat404  = float64(404)
	testHTTPKeyMethod   = "method"
	testHTTPKeyPath     = "path"
	testHTTPKeyStatus   = "status"
)

func verifyStatusLog(t *testing.T, rawVal any, expectedStatus float64) {
	t.Helper()
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request", testHTTPKeyMethod, "GET", testHTTPKeyPath, "/resource", testHTTPKeyStatus, rawVal)

	entry := parseLogMap(t, buf.Bytes())
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, expectedStatus, httpReq[testHTTPKeyStatus])
}

func TestGCPHandler_StatusTypes(t *testing.T) {
	verifyStatusLog(t, testStatusUint200, testStatusFloat200)
	verifyStatusLog(t, testStatusString404, testStatusFloat404)
}

func TestGCPHandler_FalsePositiveJobNotHTTP(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))

	logger.Info("job", testHTTPKeyMethod, "rsync", testHTTPKeyPath, "/tmp/file", testHTTPKeyStatus, "done")

	entry := parseLogMap(t, buf.Bytes())
	assert.Nil(t, entry[testKeyHTTPRequest])
	assert.Equal(t, "rsync", entry[testHTTPKeyMethod])
	assert.Equal(t, "/tmp/file", entry[testHTTPKeyPath])
	assert.Equal(t, "done", entry[testHTTPKeyStatus])
}

func TestGCPHandler_EmptyHTTPRequestMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))

	logger.Info("HTTP Request")

	entry := parseLogMap(t, buf.Bytes())
	assert.Nil(t, entry[testKeyHTTPRequest])
}
