package core_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	testAttrSecret   = "secret"
	testRedactedMask = "[REDACTED]"
)

func customFilterAttr(groups []string, a slog.Attr) slog.Attr {
	_ = groups
	if a.Key == testAttrSecret {
		return slog.String(testAttrSecret, testRedactedMask)
	}
	return a
}

func verifyBuilderEntry(t *testing.T, data []byte) {
	t.Helper()
	entry := parseLogMap(t, data)
	assert.Equal(t, core.GCPSeverityWarning, entry[testKeySeverity])
	assert.Equal(t, "disk full warning", entry[testKeyMessage])
	assert.Equal(t, testRedactedMask, entry[testAttrSecret])
	assert.Equal(t, "98%", entry["capacity"])
	assert.NotNil(t, entry[testKeyTimestamp])

	srcObj, ok := entry["logging.googleapis.com/sourceLocation"].(map[string]any)
	require.True(t, ok)
	assert.NotEmpty(t, srcObj["file"])
}

func TestBuildHandler_GCPModeOutput(t *testing.T) {
	t.Setenv("OHM9996_CLOUD_PROVIDER", "gcp")
	t.Setenv("OHM9996_OTEL_ENABLED", "false")

	cfg := &config.Config{
		CloudProvider: "gcp",
		Level:         "warn",
		AddSource:     true,
	}

	var buf bytes.Buffer
	logger := slog.New(core.BuildHandler(cfg, &buf, customFilterAttr))

	logger.Debug("should be ignored")
	assert.Empty(t, buf.Bytes())

	logger.Warn("disk full warning", testAttrSecret, "top_secret_token", "capacity", "98%")
	verifyBuilderEntry(t, buf.Bytes())
}
