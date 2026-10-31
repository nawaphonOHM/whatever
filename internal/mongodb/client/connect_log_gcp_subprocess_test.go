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

func runMongoGCPProcess(t *testing.T, helperName, envFlag string) string {
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

func extractGCPJSONLines(t *testing.T, output string) []map[string]any {
	t.Helper()
	var entries []map[string]any
	for line := range strings.SplitSeq(output, "\n") {
		if isGoTestFrameworkLine(line) {
			continue
		}
		trimmed := strings.TrimSpace(line)
		validateGCPLineFormat(t, trimmed)
		entries = append(entries, parseAndValidateGCPJSONEntry(t, trimmed))
	}
	return entries
}

func TestMongoDB_GCPSubprocess_Lifecycle(t *testing.T) {
	output := runMongoGCPProcess(t, "TestMongoDBGCPLifecycleHelperProcess", envMongoGCPLifecycleHelper)
	entries := extractGCPJSONLines(t, output)
	require.NotEmpty(t, entries)
	assertMongoInitEntryGCP(t, findClientEntry(entries, msgInitClient))
	assertMongoOptsEntryGCP(t, findClientEntry(entries, msgResolveOpts))
	assertMongoPingEntryGCP(t, findClientEntry(entries, msgPingSuccess))
	assertMongoCompleteEntryGCP(t, findClientEntry(entries, msgInitSuccess))
}

func TestMongoDB_GCPSubprocess_Failure(t *testing.T) {
	output := runMongoGCPProcess(t, "TestMongoDBGCPFailureHelperProcess", envMongoGCPFailureHelper)
	entries := extractGCPJSONLines(t, output)
	require.NotEmpty(t, entries)
	failEntry := findClientEntry(entries, msgMongoFail)
	require.NotNil(t, failEntry)
	assertGCPTimestamp(t, failEntry)
	assertGCPSeverity(t, failEntry, logLevelError)
	assert.Equal(t, errSubprocessMongoFail, failEntry["error"])
	assert.Equal(t, "ExitAbnormal", failEntry["exit_flag"])
}
