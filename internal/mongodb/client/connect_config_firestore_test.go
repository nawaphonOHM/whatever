package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
)

const (
	testFirestorePort = 27017
)

func createFirestoreTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Host = "test-project.firestore.goog"
	cfg.Port = testFirestorePort
	cfg.Username = testUnreachUser
	cfg.Password = testUnreachPass
	cfg.Database = "firestore_config_db"
	cfg.Protocol = config.ProtocolMongoDB
	cfg.UUIDRepresentation = config.UUIDRepresentationUnspecified
	return cfg
}

func setupMockPingFirestoreFailure(err error) func() {
	cleanupFirestore := SetMockPingFirestore(func(context.Context, *mongo.Client, string) error {
		return err
	})
	cleanupProbe := SetMockProbe(func(context.Context, *Client, string) error {
		return nil
	})
	return func() {
		cleanupFirestore()
		cleanupProbe()
	}
}

func setupMockFirestoreConfigSuccess() (*bool, func()) {
	var called bool
	cleanupPing := SetMockPing(func(context.Context, *mongo.Client) error {
		return nil
	})
	cleanupFirestore := SetMockPingFirestore(func(ctx context.Context, rc *mongo.Client, dbName string) error {
		called = ctx != nil && rc != nil && dbName == "firestore_config_db"
		return nil
	})
	cleanupProbe := SetMockProbe(func(context.Context, *Client, string) error {
		return nil
	})
	return &called, func() {
		cleanupPing()
		cleanupFirestore()
		cleanupProbe()
	}
}

func TestConnectWithConfig_Firestore_Success(t *testing.T) {
	called, cleanup := setupMockFirestoreConfigSuccess()
	defer cleanup()

	client, err := ConnectWithConfig(context.Background(), createFirestoreTestConfig())
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.True(t, client.IsFirestore())
	assert.True(t, *called)
}

func TestConnectWithConfig_Firestore_PingFailure(t *testing.T) {
	exitCalled, restore := interceptExitHook()
	defer restore()

	restoreMock := setupMockPingFirestoreFailure(errors.New("simulated firestore ping failure"))
	defer restoreMock()

	client, err := ConnectWithConfig(context.Background(), createFirestoreTestConfig())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}
