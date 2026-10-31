package client

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runStorageGCPProcess(t *testing.T, helperName, envFlag string) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^"+helperName+"$", "-test.v")
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"TMPDIR=" + os.Getenv("TMPDIR"),
		"HOME=" + os.Getenv("HOME"),
		envFlag + "=1",
		"OHM9996_CLOUD_PROVIDER=gcp",
		"OHM9996_OTEL_ENABLED=false",
		"OTEL_SDK_DISABLED=true",
		"OHM9996_LOGGING_OUTPUT=stdout",
		"OHM9996_LOGGING_LEVEL=info",
		"OHM9996_LOGGING_FORMAT=json",
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	require.NoError(t, cmd.Run())
	return out.String()
}

func extractStorageGCPJSONLines(t *testing.T, output string) []map[string]any {
	t.Helper()
	var entries []map[string]any
	for line := range strings.SplitSeq(output, "\n") {
		if isGoTestFrameworkLine(line) {
			continue
		}
		trimmed := strings.TrimSpace(line)
		validateGCPLineFormat(t, trimmed)
		entries = append(entries, parseAndValidateStorageGCPJSONEntry(t, trimmed))
	}
	return entries
}

func TestStorage_GCPSubprocess_Lifecycle(t *testing.T) {
	output := runStorageGCPProcess(t, "TestStorageGCPLifecycleHelperProcess", envStorageGCPLifecycleHelper)
	entries := extractStorageGCPJSONLines(t, output)
	require.NotEmpty(t, entries)
	assertStorageInitEntryGCP(t, findStorageLog(entries, msgStorageInit))
	assertStorageCompleteEntryGCP(t, findStorageLog(entries, msgStorageComplete))
}

func TestStorage_GCPSubprocess_Failure(t *testing.T) {
	output := runStorageGCPProcess(t, "TestStorageGCPFailureHelperProcess", envStorageGCPFailureHelper)
	entries := extractStorageGCPJSONLines(t, output)
	require.NotEmpty(t, entries)
	failEntry := findStorageLog(entries, msgStorageFail)
	require.NotNil(t, failEntry)
	assertStorageGCPTimestamp(t, failEntry)
	assertStorageGCPSeverity(t, failEntry, "ERROR")
	assert.Equal(t, errSubprocessStorageFail, failEntry["error"])
	assert.Equal(t, "ExitAbnormal", failEntry["exit_flag"])
}
