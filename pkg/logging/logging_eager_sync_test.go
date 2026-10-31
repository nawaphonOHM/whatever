package logging_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envPkgEagerHelper = "TEST_PKG_LOGGING_GCP_EAGER"
	msgPkgLogging     = "pkg logging eager startup"
	msgPkgSlog        = "pkg slog eager startup"
	keyModule         = "module"
)

func runPkgGCPEagerSubprocess(t *testing.T) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestLoggingGCPEagerHelperProcess$", "-test.v")
	cmd.Env = append(os.Environ(),
		envPkgEagerHelper+"=1",
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

func TestLoggingGCPEagerHelperProcess(t *testing.T) {
	if os.Getenv(envPkgEagerHelper) != "1" {
		t.Skip("not running as pkg eager helper process")
	}
	logging.Info(msgPkgLogging, keyModule, "pkg_logging")
	slog.Info(msgPkgSlog, keyModule, "pkg_slog")
	logging.Flush()
}

func isJSONString(s string) bool {
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}

func extractJSONLogLines(output string) []string {
	var jsonLines []string
	for line := range strings.SplitSeq(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if isJSONString(trimmed) {
			jsonLines = append(jsonLines, trimmed)
		}
	}
	return jsonLines
}

func parseEntries(t *testing.T, jsonLines []string) []map[string]any {
	t.Helper()
	entries := make([]map[string]any, len(jsonLines))
	for i, line := range jsonLines {
		require.NoError(t, json.Unmarshal([]byte(line), &entries[i]))
		assert.Equal(t, "INFO", entries[i]["severity"])
	}
	return entries
}

func checkEntryMatch(entry map[string]any, expMsg, expMod string) bool {
	return entry["message"] == expMsg && entry[keyModule] == expMod
}

func findMatchingEntry(entries []map[string]any, expMsg, expMod string) bool {
	for _, entry := range entries {
		if checkEntryMatch(entry, expMsg, expMod) {
			return true
		}
	}
	return false
}

func verifyPkgEntries(t *testing.T, jsonLines []string) {
	t.Helper()
	require.Len(t, jsonLines, 2)
	entries := parseEntries(t, jsonLines)

	assert.True(t, findMatchingEntry(entries, msgPkgLogging, "pkg_logging"))
	assert.True(t, findMatchingEntry(entries, msgPkgSlog, "pkg_slog"))
}

func TestLogging_GCPEagerInit(t *testing.T) {
	output := runPkgGCPEagerSubprocess(t)
	lines := extractJSONLogLines(output)
	verifyPkgEntries(t, lines)
}
