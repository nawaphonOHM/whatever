package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	msgInitBlueprint = "initializing REST server from blueprint"
	msgInitRegs      = "initializing REST server from registrations"
)

func assertBlueprintJSONLogs(t *testing.T, entries []map[string]any, count int) {
	t.Helper()
	regEntry := findLogEntry(entries, msgInitRegs)
	require.NotNil(t, regEntry)
	assert.Equal(t, "INFO", regEntry[levelField])
	assert.Equal(t, float64(count), regEntry["registration_count"])

	bpEntry := findLogEntry(entries, msgInitBlueprint)
	require.NotNil(t, bpEntry)
	assert.Equal(t, "INFO", bpEntry[levelField])
}

func TestServer_BlueprintLogging_JSON(t *testing.T) {
	buf, reset := setupLifecycleLogger(t)
	defer reset()

	regs := sampleRegs()
	srv, err := NewFromRegistrations(regs)
	require.NoError(t, err)
	require.NotNil(t, srv)
	assertBlueprintJSONLogs(t, parseJSONLogs(t, buf.Bytes()), len(regs))
}

func assertBlueprintGCPLogs(t *testing.T, entries []map[string]any, count int) {
	t.Helper()
	regEntry := findGCPLogEntry(entries, msgInitRegs)
	require.NotNil(t, regEntry)
	assert.Equal(t, gcpInfoLevel, regEntry[sevField])
	assert.Equal(t, float64(count), regEntry["registration_count"])

	bpEntry := findGCPLogEntry(entries, msgInitBlueprint)
	require.NotNil(t, bpEntry)
	assert.Equal(t, gcpInfoLevel, bpEntry[sevField])
}

func TestServer_BlueprintLogging_GCP(t *testing.T) {
	buf, reset := setupGCPLifecycleLogger(t)
	defer reset()

	regs := sampleRegs()
	srv, err := NewFromRegistrations(regs)
	require.NoError(t, err)
	require.NotNil(t, srv)
	assertBlueprintGCPLogs(t, parseJSONLogs(t, buf.Bytes()), len(regs))
}
