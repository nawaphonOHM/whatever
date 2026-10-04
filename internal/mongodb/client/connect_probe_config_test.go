package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectWithConfig_ProbeSequence_Step2_Failure(t *testing.T) {
	exitCalled, restore := interceptExitHook()
	defer restore()

	cleanup := setupMockPingAndProbe(errors.New(errStep2Substr + ": timeout"))
	defer cleanup()

	client, err := ConnectWithConfig(context.Background(), createUnreachableConfig())
	assert.Nil(t, client)
	assertProbeFailure(t, err, exitCalled, errStep2Substr)
}

func TestConnectWithConfig_ProbeSequence_Step3_Failure(t *testing.T) {
	exitCalled, restore := interceptExitHook()
	defer restore()

	cleanup := setupMockPingAndProbe(errors.New(errStep3Substr + ": maxTimeMS exceeded"))
	defer cleanup()

	client, err := ConnectWithConfig(context.Background(), createUnreachableConfig())
	assert.Nil(t, client)
	assertProbeFailure(t, err, exitCalled, errStep3Substr)
}

func TestConnectWithConfig_ProbeSequence_Step4_Failure(t *testing.T) {
	exitCalled, restore := interceptExitHook()
	defer restore()

	cleanup := setupMockPingAndProbe(errors.New(errStep4Substr + " \"users\": err"))
	defer cleanup()

	client, err := ConnectWithConfig(context.Background(), createUnreachableConfig())
	assert.Nil(t, client)
	assertProbeFailure(t, err, exitCalled, errStep4Substr)
}

func TestConnectWithConfig_ProbeSequence_Success(t *testing.T) {
	cleanup := setupMockPingAndProbe(nil)
	defer cleanup()

	client, err := ConnectWithConfig(context.Background(), createUnreachableConfig())
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NoError(t, client.Disconnect(context.Background()))
}
