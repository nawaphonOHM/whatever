# MongoDB Testing Toolkit (`pkg/testing/mongodb`)

The `pkg/testing/mongodb` package provides a dedicated, highly ergonomic MongoDB client and fixture management toolkit tailored specifically for automated test suites, integration tests, and test fixture isolation.

Unlike the production `pkg/mongodb` package—which strictly requires `OHM9996_MONGODB_*` environment variables and initiates process termination on fatal connection errors—`pkg/testing/mongodb` is built for maximum developer utility, zero global environment mutations, explicit Go error returns, and deterministic test state isolation.

---

## Architecture & Package Roles

| Package | Primary Purpose | Configuration Source | Error / Failure Mode | Cleanup / Teardown |
|---|---|---|---|---|
| **`pkg/mongodb`** | Production microservice runtime client | `OHM9996_MONGODB_*` environment variables exclusively | Logs error and initiates peaceful process exit (`os.Exit(0)`) | `defer client.Disconnect(ctx)` in application main |
| **`pkg/testcontainers/mongodb`** | Ephemeral Docker container lifecycle management | Functional container options (`WithImage`, `WithReplicaSet`, etc.) | Returns Go `error` on container startup failure | `defer container.Terminate(ctx)` |
| **`pkg/testing/mongodb`** | Test client connections, fixture cleanup, & `testing.TB` automation | Connection strings (URIs) or programmatic functional options | Returns Go `error` or invokes `tb.Fatalf` with automatic `tb.Cleanup` | Auto-registers `client.Close` to `tb.Cleanup` |

---

## Key Features

- **Direct URI Connection (`ConnectURI`, `NewTestClientURI`)**: Connect immediately from connection strings (such as `container.ConnectionString(ctx)`) with automatic database resolution and query parameter parsing.
- **Programmatic Functional Options (`Connect`, `NewTestClient`)**: Configure host, port, database, credentials, timeouts, connection pooling, and driver options without touching environment variables.
- **Automated `testing.TB` Lifecycle Integration**: `NewTestClient(t, ...)` and `NewTestClientURI(t, uri, ...)` assert connectivity via `tb.Fatalf` and automatically attach connection teardown to `tb.Cleanup`, guaranteeing zero leaked sockets even on early test failures.
- **Zero Process Termination**: Guarantees zero `os.Exit` calls; all connection and validation failures are returned as standard Go errors.
- **Rich Fixture Management**: Built-in helpers to drop databases (`DropDatabase`), drop collections (`DropCollection`), truncate documents (`TruncateCollections`), and list non-system collections (`ListCollectionNames`).
- **Official Driver Compatibility**: Native `Database(name...)` and `Collection(name, dbName...)` accessors returning official MongoDB Go driver v2 handles (`*mongo.Database`, `*mongo.Collection`), plus raw client access (`RawClient()`).
- **Thread Safety**: All `TestClient` operations are safe for concurrent execution across parallel subtests (`t.Parallel()`).

---

## Quick Start

### 1. `testing.TB` Integration with Testcontainers

The most common workflow combines `pkg/testcontainers/mongodb` with `pkg/testing/mongodb`. `NewTestClientURI` connects to the running container and registers automatic teardown with `t.Cleanup`:

```go
package repository_test

import (
	"context"
	"testing"

	tcmongo "github.com/nawaphonOHM/whatever/pkg/testcontainers/mongodb"
	testmongo "github.com/nawaphonOHM/whatever/pkg/testing/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestUserRepository(t *testing.T) {
	ctx := context.Background()

	// 1. Spin up ephemeral MongoDB container
	container, err := tcmongo.Run(ctx, tcmongo.WithDatabase("test_db"))
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to resolve connection string: %v", err)
	}

	// 2. Initialize test client with automated cleanup
	client := testmongo.NewTestClientURI(t, connStr, testmongo.WithDatabase("test_db"))

	// 3. Clean state before running assertions
	if err := client.TruncateCollections(ctx); err != nil {
		t.Fatalf("Failed to truncate collections: %v", err)
	}

	// 4. Perform database operations
	users := client.Collection("users")
	_, err = users.InsertOne(ctx, bson.M{"username": "johndoe", "email": "john@example.com"})
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}
}
```

---

## Usage Patterns

### 2. Programmatic Construction with Functional Options

Configure endpoints, credentials, timeouts, and connection pool sizing programmatically without mutating global environment variables:

```go
func TestCustomEndpoint(t *testing.T) {
	client := testmongo.NewTestClient(t,
		testmongo.WithHost("127.0.0.1"),
		testmongo.WithPort(27017),
		testmongo.WithDatabase("test_catalog"),
		testmongo.WithUsername("testuser"),
		testmongo.WithPassword("testpass"),
		testmongo.WithAuthSource("admin"),
		testmongo.WithConnectTimeout(5*time.Second),
		testmongo.WithPoolLimits(2, 20),
	)

	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}
```

