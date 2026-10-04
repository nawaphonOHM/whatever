package logging_test

import (
	"bytes"
	"testing"

	"github.com/nawaphonOHM/whatever/v2/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v2/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func verifyTextOutput(t *testing.T, output string) {
	t.Helper()
	assert.Contains(t, output, "level=INFO")
	assert.Contains(t, output, "msg=\"server starting\"")
	assert.Contains(t, output, "port=8080")
	assert.Contains(t, output, "level=WARN")
	assert.Contains(t, output, "msg=\"high latency detected\"")
	assert.Contains(t, output, "ms=250")
}

func TestLogging_TextOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	cleanup := core.CaptureTextLogs(buf, logging.LevelInfo)
	defer cleanup()

	logging.Info("server starting", "port", testPort)
	logging.Warn("high latency detected", "ms", testLatencyMs)

	verifyTextOutput(t, buf.String())
}
