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

| Environment variable | Description | Required / Default |
| --- | --- | --- |
| `OHM9996_MONGODB_HOST` | Hostname or IP address of the MongoDB server | Required |
| `OHM9996_MONGODB_PORT` | Network port (1–65535, optional for `mongodb+srv`) | `27017` |
| `OHM9996_MONGODB_PROTOCOL` | Connection protocol (`mongodb` or `mongodb+srv`) | `mongodb` |
| `OHM9996_MONGODB_DATABASE` | Default application database for direct collection access | `""` (empty) |
| `OHM9996_MONGODB_USERNAME` | Username for authentication (optional, paired with password) | `""` (empty) |
| `OHM9996_MONGODB_PASSWORD` | Password for authentication (optional, paired with username) | `""` (empty) |
| `OHM9996_MONGODB_AUTH_SOURCE` | Authentication database name (e.g., `admin`) | `""` (empty) |
| `OHM9996_MONGODB_APP_NAME` | Application name for connection metadata and diagnostics | `""` (empty) |
| `OHM9996_MONGODB_UUID_REPRESENTATION` | UUID binary representation (`unspecified`, `standard`, `csharpLegacy`, `javaLegacy`, `pythonLegacy`) | `unspecified` |
| `OHM9996_MONGODB_CONNECT_TIMEOUT` | Maximum duration for initial TCP connection establishment | `10s` |
| `OHM9996_MONGODB_SERVER_SELECTION_TIMEOUT` | Timeout for cluster server discovery and primary election | `5s` |
| `OHM9996_MONGODB_SOCKET_TIMEOUT` | Socket read and write operation timeout | `10s` |
| `OHM9996_MONGODB_MAX_CONN_IDLE_TIME` | Maximum idle duration before an unused connection is closed | `10m` |
| `OHM9996_MONGODB_MAX_POOL_SIZE` | Maximum number of concurrent connections in the pool | `100` |
| `OHM9996_MONGODB_MIN_POOL_SIZE` | Minimum number of idle connections maintained in the pool | `5` |

If mandatory configuration (`HOST`) is missing, or if connection fails after TLS fallback, the client logs a descriptive error and triggers peaceful termination (`exitFunc(0)`).

The environment-backed configuration, defaults, validation, and driver option
construction live in `internal/mongodb/config` and `internal/mongodb/client`. They are not part of the API available
to importing projects.

## Connection Lifecycle & TLS Fallback

Calling `mongodb.Connect(ctx)` executes a two-phase connection flow:
1. **Unencrypted Connection Attempt**: First attempts to connect and ping the MongoDB server without TLS (`tls=false`).
2. **Automatic TLS Fallback**: If the server rejects the unencrypted connection indicating TLS/SSL is required (e.g., MongoDB Atlas or secured clusters), the client automatically retries connection with TLS enabled (`tls=true`).
3. **Graceful Termination on Failure**: If connectivity cannot be established after retry or due to fatal configuration errors, the client logs the error, triggers graceful program exit, and returns the underlying error.

## Usage

Example usage:

```go
client, err := mongodb.Connect(ctx)
if err != nil {
	return err
}
defer client.Disconnect(ctx)

// Access collections using the default configured database:
users := client.Collection("users")

// Or explicitly specify a database override:
orders := client.Collection("orders", "custom_db")
```

The public client exposes `Database`, `Collection`, `RawClient`, `Ping`,
and `Disconnect` methods. `Ping` can be used by readiness probes to verify
connectivity. The process exit hook can be intercepted during testing via `SetExitFunc`.
