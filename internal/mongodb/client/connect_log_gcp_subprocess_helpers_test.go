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
	envMongoGCPLifecycleHelper = "TEST_MONGO_GCP_LIFECYCLE_HELPER"
	envMongoGCPFailureHelper   = "TEST_MONGO_GCP_FAILURE_HELPER"
	errSubprocessMongoFail     = "mongodb subprocess network error"
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

func TestMongoDBGCPLifecycleHelperProcess(t *testing.T) {
	if os.Getenv(envMongoGCPLifecycleHelper) != "1" {
		t.Skip("not running as mongodb gcp lifecycle helper process")
	}
	stubMongoPingAndProbe(t)
	_, err := ConnectWithConfig(context.Background(), gcpMongoTestConfig())
	require.NoError(t, err)
}

func TestMongoDBGCPFailureHelperProcess(t *testing.T) {
	if os.Getenv(envMongoGCPFailureHelper) != "1" {
		t.Skip("not running as mongodb gcp failure helper process")
	}
	SetExitFunc(func(int) {})
	require.Error(t, handleConnectionError(context.Background(), errors.New(errSubprocessMongoFail)))
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

func validateParsedGCPSeverity(t *testing.T, entry map[string]any) {
	t.Helper()
	sev, ok := entry[logFieldSeverity].(string)
	require.True(t, ok, "expected severity in log entry: %v", entry)
	assert.NotEqual(t, "DEFAULT", sev, "GCP severity must not be DEFAULT: %v", entry)
	_, valid := validGCPSeverities[sev]
	assert.True(t, valid, "severity %q is not in valid GCP severity enum: %v", sev, entry)
}

func parseAndValidateGCPJSONEntry(t *testing.T, line string) map[string]any {
	t.Helper()
	var entry map[string]any
	err := json.Unmarshal([]byte(line), &entry)
	require.NoError(t, err, "non-framework line failed to parse as valid JSON: %q", line)
	assertGCPTimestamp(t, entry)
	validateParsedGCPSeverity(t, entry)
	return entry
}
