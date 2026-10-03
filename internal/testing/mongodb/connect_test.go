package mongodb

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const testInvalidPort = 99999

func setMockPingError(t *testing.T, expectedErr error) {
	origPing := pingClient
	t.Cleanup(func() { pingClient = origPing })
	pingClient = func(context.Context, *mongo.Client) error {
		return expectedErr
	}
}

func setMockPingSuccess(t *testing.T, pingCalled *bool) {
	origPing := pingClient
	t.Cleanup(func() { pingClient = origPing })
	pingClient = func(context.Context, *mongo.Client) error {
		*pingCalled = true
		return nil
	}
}

func TestConnect_InvalidPort(t *testing.T) {
	ctx := context.Background()
	client, err := Connect(ctx, WithPort(testInvalidPort))
	require.ErrorIs(t, err, ErrInvalidPort)
	assert.Nil(t, client)
}

func TestConnect_NoPing_Success(t *testing.T) {
	ctx := context.Background()
	client, err := Connect(ctx, WithPing(false), WithDatabase("test_db"))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "test_db", client.Database().Name())
}

func TestConnect_PingFailure(t *testing.T) {
	expectedErr := errors.New("simulated ping error")
	setMockPingError(t, expectedErr)

	ctx := context.Background()
	client, err := Connect(ctx, WithPing(true))
	require.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "simulated ping error")
}

func TestConnect_PingSuccess(t *testing.T) {
	pingCalled := false
	setMockPingSuccess(t, &pingCalled)

	ctx := context.Background()
	client, err := Connect(ctx, WithPing(true), WithDatabase("ping_db"))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.True(t, pingCalled)
	assert.Equal(t, "ping_db", client.Database().Name())
}
