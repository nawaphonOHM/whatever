package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestNewClient_Accessors(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	require.NotNil(t, rawClient)
	defer func() {
		_ = rawClient.Disconnect(context.Background())
	}()

	client := NewClient(rawClient, "default_db")
	require.NotNil(t, client)

	// Database with default name
	dbDefault := client.Database()
	require.NotNil(t, dbDefault)
	assert.Equal(t, "default_db", dbDefault.Name())

	// Database with empty string fallback
	dbEmpty := client.Database("")
	require.NotNil(t, dbEmpty)
	assert.Equal(t, "default_db", dbEmpty.Name())

	// Database with custom name
	dbCustom := client.Database("custom_db")
	require.NotNil(t, dbCustom)
	assert.Equal(t, "custom_db", dbCustom.Name())

	// Collection with default db
	collDefault := client.Collection("users")
	require.NotNil(t, collDefault)
	assert.Equal(t, "users", collDefault.Name())
	assert.Equal(t, "default_db", collDefault.Database().Name())

	// Collection with custom db
	collCustom := client.Collection("orders", "custom_db")
	require.NotNil(t, collCustom)
	assert.Equal(t, "orders", collCustom.Name())
	assert.Equal(t, "custom_db", collCustom.Database().Name())

	// Collection with empty custom db (fallback to default)
	collEmptyDB := client.Collection("products", "")
	require.NotNil(t, collEmptyDB)
	assert.Equal(t, "products", collEmptyDB.Name())
	assert.Equal(t, "default_db", collEmptyDB.Database().Name())

	// RawClient escape hatch
	assert.Same(t, rawClient, client.RawClient())
}

func TestNewClient_EmptyDefaultDatabase_Accessors(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	require.NotNil(t, rawClient)
	defer func() {
		_ = rawClient.Disconnect(context.Background())
	}()

	client := NewClient(rawClient, "")
	require.NotNil(t, client)

	// Database with empty default name returns database handle with empty name
	dbDefault := client.Database()
	require.NotNil(t, dbDefault)
	assert.Equal(t, "", dbDefault.Name())

	dbEmpty := client.Database("")
	require.NotNil(t, dbEmpty)
	assert.Equal(t, "", dbEmpty.Name())

	dbCustom := client.Database("custom_db")
	require.NotNil(t, dbCustom)
	assert.Equal(t, "custom_db", dbCustom.Name())

	// Collection with empty default database
	collDefault := client.Collection("users")
	require.NotNil(t, collDefault)
	assert.Equal(t, "users", collDefault.Name())
	assert.Equal(t, "", collDefault.Database().Name())

	collEmptyDB := client.Collection("users", "")
	require.NotNil(t, collEmptyDB)
	assert.Equal(t, "users", collEmptyDB.Name())
	assert.Equal(t, "", collEmptyDB.Database().Name())

	collCustomDB := client.Collection("users", "custom_db")
	require.NotNil(t, collCustomDB)
	assert.Equal(t, "users", collCustomDB.Name())
	assert.Equal(t, "custom_db", collCustomDB.Database().Name())

	collEmptyName := client.Collection("", "custom_db")
	require.NotNil(t, collEmptyName)
	assert.Equal(t, "", collEmptyName.Name())
	assert.Equal(t, "custom_db", collEmptyName.Database().Name())
}

func TestClient_NilSafety(t *testing.T) {
	ctx := context.Background()

	var nilClient *Client
	assert.Nil(t, nilClient.Database())
	assert.Nil(t, nilClient.Database("test"))
	assert.Nil(t, nilClient.Collection("users"))
	assert.Nil(t, nilClient.Collection("users", "test"))
	assert.Nil(t, nilClient.RawClient())
	assert.ErrorIs(t, nilClient.Ping(ctx), ErrNilClient)
	assert.ErrorIs(t, nilClient.Disconnect(ctx), ErrNilClient)

	uninitializedClient := &Client{rawClient: nil, defaultDatabase: "default_db"}
	assert.Nil(t, uninitializedClient.Database())
	assert.Nil(t, uninitializedClient.Database("test"))
	assert.Nil(t, uninitializedClient.Collection("users"))
	assert.Nil(t, uninitializedClient.Collection("users", "test"))
	assert.Nil(t, uninitializedClient.RawClient())
	assert.ErrorIs(t, uninitializedClient.Ping(ctx), ErrNilClient)
	assert.ErrorIs(t, uninitializedClient.Disconnect(ctx), ErrNilClient)
}

