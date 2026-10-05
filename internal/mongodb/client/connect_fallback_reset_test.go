package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func runFallbackResetTest(t *testing.T, initialErr error) {
	var attempts int
	cleanup := setupMockPing(func(context.Context, *mongo.Client) error {
		attempts++
		if attempts == 1 {
			return initialErr
		}
		return nil
	})
	defer cleanup()

	client, err := ConnectWithConfig(context.Background(), createValidTestConfig())
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, 2, attempts)
}

func TestConnect_TLSFallback_ConnectionReset(t *testing.T) {
	runFallbackResetTest(t, errors.New("read tcp 127.0.0.1:59999: read: connection reset by peer"))
}

func TestConnect_TLSFallback_IncompleteRead(t *testing.T) {
	runFallbackResetTest(t, errors.New("incomplete read of full message: connection reset by peer"))
}

func TestConnect_TLSFallback_BrokenPipe(t *testing.T) {
	runFallbackResetTest(t, errors.New("write: broken pipe"))
}

func TestConnect_TLSFallback_ServerSelectionError(t *testing.T) {
	runFallbackResetTest(t, errors.New("server selection error: context deadline exceeded"))
}
