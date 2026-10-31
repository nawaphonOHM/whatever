package storage_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envStorageEagerHelper = "TEST_PKG_STORAGE_GCP_EAGER"
	msgStorageSlog        = "bare storage import startup slog"
)

func runStorageGCPEagerSubprocess(t *testing.T) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestStorageEagerInitHelperProcess$", "-test.v")
	cmd.Env = append(os.Environ(),
		envStorageEagerHelper+"=1",
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

func TestStorageEagerInitHelperProcess(t *testing.T) {
	if os.Getenv(envStorageEagerHelper) != "1" {
		t.Skip("not running as storage eager helper process")
	}
	_ = storage.DefaultConfig()
	slog.Info(msgStorageSlog, "component", "storage_client")
}

func isJSONLine(line string) bool {
	return strings.HasPrefix(line, "{") && strings.HasSuffix(line, "}")
}

func parseJSONCandidate(line string) (map[string]any, bool) {
	if !isJSONLine(line) {
		return nil, false
	}
	var entry map[string]any
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return nil, false
	}
	return entry, true
}

func extractJSONLine(output string) (map[string]any, bool) {
	for line := range strings.SplitSeq(output, "\n") {
		if entry, ok := parseJSONCandidate(strings.TrimSpace(line)); ok {
			return entry, true
		}
	}
	return nil, false
}

func TestStorage_BareImportEagerGCPLogging(t *testing.T) {
	output := runStorageGCPEagerSubprocess(t)
	entry, ok := extractJSONLine(output)
	require.True(t, ok, "expected valid GCP JSON log entry in output:\n%s", output)
	assert.Equal(t, "INFO", entry["severity"])
	assert.Equal(t, msgStorageSlog, entry["message"])
	assert.Equal(t, "storage_client", entry["component"])
	assert.NotEmpty(t, entry["timestamp"])
}

func TestStorage_ProductionDependenciesIncludeCore(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", "github.com/nawaphonOHM/whatever/v3/pkg/storage")
	out, err := cmd.Output()
	require.NoError(t, err)
	assert.Contains(t, string(out), "github.com/nawaphonOHM/whatever/v3/internal/logging/core")
}
