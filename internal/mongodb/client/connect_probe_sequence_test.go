package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	errProbePrefix = "failed to probe mongodb"
	errStep2Substr = "failed to list collections without maxTimeMS"
	errStep3Substr = "failed to list collections with maxTimeMS"
	errStep4Substr = "failed to probe document in collection"
)

func assertProbeFailure(t *testing.T, err error, exitCalled *bool, substr string) {
	t.Helper()
	require.Error(t, err)
	assert.Contains(t, err.Error(), errProbePrefix)
	assert.Contains(t, err.Error(), substr)
	assert.True(t, *exitCalled)
}

func TestConnect_ProbeSequence_Step2_Failure(t *testing.T) {
	setupConnectPingEnv(t)
	exitCalled, cleanupExit := setupExitCapture()
	defer cleanupExit()

	cleanup := setupMockPingAndProbe(errors.New(errStep2Substr + ": deadline"))
	defer cleanup()

	client, err := Connect(context.Background())
	assert.Nil(t, client)
	assertProbeFailure(t, err, exitCalled, errStep2Substr)
}

func TestConnect_ProbeSequence_Step3_Failure(t *testing.T) {
	setupConnectPingEnv(t)
	exitCalled, cleanupExit := setupExitCapture()
	defer cleanupExit()

	cleanup := setupMockPingAndProbe(errors.New(errStep3Substr + ": server error"))
	defer cleanup()

	client, err := Connect(context.Background())
	assert.Nil(t, client)
	assertProbeFailure(t, err, exitCalled, errStep3Substr)
}

func TestConnect_ProbeSequence_Step4_Failure(t *testing.T) {
	setupConnectPingEnv(t)
	exitCalled, cleanupExit := setupExitCapture()
	defer cleanupExit()

	cleanup := setupMockPingAndProbe(errors.New(errStep4Substr + " \"users\": err"))
	defer cleanup()

	client, err := Connect(context.Background())
	assert.Nil(t, client)
	assertProbeFailure(t, err, exitCalled, errStep4Substr)
}

func TestConnect_ProbeSequence_Success(t *testing.T) {
	setupConnectPingEnv(t)
	cleanup := setupMockPingAndProbe(nil)
	defer cleanup()

	client, err := Connect(context.Background())
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.NoError(t, client.Disconnect(context.Background()))
}
