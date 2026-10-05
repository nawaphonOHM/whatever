package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupConnectFirestoreEnv(t *testing.T) {
	t.Setenv(envMongoHost, "demo-app.firestore.goog")
	t.Setenv(envMongoPort, "27017")
	t.Setenv(envMongoUsername, testConnUser)
	t.Setenv(envMongoPassword, testConnPass)
	t.Setenv(envMongoDatabase, "firestore_db")
}

func setupMockFirestoreSuccess() (*int, func()) {
	var firestorePingCalled int
	cleanupPing := SetMockPing(func(context.Context, *mongo.Client) error {
		return nil
	})
	cleanupFirestore := SetMockPingFirestore(func(context.Context, *mongo.Client, string) error {
		firestorePingCalled++
		return nil
	})
	cleanupProbe := SetMockProbe(func(context.Context, *Client, string) error {
		return nil
	})
	return &firestorePingCalled, func() {
		cleanupPing()
		cleanupFirestore()
		cleanupProbe()
	}
}

func setupMockFirestoreFail(pingErr, probeErr error) func() {
	cleanupPing := SetMockPingFirestore(func(context.Context, *mongo.Client, string) error {
		return pingErr
	})
	cleanupProbe := SetMockProbe(func(context.Context, *Client, string) error {
		return probeErr
	})
	return func() {
		cleanupPing()
		cleanupProbe()
	}
}

// TestConnect_Firestore_Success tests successful connection initialization targeting Firestore.
func TestConnect_Firestore_Success(t *testing.T) {
	setupConnectFirestoreEnv(t)
	firestorePingCalled, cleanup := setupMockFirestoreSuccess()
	defer cleanup()

	client, err := Connect(context.Background())
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.True(t, client.IsFirestore())
	assert.Equal(t, "firestore_db", client.Database().Name())
	assert.Equal(t, 1, *firestorePingCalled)
}

// TestConnect_Firestore_PingFailure tests connection failure when Firestore ping fails.
func TestConnect_Firestore_PingFailure(t *testing.T) {
	setupConnectFirestoreEnv(t)
	exitCalled, cleanupExit := setupExitCapture()
	defer cleanupExit()

	cleanup := setupMockFirestoreFail(errors.New("simulated firestore ping failed"), nil)
	defer cleanup()

	client, err := Connect(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}

// TestConnect_Firestore_ProbeFailure tests connection failure when probe fails under Firestore.
func TestConnect_Firestore_ProbeFailure(t *testing.T) {
	setupConnectFirestoreEnv(t)
	exitCalled, cleanupExit := setupExitCapture()
	defer cleanupExit()

	cleanup := setupMockFirestoreFail(nil, errors.New("simulated firestore probe failed"))
	defer cleanup()

	client, err := Connect(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to probe mongodb")
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}