### 3. Standalone Client without `testing.TB` (`Connect` / `ConnectURI`)

For setup scripts, benchmark harnesses, or test suites managing custom teardown routines, use `Connect` or `ConnectURI` directly:

```go
func TestStandaloneConnection(ctx context.Context, uri string) error {
	client, err := testmongo.ConnectURI(ctx, uri,
		testmongo.WithDatabase("analytics"),
		testmongo.WithPing(true),
	)
	if err != nil {
		return err
	}
	defer client.Close(ctx)

	// Work with the client
	coll := client.Collection("events")
	_ = coll
	return nil
}
```

### 4. Test Fixture Isolation & Cleanup

Ensure each test starts with a clean database state using fixture helpers:

```go
func TestOrderProcessing(t *testing.T) {
	ctx := context.Background()
	client := testmongo.NewTestClientURI(t, "mongodb://localhost:27017/test_orders")

	// Reset state between test runs:
	// Option A: Truncate specific collections
	_ = client.TruncateCollections(ctx, "orders", "line_items")

	// Option B: Truncate ALL non-system collections in the default database
	_ = client.TruncateCollections(ctx)

	// Option C: Drop a specific collection
	_ = client.DropCollection(ctx, "temp_data")

	// Option D: Drop the entire default or specified database
	_ = client.DropDatabase(ctx)
	_ = client.DropDatabase(ctx, "other_db")

	// Inspect existing non-system collections
	collections, err := client.ListCollectionNames(ctx)
	if err != nil {
		t.Fatalf("Failed to list collections: %v", err)
	}
	t.Logf("Active collections: %v", collections)
}
```

### 5. Parallel Subtests with Isolated Databases

Run subtests concurrently by assigning unique database names per subtest:

```go
func TestParallelOperations(t *testing.T) {
	connURI := "mongodb://localhost:27017"

	tests := []struct {
		name string
		db   string
	}{
		{name: "ScenarioA", db: "test_parallel_a"},
		{name: "ScenarioB", db: "test_parallel_b"},
		{name: "ScenarioC", db: "test_parallel_c"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			client := testmongo.NewTestClientURI(t, connURI, testmongo.WithDatabase(tt.db))
			defer func() { _ = client.DropDatabase(ctx, tt.db) }()

			// Execute test logic isolated to tt.db
			coll := client.Collection("items")
			_, err := coll.InsertOne(ctx, bson.M{"test": tt.name})
			if err != nil {
				t.Fatalf("Insert failed: %v", err)
			}
		})
	}
}
```

---

## API Reference

### Constructors & Helpers

| Function | Signature | Description |
|---|---|---|
| `Connect` | `Connect(ctx context.Context, opts ...Option) (*TestClient, error)` | Connects via functional options; returns standard error on failure (no process exit) |
| `ConnectURI` | `ConnectURI(ctx context.Context, uri string, opts ...Option) (*TestClient, error)` | Connects via URI string with optional overrides; returns standard error on failure |
| `NewTestClient` | `NewTestClient(tb testing.TB, opts ...Option) *TestClient` | Connects via functional options, asserts connectivity (`tb.Fatalf`), and binds `client.Close` to `tb.Cleanup` |
| `NewTestClientURI` | `NewTestClientURI(tb testing.TB, uri string, opts ...Option) *TestClient` | Connects via URI string, asserts connectivity (`tb.Fatalf`), and binds `client.Close` to `tb.Cleanup` |
| `NewClient` | `NewClient(rawClient *mongo.Client, defaultDatabase string) *TestClient` | Wraps an existing official `*mongo.Client` with `TestClient` fixture utilities |
| `DefaultOptions` | `DefaultOptions() *Options` | Returns an `Options` struct initialized with recommended test defaults |
| `NewOptions` | `NewOptions(opts ...Option) *Options` | Evaluates functional options over test defaults and returns the resulting `*Options` |

### Functional Options

