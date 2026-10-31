package core_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envEagerInitHelper = "TEST_GCP_EAGER_INIT_HELPER"
	msgEagerSlog       = "eager startup slog message"
	msgEagerCentral    = "eager startup central message"
)

type eagerExpectation struct {
	msg string
	key string
	val string
}

func runGCPEagerSubprocess(t *testing.T) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestGCPEagerInitHelperProcess$", "-test.v")
	cmd.Env = append(os.Environ(),
		envEagerInitHelper+"=1",
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

func TestGCPEagerInitHelperProcess(t *testing.T) {
	if os.Getenv(envEagerInitHelper) != "1" {
		t.Skip("not running as eager helper process")
	}
	slog.Info(msgEagerSlog, "subsystem", "client")
	central.DefaultWorker().Logger().InfoContext(
		context.Background(),
		msgEagerCentral,
		slog.String("worker", "default"),
	)
	core.Flush()
}

func verifyEagerEntry(t *testing.T, line string, exp eagerExpectation) {
	t.Helper()
	var entry map[string]any
	require.NoError(t, json.Unmarshal([]byte(line), &entry))
	assert.Equal(t, "INFO", entry["severity"])
	assert.Equal(t, exp.msg, entry["message"])
	assert.Equal(t, exp.val, entry[exp.key])
}

func verifyEagerLines(t *testing.T, jsonLines []string) {
	t.Helper()
	require.Len(t, jsonLines, 2)
	verifyEagerEntry(t, jsonLines[0], eagerExpectation{msg: msgEagerSlog, key: "subsystem", val: "client"})
	verifyEagerEntry(t, jsonLines[1], eagerExpectation{msg: msgEagerCentral, key: "worker", val: "default"})
}

func TestGCPEagerInit_DirectSlogOutputsGCPJSON(t *testing.T) {
	output := runGCPEagerSubprocess(t)
	assertNoPlainTextLogLeaks(t, output)
	verifyEagerLines(t, filterLogJSONLines(output))
}
