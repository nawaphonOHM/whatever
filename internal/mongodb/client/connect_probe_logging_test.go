package client

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProbeSuccessLogsTargetCollection_JSONAndGCP(t *testing.T) {
	for _, mode := range clientLogModes() {
		t.Run(mode.name, func(t *testing.T) { assertProbeSuccessLog(t, mode) })
	}
}

func assertProbeSuccessLog(t *testing.T, mode clientLogMode) {
	t.Helper()
	var output bytes.Buffer
	setupWorkerLogger(t, mode.handler(&output))
	logProbeSuccess(t.Context(), "orders")
	entry := findClientEntry(parseClientLogs(t, output.Bytes()), msgProbeSuccess)
	require.NotNil(t, entry)
	assert.Equal(t, logLevelInfo, entry[mode.levelKey])
	assert.Equal(t, "orders", entry["target_collection"])
}