| Option | Signature | Description | Default |
|---|---|---|---|
| `WithURI` | `WithURI(uri string) Option` | Direct connection URI string | `""` |
| `WithHost` | `WithHost(host string) Option` | Target MongoDB hostname or IP | `"localhost"` |
| `WithPort` | `WithPort(port int) Option` | Target MongoDB network port | `27017` |
| `WithProtocol` | `WithProtocol(protocol string) Option` | Connection scheme (`mongodb` or `mongodb+srv`) | `"mongodb"` |
| `WithDatabase` | `WithDatabase(database string) Option` | Default database name for collections & fixtures | `""` |
| `WithUsername` | `WithUsername(username string) Option` | Authentication username | `""` |
| `WithPassword` | `WithPassword(password string) Option` | Authentication password | `""` |
| `WithAuthSource` | `WithAuthSource(authSource string) Option` | Authentication database name (e.g., `admin`) | `""` |
| `WithAppName` | `WithAppName(appName string) Option` | Client application name for diagnostics | `""` |
| `WithConnectTimeout` | `WithConnectTimeout(timeout time.Duration) Option` | Initial TCP connection establishment timeout | `10s` |
| `WithServerSelectionTimeout` | `WithServerSelectionTimeout(timeout time.Duration) Option` | Cluster primary / server election timeout | `5s` |
| `WithSocketTimeout` | `WithSocketTimeout(timeout time.Duration) Option` | Socket read/write operation timeout | `10s` |
| `WithMaxConnIdleTime` | `WithMaxConnIdleTime(idleTime time.Duration) Option` | Maximum idle duration for pooled connections | `10m` |
| `WithPoolLimits` | `WithPoolLimits(min, max uint64) Option` | Minimum and maximum connection pool sizes | Min: `5`, Max: `100` |
| `WithPing` | `WithPing(enablePing bool) Option` | Enable/disable connectivity verification on connect | `true` |
| `WithTLS` | `WithTLS(enableTLS bool) Option` | Enable/disable TLS/SSL encryption | `false` |
| `WithDirectConnection` | `WithDirectConnection(direct bool) Option` | Force direct connection to a single host instance | `false` |
| `WithUUIDRepresentation` | `WithUUIDRepresentation(repr string) Option` | Binary UUID representation format | `"unspecified"` |
| `WithDriverOptions` | `WithDriverOptions(opts ...*options.ClientOptions) Option` | Appends raw official mongo-driver `ClientOptions` | `nil` |

### `TestClient` Methods

| Method | Signature | Description |
|---|---|---|
| `Database` | `Database(name ...string) *mongo.Database` | Returns `*mongo.Database` for the specified override or default database |
| `Collection` | `Collection(name string, dbName ...string) *mongo.Collection` | Returns `*mongo.Collection` for the given collection in target/default database |
| `Ping` | `Ping(ctx context.Context) error` | Verifies cluster reachability with `readpref.Primary()` |
| `Disconnect` | `Disconnect(ctx context.Context) error` | Gracefully closes all connections in the pool |
| `Close` | `Close(ctx ...context.Context) error` | Gracefully closes the client (uses default `10s` timeout if context is omitted) |
| `RawClient` | `RawClient() *mongo.Client` | Returns the underlying official driver `*mongo.Client` handle |
| `Client` | `Client() *client.Client` | Returns the underlying internal client handle |
| `DropDatabase` | `DropDatabase(ctx context.Context, name ...string) error` | Drops the specified or default database |
| `DropCollection` | `DropCollection(ctx context.Context, collection string, dbName ...string) error` | Drops a specific collection from the resolved database |
| `TruncateCollections` | `TruncateCollections(ctx context.Context, collections ...string) error` | Deletes all documents from specified collections or all non-system collections |
| `ListCollectionNames` | `ListCollectionNames(ctx context.Context, dbName ...string) ([]string, error)` | Returns all non-system collection names in the specified or default database |

### Sentinel Errors

| Error Sentinel | String Value | Cause |
|---|---|---|
| `ErrNilClient` | `"mongodb test client is nil"` | Invoking methods on a `nil` `TestClient` receiver |
| `ErrNilConfig` | `"mongodb config cannot be nil"` | Providing `nil` options configuration |
| `ErrEmptyURI` | `"mongodb uri cannot be empty"` | Passing an empty URI to `ConnectURI` or `NewTestClientURI` |
| `ErrInvalidPort` | `"mongodb port must be between 1 and 65535"` | Specifying a port number outside `1..65535` under `mongodb://` |
| `ErrNilTestingTB` | `"testing.TB cannot be nil"` | Passing `nil` `testing.TB` to `NewTestClient` or `NewTestClientURI` |

### Default Constants

| Constant | Value | Description |
|---|---|---|
| `DefaultHost` | `"localhost"` | Default MongoDB host address |
| `DefaultPort` | `27017` | Default MongoDB network port |
| `DefaultProtocol` | `"mongodb"` | Default MongoDB connection scheme |
| `DefaultConnectTimeout` | `10 * time.Second` | Default connection establishment timeout |
| `DefaultServerSelectionTimeout` | `5 * time.Second` | Default server selection timeout |
| `DefaultSocketTimeout` | `10 * time.Second` | Default socket operation timeout |
| `DefaultMaxConnIdleTime` | `10 * time.Minute` | Default maximum connection idle duration |
| `DefaultMaxPoolSize` | `100` | Default maximum connection pool capacity |
| `DefaultMinPoolSize` | `5` | Default minimum connection pool capacity |
| `DefaultUUIDRepresentation` | `"unspecified"` | Default BSON UUID binary representation |
