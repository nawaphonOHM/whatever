# MongoDB

The `mongodb` package is the public entrypoint for connecting consuming services
to MongoDB. It follows the same modular library pattern as `pkg/rest`:
applications interact with clean public APIs, while configuration and implementation
details remain internal to this library.

## Configuration

The client loads and validates configuration from environment variables when it
is started. Configuration values are intentionally not exposed as public
structs or functional options, so importing applications use the standardized
`OHM9996_MONGODB_*` environment variables.

| Environment variable | Default |
| --- | --- |
| `OHM9996_MONGODB_URI` | `mongodb://localhost:27017` |
| `OHM9996_MONGODB_DATABASE` | empty |
| `OHM9996_MONGODB_CONNECT_TIMEOUT` | `10s` |
| `OHM9996_MONGODB_SERVER_SELECTION_TIMEOUT` | `5s` |
| `OHM9996_MONGODB_SOCKET_TIMEOUT` | `10s` |
| `OHM9996_MONGODB_MAX_POOL_SIZE` | `100` |
| `OHM9996_MONGODB_MIN_POOL_SIZE` | `5` |
| `OHM9996_MONGODB_MAX_CONN_IDLE_TIME` | `10m` |
| `OHM9996_MONGODB_APP_NAME` | empty |

The environment-backed configuration, defaults, validation, and driver option
construction live in `internal/mongodb`. They are not part of the API available
to importing projects.

## Usage

Example usage:

```go
client, err := mongodb.Connect(ctx)
if err != nil {
	return err
}
defer client.Disconnect(ctx)

users := client.Collection("users")
// Execute native mongo-driver queries directly on users (*mongo.Collection).
```

The public client exposes `Database`, `Collection`, `RawClient`, `Ping`,
and `Disconnect` methods. `Ping` can be used by readiness probes to verify
connectivity.
