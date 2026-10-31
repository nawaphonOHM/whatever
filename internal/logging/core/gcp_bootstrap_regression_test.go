package core_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const (
	envBootstrapHelper = "TEST_GCP_BOOTSTRAP_HELPER"
	testBootstrapMsg   = "bootstrap test message"
)

func runGCPBootstrapSubprocess(t *testing.T) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestGCPBootstrapHelperProcess$", "-test.v")
	cmd.Env = append(os.Environ(),
		envBootstrapHelper+"=1",
		"OHM9996_CLOUD_PROVIDER=gcp",
		"OHM9996_LOGGING_OUTPUT=stdout",
		"OHM9996_LOGGING_LEVEL=info",
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	require.NoError(t, cmd.Run())
	return out.String()
}

func isJSONLine(s string) bool {
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}

func filterLogJSONLines(output string) []string {
	var jsonLines []string
	for line := range strings.SplitSeq(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if isJSONLine(trimmed) {
			jsonLines = append(jsonLines, trimmed)
		}
	}
	return jsonLines
}

func assertSingleGCPLine(t *testing.T, trimmed string) {
	t.Helper()
	var entry map[string]any
	require.NoError(t, json.Unmarshal([]byte(trimmed), &entry), "must be JSON: %s", trimmed)
	assert.NotEmpty(t, entry["severity"], "every GCP log must have severity: %s", trimmed)
	assert.NotContains(t, trimmed, "INFO loading REST configuration")
}

func checkLineForLeaks(t *testing.T, line string) {
	t.Helper()
	trimmed := strings.TrimSpace(line)
	if isTestHarnessLine(trimmed) {
		return
	}
	assertSingleGCPLine(t, trimmed)
}

func assertNoPlainTextLogLeaks(t *testing.T, output string) {
	t.Helper()
	for line := range strings.SplitSeq(output, "\n") {
		checkLineForLeaks(t, line)
	}
}

func TestGCPBootstrapHelperProcess(t *testing.T) {
	if os.Getenv(envBootstrapHelper) != "1" {
		t.Skip("not running as helper process")
	}
	core.ResetDefault()
	core.Info(testBootstrapMsg)
	core.Flush()
}

func TestGCPBootstrap_NoPlainTextPreInitLines(t *testing.T) {
	output := runGCPBootstrapSubprocess(t)
	assertNoPlainTextLogLeaks(t, output)
	jsonLines := filterLogJSONLines(output)
	require.NotEmpty(t, jsonLines, "expected at least one JSON log line")
	assert.Contains(t, jsonLines[len(jsonLines)-1], testBootstrapMsg)
}