func TestClient_Disconnect(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	client := NewClient(rawClient, "test_db")
	require.NotNil(t, client)

	ctx := context.Background()
	err = client.Disconnect(ctx)
	assert.NoError(t, err)

	// Second disconnect returns an error because client is already disconnected
	err = client.Disconnect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "disconnected")
}

func TestClient_Ping_CanceledContext(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer func() {
		_ = rawClient.Disconnect(context.Background())
	}()

	client := NewClient(rawClient, "test_db")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = client.Ping(ctx)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestClient_Ping_Unreachable(t *testing.T) {
	// Use 127.0.0.1 on an unlikely port with a very short timeout
	rawClient, err := mongo.Connect(
		options.Client().
			ApplyURI("mongodb://127.0.0.1:59999").
			SetServerSelectionTimeout(50 * time.Millisecond).
			SetTimeout(50 * time.Millisecond),
	)
	require.NoError(t, err)
	defer func() {
		_ = rawClient.Disconnect(context.Background())
	}()

	client := NewClient(rawClient, "test_db")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = client.Ping(ctx)
	assert.Error(t, err)
}

func TestConnectWithConfig_NilConfig(t *testing.T) {
	ctx := context.Background()
	client, err := ConnectWithConfig(ctx, nil)
	assert.ErrorIs(t, err, ErrNilConfig)
	assert.Nil(t, client)
}

func TestConnectWithConfig_InvalidConfig(t *testing.T) {
	ctx := context.Background()
	invalidCfg := &Config{
		URI: "",
	}
	client, err := ConnectWithConfig(ctx, invalidCfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid mongodb config")
	assert.Nil(t, client)
}

func TestConnectWithConfig_InvalidURIFormat(t *testing.T) {
	ctx := context.Background()
	invalidURICfg := &Config{
		URI:            "://invalid uri",
		ConnectTimeout: 1 * time.Second,
	}
	client, err := ConnectWithConfig(ctx, invalidURICfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create mongodb client")
	assert.Nil(t, client)
}

func TestConnectWithConfig_PingFailure(t *testing.T) {
	cfg := &Config{
		URI:                    "mongodb://127.0.0.1:59999",
		Database:               "test_db",
		ConnectTimeout:         50 * time.Millisecond,
		ServerSelectionTimeout: 50 * time.Millisecond,
		SocketTimeout:          50 * time.Millisecond,
		MaxPoolSize:            10,
		MinPoolSize:            1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	client, err := ConnectWithConfig(ctx, cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}

func TestConnectWithConfig_WithOptions(t *testing.T) {
	cfg := &Config{
		URI:                    "mongodb://127.0.0.1:59999",
		Database:               "test_db",
		ConnectTimeout:         50 * time.Millisecond,
		ServerSelectionTimeout: 50 * time.Millisecond,
		SocketTimeout:          50 * time.Millisecond,
	}

	extraOpt := options.Client().SetAppName("custom-connect-app")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	client, err := ConnectWithConfig(ctx, cfg, WithDriverOptions(extraOpt))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}

func TestConnect_LoadConfigFailure(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_CONNECT_TIMEOUT", "-5s")

	ctx := context.Background()
	client, err := Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connect timeout cannot be negative")
	assert.Nil(t, client)
}

func TestConnect_PingFailure(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_URI", "mongodb://127.0.0.1:59999")
	t.Setenv("OHM9969_MONGODB_CONNECT_TIMEOUT", "50ms")
	t.Setenv("OHM9969_MONGODB_SERVER_SELECTION_TIMEOUT", "50ms")
	t.Setenv("OHM9969_MONGODB_SOCKET_TIMEOUT", "50ms")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	client, err := Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}

func TestConnect_CanceledContext(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_URI", "mongodb://localhost:27017")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client, err := Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}
