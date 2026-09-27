package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	testMockUser = "mockuser"
	testMockPass = "mockpass"
	testMockPort = 59999
)

func createValidTestConfig() *config.Config {
	return &config.Config{
		Host:               "127.0.0.1",
		Port:               testMockPort,
		Username:           testMockUser,
		Password:           testMockPass,
		Protocol:           config.ProtocolMongoDB,
		UUIDRepresentation: config.UUIDRepresentationUnspecified,
	}
}

func setupMockPing(fn func(context.Context, *mongo.Client) error) func() {
	origPing := pingClient
	pingClient = fn
	return func() { pingClient = origPing }
}

func setupExitCapture() (*bool, func()) {
	var exitCalled bool
	prev := SetExitFunc(func(int) { exitCalled = true })
	return &exitCalled, func() { SetExitFunc(prev) }
}

// TestConnect_TLSFallback_Success tests fallback to TLS when initial connect fails with TLS error.
func TestConnect_TLSFallback_Success(t *testing.T) {
	var attempts int
	cleanup := setupMockPing(func(context.Context, *mongo.Client) error {
		attempts++
		if attempts == 1 {
			return errors.New(testErrMsgTLS)
		}
		return nil
	})
	defer cleanup()

	client, err := ConnectWithConfig(context.Background(), createValidTestConfig())
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, 2, attempts)
}

// TestConnect_TLSFallback_RetryFails tests failure when both plain and TLS attempts fail.
func TestConnect_TLSFallback_RetryFails(t *testing.T) {
	var attempts int
	cleanupPing := setupMockPing(func(context.Context, *mongo.Client) error {
		attempts++
		return errors.New(testErrMsgTLS)
	})
	defer cleanupPing()

	exitCalled, cleanupExit := setupExitCapture()
	defer cleanupExit()

	client, err := ConnectWithConfig(context.Background(), createValidTestConfig())
	require.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, 2, attempts)
	assert.True(t, *exitCalled)
}

// TestConnect_NonTLSError_NoRetry tests immediate termination on non-TLS error without retry.
func TestConnect_NonTLSError_NoRetry(t *testing.T) {
	var attempts int
	cleanupPing := setupMockPing(func(context.Context, *mongo.Client) error {
		attempts++
		return errors.New(testErrMsgNonTLS)
	})
	defer cleanupPing()

	exitCalled, cleanupExit := setupExitCapture()
	defer cleanupExit()

	client, err := ConnectWithConfig(context.Background(), createValidTestConfig())
	require.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, 1, attempts)
	assert.True(t, *exitCalled)
}

// TestConnect_ImmediateSuccess tests direct successful unencrypted connection without fallback.
func TestConnect_ImmediateSuccess(t *testing.T) {
	var attempts int
	cleanup := setupMockPing(func(context.Context, *mongo.Client) error {
		attempts++
		return nil
	})
	defer cleanup()

	client, err := ConnectWithConfig(context.Background(), createValidTestConfig())
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, 1, attempts)
}
