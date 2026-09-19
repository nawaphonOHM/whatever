# whatever

[![CI](https://github.com/nawaphonOHM/whatever/actions/workflows/ci.yml/badge.svg)](https://github.com/nawaphonOHM/whatever/actions/workflows/ci.yml)

A production-ready, modular Go library designed to be imported by microservices and API applications. It provides declarative REST API registration scanning, pre-registered health/readiness endpoints, encapsulated Gin HTTP server lifecycle management with graceful shutdown, zero-boilerplate managed MongoDB client connectivity, production-grade middlewares, and uniform JSON API response envelopes.

---

## Table of Contents

- [Architecture & Directory Layout](#architecture--directory-layout)
- [Public Packages](#public-packages)
- [REST Registration Contract](#rest-registration-contract)
- [MongoDB Client (`pkg/mongodb`)](#mongodb-client-pkgmongodb)
- [Reserved Framework Endpoints](#reserved-framework-endpoints)
- [Configuration](#configuration)
- [Consumer Bootstrap](#consumer-bootstrap)
- [Development Workflows](#development-workflows)
- [Continuous Integration](#continuous-integration)

---

## Architecture & Directory Layout

This repository is a library, not an application binary. An importing service owns
its `main` package and supplies its own domain registrations.

```
├── internal/                  # Private configuration, middleware, server, and mongodb implementation
├── pkg/
│   ├── logger/                # Public structured logger helpers
│   ├── mongodb/               # Public MongoDB connection entrypoint and managed client
│   └── rest/                  # Public REST API registration contracts & JSON response factories
├── .github/workflows/ci.yml   # Library test, lint, and build verification
├── Makefile                   # Local verification commands
└── README.md
```

## Public Packages

### `pkg/rest`

The unified REST package provides declarative route registration contracts and standardized response factories. Server execution, engine bootstrap, and lifecycle management are encapsulated internally within `internal/rest/server`.

- **API Registration**: Declarative route registration structs (`RRestAPIRegistration`, `ExportableAPI`), handler signatures (`Handler`, `Middleware`), request context wrapper (`Context`), HTTP method constants (`GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `HEAD`, `CONNECT`, `TRACE`), and path types (`Pathz`, `APIVersioning`).
- **Standardized Response Factories**: The `Response` interface, generic JSON success envelope builders (`OK`, `Created`, `NoContent`, `JSON`, `SuccessResponse[T]`, `Envelope[T]`), and RFC 9457 Problem Details error constructors (`BadRequest`, `Unauthorized`, `Forbidden`, `NotFound`, `InternalServerError`, `Error`).

- **Success Envelope**: `{"success": true, "message": "...", "data": ..., "timestamp": "..."}`
- **Error Envelope**: RFC 9457 Problem Details (`application/problem+json`)

```go
import "github.com/nawaphonOHM/whatever/pkg/rest"

// HTTP 200 OK with data and optional message (returns rest.Response)
return rest.OK(data)
return rest.OK(data, "Fetched successfully")

// HTTP 201 Created with data and optional message (returns rest.Response)
return rest.Created(newResource)
return rest.Created(newResource, "Created successfully")

// HTTP 204 No Content
return rest.NoContent()

// HTTP 400 Bad Request (returns RFC 9457 ProblemDetails rest.Response)
return rest.BadRequest("VALIDATION_FAILED", "Invalid payload", validationDetails)

// HTTP 404 Not Found (returns RFC 9457 ProblemDetails rest.Response)
return rest.NotFound("NOT_FOUND", "Resource not found")
```

### `pkg/logger`

The structured logger package provides HTTP request logging middleware and request-tracing utilities with `log/slog` support. It exports middleware constructors (`Logger`, `WithLogger`, `WithConfig`), request tracing helper (`GetRequestID`), and configuration struct (`Config`). Request-ID generation, panic recovery, CORS, and framework health probe handling are installed internally by the framework server engine.

### `pkg/mongodb`

The MongoDB package provides a zero-boilerplate entrypoint for connecting microservices to MongoDB clusters. Calling `mongodb.Connect(ctx)` loads configuration automatically from `OHM9969_MONGODB_*` environment variables, connects to the cluster, validates connectivity via ping, and returns a managed `*mongodb.Client`. The client exposes native `*mongo.Database` and `*mongo.Collection` handles for executing queries directly via the official MongoDB Go driver v2 (`go.mongodb.org/mongo-driver/v2`).

---

## REST Registration Contract

External projects define endpoints declaratively using `RRestAPIRegistration` and `ExportableAPI`:

```go
package myfeature

import (
    "github.com/nawaphonOHM/whatever/pkg/rest"
)

func authMiddleware(c rest.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        return
    }
}

func NewFeatureAPIs() *rest.RRestAPIRegistration {
    return &rest.RRestAPIRegistration{
        Version: 1,           // Generates /v1 prefix
        Prefix:  "/items",     // Base path for this group; no /api is added
        Apis: []*rest.ExportableAPI{
            {
                Path:   "",
                Method: rest.GET,
                Handler: func(c rest.Context) rest.Response {
                    return rest.OK([]string{"item1", "item2"})
                },
            },
            {
                Path:       "/:id",
                Method:     rest.GET,
                Middleware: []rest.Middleware{authMiddleware},
                Handler: func(c rest.Context) rest.Response {
                    id := c.Param("id")
                    return rest.OK(map[string]string{"id": id})
                },
            },
        },
    }
}
```

### Route Validation Guarantees
When registrations are evaluated, the library performs fail-fast preflight validation before mutating Gin:
1. **Reserved Path Check**: Rejects any registration mapping to `/health` or `/ready`.
2. **Duplicate Route Prevention**: Detects duplicate method + path combinations and returns a descriptive error rather than allowing Gin to panic.
3. **Nil Safety**: Validates that registrations, APIs, and handlers are non-nil.
4. **Method Validation**: Verifies that HTTP methods are valid recognized methods (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`, etc.).

---

## MongoDB Client (`pkg/mongodb`)

The `pkg/mongodb` package encapsulates MongoDB connection establishment, connection pooling, and lifecycle management while directly exposing official driver `*mongo.Database` and `*mongo.Collection` types for zero-overhead querying.

### Connecting & Lifecycle

Consuming applications connect to MongoDB using `mongodb.Connect(ctx)`. The library automatically reads and validates `OHM9969_MONGODB_*` environment variables, initializes connection pools, and verifies connectivity via an initial ping:

```go
package database

import (
    "context"
    "log"
    "time"

    "github.com/nawaphonOHM/whatever/pkg/mongodb"
)

func InitMongoDB(ctx context.Context) (*mongodb.Client, func()) {
    client, err := mongodb.Connect(ctx)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    cleanup := func() {
        disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := client.Disconnect(disconnectCtx); err != nil {
            log.Printf("Failed to gracefully disconnect MongoDB: %v", err)
        }
    }

    return client, cleanup
}
```

### Database & Collection Handles

Once connected, access collections and databases directly. If no database name is specified, operations automatically use the default database configured in `OHM9969_MONGODB_DATABASE`:

```go
package repository

import (
    "context"

    "github.com/nawaphonOHM/whatever/pkg/mongodb"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
    ID    string `bson:"_id,omitempty"`
    Email string `bson:"email"`
    Name  string `bson:"name"`
}

type UserRepository struct {
    client *mongodb.Client
}

func NewUserRepository(client *mongodb.Client) *UserRepository {
    return &UserRepository{client: client}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    // Automatically targets the default database from OHM9969_MONGODB_DATABASE
    coll := r.client.Collection("users")

    var user User
    err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) RecordAudit(ctx context.Context, entry bson.M) error {
    // Explicitly target a specific database instead of the default database
    coll := r.client.Collection("audit_logs", "audit_db")

    _, err := coll.InsertOne(ctx, entry)
    return err
}
```

### Readiness & Health Verification

Use `Ping(ctx)` to verify active cluster connectivity (useful in custom readiness checks or background health pollers):

```go
// Sends a primary read-preference ping command to the MongoDB cluster
if err := client.Ping(ctx); err != nil {
    log.Printf("MongoDB health check failed: %v", err)
}
```

### Raw Driver Access

When advanced MongoDB features are needed—such as multi-document transactions, client sessions, or change streams—use `RawClient()` to access the underlying official `*mongo.Client`:

```go
rawClient := client.RawClient()
session, err := rawClient.StartSession()
if err != nil {
    return err
}
defer session.EndSession(ctx)
```

---

## Reserved Framework Endpoints

The library pre-registers and reserves the following endpoints:

| Method | Path | Description | Response Data |
|---|---|---|---|
| `GET` | `/health` | Liveness probe (verifies process is running) | `{"status": "up", "timestamp": "...", "version": "..."}` |
| `GET` | `/ready` | Readiness probe (verifies server is ready for traffic) | `{"status": "ready", "timestamp": "...", "version": "..."}` |

Any attempt by a consuming application to register a route at `/health` or `/ready` is rejected during route validation.

---

## Configuration

All environment variables read by the library use the **`OHM9969_`** prefix.

### REST Server Configuration

| Variable | Description | Default |
|---|---|---|
| `OHM9969_SERVER_HOST` | Network interface address to bind | `""` (all interfaces) |
| `OHM9969_SERVER_PORT` | HTTP server port | `8080` |
| `OHM9969_APP_VERSION` | Application version reported by health endpoints | `""` (empty) |
| `OHM9969_GIN_MODE` | Gin engine mode (`debug`, `release`, `test`) | `release` |
| `OHM9969_SERVER_READ_TIMEOUT` | Maximum duration for reading request | `10s` |
| `OHM9969_SERVER_WRITE_TIMEOUT` | Maximum duration for writing response | `10s` |
| `OHM9969_SERVER_IDLE_TIMEOUT` | Maximum duration for keep-alive connections | `60s` |
| `OHM9969_SERVER_SHUTDOWN_TIMEOUT` | Graceful shutdown timeout before forcing exit | `10s` |

### MongoDB Configuration

| Variable | Description | Default |
|---|---|---|
| `OHM9969_MONGODB_URI` | MongoDB connection URI string | `mongodb://localhost:27017` |
| `OHM9969_MONGODB_DATABASE` | Default database name for collections | `""` (empty) |
| `OHM9969_MONGODB_CONNECT_TIMEOUT` | Initial connection timeout | `10s` |
| `OHM9969_MONGODB_SERVER_SELECTION_TIMEOUT` | Server selection timeout | `5s` |
| `OHM9969_MONGODB_SOCKET_TIMEOUT` | Socket read/write timeout | `10s` |
| `OHM9969_MONGODB_MAX_POOL_SIZE` | Maximum connection pool size | `100` |
| `OHM9969_MONGODB_MIN_POOL_SIZE` | Minimum connection pool size | `5` |
| `OHM9969_MONGODB_MAX_CONN_IDLE_TIME` | Maximum duration a connection remains idle | `10m` |
| `OHM9969_MONGODB_APP_NAME` | Client metadata application name sent to MongoDB | `""` (empty) |

Both the REST server engine and `mongodb.Connect` load an optional `.env` file when present. Configuration precedence
is system environment variables, then `.env`, then library defaults. An absent `.env` file is not an error.

---

## Consumer Bootstrap

Here is how an importing application bootstraps MongoDB and defines REST API registrations:

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/nawaphonOHM/whatever/pkg/mongodb"
    "github.com/nawaphonOHM/whatever/pkg/rest"
    "github.com/myorg/myapp/internal/items"
)

func main() {
    ctx := context.Background()

    // 1. Initialize MongoDB connection with zero-boilerplate environment configuration
    mongoClient, err := mongodb.Connect(ctx)
    if err != nil {
        log.Fatalf("Failed to initialize MongoDB: %v", err)
    }
    defer func() {
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := mongoClient.Disconnect(shutdownCtx); err != nil {
            log.Printf("MongoDB disconnect error: %v", err)
        }
    }()

    // 2. Initialize domain items API with MongoDB collection
    itemsColl := mongoClient.Collection("items")
    itemAPI := items.NewItemAPIRegistration(itemsColl)

    // 3. Collect domain API registrations
    registrations := []*rest.RRestAPIRegistration{
        itemAPI,
    }
    _ = registrations
}
```

HTTP server execution, middleware orchestration, and graceful shutdown are managed internally by the framework engine.

---

## Development Workflows

### Makefile Targets

```bash
# Run unit and integration tests with race detection
make test

# Generate HTML code coverage report
make test-coverage

# Run linters (golangci-lint or go vet fallback)
make lint

# Run go vet analysis
make vet

# Verify compilation of all packages
make build

# Tidy and verify module dependencies
make tidy

# Run test, vet, and build
make all

# Clean temporary test and build artifacts
make clean
```

### Testing & Coverage

Execute all package tests with data race detector enabled:
```bash
make test
```

Generate and view coverage breakdown:
```bash
make test-coverage
```

---

## Continuous Integration

The repository includes a GitHub Actions workflow (`.github/workflows/ci.yml`) configured to:
- Enforce Go module integrity (`go mod verify`).
- Execute static analysis with `golangci-lint`.
- Run the full test suite across all packages with `-race` and coverage collection.
- Verify compilation of all library packages with `go build -v ./...`.
