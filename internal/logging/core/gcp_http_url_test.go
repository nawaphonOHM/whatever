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
	testPathKey = "path"
)

func TestGCPHandler_URLEscaping_PathAndQuery(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request",
		"method", "GET",
		testPathKey, "/files/a b/c+d",
		"query", "x=1&tag=a%26b",
	)

	entry := parseLogMap(t, buf.Bytes())
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "/files/a%20b/c+d?x=1&tag=a%26b", httpReq["requestUrl"])
	assert.Nil(t, entry[testPathKey])
	assert.Nil(t, entry["query"])
}

func TestGCPHandler_URLEscaping_AlreadyEncodedPath(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info("HTTP Request",
		"method", "GET",
		testPathKey, "/files/a%20b%2Fc",
		"query", "filter=yes",
	)

	entry := parseLogMap(t, buf.Bytes())
	httpReq, ok := entry[testKeyHTTPRequest].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "/files/a%20b%2Fc?filter=yes", httpReq["requestUrl"])
}
