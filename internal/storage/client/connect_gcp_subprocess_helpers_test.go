package client

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envStorageGCPLifecycleHelper = "TEST_STORAGE_GCP_LIFECYCLE_HELPER"
	envStorageGCPFailureHelper   = "TEST_STORAGE_GCP_FAILURE_HELPER"
	errSubprocessStorageFail     = "storage subprocess network error"
)

var (
	rePlainTextDatePrefix = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}`)
	reSlogTextPattern     = regexp.MustCompile(`(^|\s)level=(DEBUG|INFO|WARN|ERROR)`)
	frameworkExact        = map[string]struct{}{"PASS": {}, "FAIL": {}}
	frameworkPrefixes     = []string{
		"=== RUN", "=== PAUSE", "=== CONT", "--- PASS:", "--- SKIP:", "--- FAIL:",
		"ok ", "FAIL ", "? ", "coverage:", "testing: warning: no tests to run",
	}
)

func TestStorageGCPLifecycleHelperProcess(t *testing.T) {
	if os.Getenv(envStorageGCPLifecycleHelper) != "1" {
		t.Skip("not running as storage gcp lifecycle helper process")
	}
	_, err := ConnectWithConfig(context.Background(), gcpStorageTestConfig(), WithSkipProbe(true))
	require.NoError(t, err)
}

func TestStorageGCPFailureHelperProcess(t *testing.T) {
	if os.Getenv(envStorageGCPFailureHelper) != "1" {
		t.Skip("not running as storage gcp failure helper process")
	}
	SetExitFunc(func(int) {})
	require.Error(t, handleConnectionError(context.Background(), errors.New(errSubprocessStorageFail)))
}

func hasFrameworkPrefix(trimmed string) bool {
	for _, p := range frameworkPrefixes {
		if strings.HasPrefix(trimmed, p) {
			return true
		}
	}
	return false
}

func isFrameworkLine(trimmed string) bool {
	if _, ok := frameworkExact[trimmed]; ok {
		return true
	}
	return hasFrameworkPrefix(trimmed)
}

func isGoTestFrameworkLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return true
	}
	return isFrameworkLine(trimmed)
}

func validateGCPLineFormat(t *testing.T, line string) {
	t.Helper()
	require.False(t, rePlainTextDatePrefix.MatchString(line), "plain-text date prefix detected in output: %s", line)
	require.False(t, reSlogTextPattern.MatchString(line), "slog text handler format detected in output: %s", line)
}

func validateParsedStorageGCPSeverity(t *testing.T, entry map[string]any) {
	t.Helper()
	sev, ok := entry[logFieldSeverity].(string)
	require.True(t, ok, "expected severity in log entry: %v", entry)
	assert.NotEqual(t, "DEFAULT", sev, "GCP severity must not be DEFAULT: %v", entry)
	_, valid := validGCPSeverities[sev]
	assert.True(t, valid, "severity %q is not in valid GCP severity enum: %v", sev, entry)
}

func parseAndValidateStorageGCPJSONEntry(t *testing.T, line string) map[string]any {
	t.Helper()
	var entry map[string]any
	err := json.Unmarshal([]byte(line), &entry)
	require.NoError(t, err, "non-framework line failed to parse as valid JSON: %q", line)
	assertStorageGCPTimestamp(t, entry)
	validateParsedStorageGCPSeverity(t, entry)
	return entry
}
