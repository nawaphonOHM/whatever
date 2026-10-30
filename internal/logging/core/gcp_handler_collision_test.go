package core_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

func TestGCPHandler_MultiCollisionDeterministic(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, nil)
	logger := slog.New(handler)

	logger.Info("main log message",
		"message", "attr_msg_val",
		"custom_message", "existing_custom_val",
		"severity", "attr_sev_val",
		"custom_severity", "existing_sev_val",
	)

	entry := parseLogMap(t, buf.Bytes())
	assert.Equal(t, "main log message", entry[testKeyMessage])
	assert.Equal(t, "existing_custom_val", entry["custom_message"])
	assert.Equal(t, "attr_msg_val", entry["custom_message_1"])
	assert.Equal(t, "existing_sev_val", entry["custom_severity"])
	assert.Equal(t, "attr_sev_val", entry["custom_severity_1"])
}

func TestGCPHandler_SourceLocation(t *testing.T) {
	var buf bytes.Buffer
	handler := core.NewGCPHandler(&buf, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(handler)

	logger.Info("source test")

	entry := parseLogMap(t, buf.Bytes())
	srcObj, ok := entry["logging.googleapis.com/sourceLocation"].(map[string]any)
	assert.True(t, ok)
	assert.NotEmpty(t, srcObj["file"])
	assert.NotEmpty(t, srcObj["line"])
}
