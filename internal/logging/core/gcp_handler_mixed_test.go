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
	testMixedID = 123
	testElem1   = 1
	testElem2   = 2
	testElem3   = 3
)

func parseMixedMessageLog(t *testing.T, logMsg, expectedHeader string) any {
	t.Helper()
	var buf bytes.Buffer
	logger := slog.New(core.NewGCPHandler(&buf, nil))
	logger.Info(logMsg)
	entry := parseLogMap(t, buf.Bytes())
	msgObj, ok := entry[testKeyMessage].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, expectedHeader, msgObj["header"])
	return msgObj["content"]
}

func TestGCPHandlerMixedMessage_Object(t *testing.T) {
	content := parseMixedMessageLog(t, `payload: {"action": "create", "id": 123}`, "payload:")
	contentObj, ok := content.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "create", contentObj["action"])
	assert.Equal(t, float64(testMixedID), contentObj["id"])
}

func TestGCPHandlerMixedMessage_Array(t *testing.T) {
	content := parseMixedMessageLog(t, "status [1, 2, 3]", "status")
	contentArr, ok := content.([]any)
	require.True(t, ok)
	expected := []any{float64(testElem1), float64(testElem2), float64(testElem3)}
	assert.Equal(t, expected, contentArr)
}
