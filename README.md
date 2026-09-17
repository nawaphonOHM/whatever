# Go Boilerplate Library

[![CI](https://github.com/example/go-boilerplate/actions/workflows/ci.yml/badge.svg)](https://github.com/example/go-boilerplate/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/example/go-boilerplate)](https://goreportcard.com/report/github.com/example/go-boilerplate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A production-ready, modular Go library designed to be imported by microservices and API applications. It provides declarative REST API registration scanning, pre-registered health/readiness endpoints, encapsulated Gin HTTP server lifecycle management with graceful shutdown, production-grade middlewares, and uniform JSON API response envelopes.

---

## Table of Contents

- [Architecture & Directory Layout](#architecture--directory-layout)
- [Public Packages](#public-packages)
- [REST Registration Contract](#rest-registration-contract)
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
.
├── example/                   # Consumer-style registration and usage tests
├── internal/                  # Private configuration, middleware, and health implementation
├── pkg/
│   ├── logger/                # Public structured logger helpers
│   ├── response/              # JSON response envelopes and helpers
│   └── server/                # StartREST, declarative registration contract
├── .github/workflows/ci.yml   # Library test, lint, and build verification
├── Makefile                   # Local verification commands
└── README.md
```

## Public Packages

### `pkg/rest/server`

This is the primary package. `server.StartREST([]*server.RestApiRegistration)`
loads the library configuration, installs the default HTTP middleware, registers
the framework health endpoints, validates the supplied registrations, and starts
Gin with graceful shutdown on `SIGINT` or `SIGTERM`.

The registration slice is explicit: Go cannot discover arbitrary packages that
import a module at runtime. Each consuming service collects registrations from its
domain packages and passes them to `StartREST`.

### `pkg/rest/response`

Generic JSON response formatting utilities adhering to a consistent API contract and RFC 9457 Problem Details:

- **Success Envelope**: `{"success": true, "message": "...", "data": ..., "timestamp": "..."}`
- **Error Envelope**: RFC 9457 Problem Details (`application/problem+json`)

```go
import "github.com/nawaphonOHM/whatever/pkg/rest/response"

// HTTP 200 OK with data and optional message (returns response.Response)
return response.OK(data)
return response.OK(data, "Fetched successfully")

// HTTP 201 Created with data and optional message (returns response.Response)
return response.Created(newResource)
return response.Created(newResource, "Created successfully")

// HTTP 204 No Content
return response.NoContent()

// HTTP 400 Bad Request (returns RFC 9457 ProblemDetails response.Response)
return response.BadRequest("VALIDATION_FAILED", "Invalid payload", validationDetails)

// HTTP 404 Not Found (returns RFC 9457 ProblemDetails response.Response)
return response.NotFound("NOT_FOUND", "Resource not found")
```

### `pkg/logger`

The structured logger package provides logging middleware and request-tracing
utilities with slog support. Request-ID generation, panic recovery, CORS, and
framework health probe handling are installed internally by `pkg/rest/server`.

---

## REST Registration Contract

External projects define endpoints declaratively using `RestApiRegistration` and `ExportableApi`:

```go
package myfeature

import (
    "github.com/example/go-boilerplate/pkg/response"
    "github.com/example/go-boilerplate/pkg/server"
)

func NewFeatureAPIs() *server.RestApiRegistration {
    return &server.RestApiRegistration{
        Version: 1,           // Generates /v1 prefix
        Prefix:  "/items",     // Base path for this group; no /api is added
        Apis: []*server.ExportableApi{
            {
                Path:   "",
                Method: server.GET,
                Handler: func(c *server.Context) response.Response {
                    return response.OK([]string{"item1", "item2"})
                },
            },
            {
                Path:   "/:id",
                Method: server.GET,
                Handler: func(c *server.Context) response.Response {
                    id := c.Param("id")
                    return response.OK(map[string]string{"id": id})
                },
            },
        },
    }
}
```

### Route Validation Guarantees
When `StartREST` is called, the library performs fail-fast preflight validation before mutating Gin:
1. **Reserved Path Check**: Rejects any registration mapping to `/health` or `/ready`.
2. **Duplicate Route Prevention**: Detects duplicate method + path combinations and returns a descriptive error rather than allowing Gin to panic.
3. **Nil Safety**: Validates that registrations, APIs, and handlers are non-nil.
4. **Method Validation**: Verifies that HTTP methods are valid recognized methods (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`, etc.).

---

## Reserved Framework Endpoints

The library pre-registers and reserves the following endpoints:

| Method | Path | Description | Response Data |
|---|---|---|---|
| `GET` | `/health` | Liveness probe (verifies process is running) | `{"status": "up", "timestamp": "...", "version": "..."}` |
| `GET` | `/ready` | Readiness probe (verifies server is ready for traffic) | `{"status": "ready", "timestamp": "...", "version": "..."}` |

Any attempt by a consuming application to register a route at `/health` or `/ready` is rejected with `server.ErrReservedPath`.

---

## Configuration

All environment variables read by the library use the **`OHM9969_`** prefix:

| Variable | Description | Default |
|---|---|---|
| `OHM9969_SERVER_HOST` | Network interface address to bind | `""` (all interfaces) |
| `OHM9969_SERVER_PORT` | HTTP server port | `8080` |
| `OHM9969_GIN_MODE` | Gin engine mode (`debug`, `release`, `test`) | `release` |
| `OHM9969_SERVER_READ_TIMEOUT` | Maximum duration for reading request | `10s` |
| `OHM9969_SERVER_WRITE_TIMEOUT` | Maximum duration for writing response | `10s` |
| `OHM9969_SERVER_IDLE_TIMEOUT` | Maximum duration for keep-alive connections | `60s` |
| `OHM9969_SERVER_SHUTDOWN_TIMEOUT` | Graceful shutdown timeout before forcing exit | `10s` |

`StartREST` loads an optional `.env` file when present. Configuration precedence
is system environment, then `.env`, then the defaults declared by the library.
An absent `.env` file is not an error.

---

## Consumer Bootstrap

Here is how an importing application bootstraps a service using this library:

```go
package main

import (
    "log"

    "github.com/example/go-boilerplate/pkg/server"
    "github.com/myorg/myapp/internal/items"
)

func main() {
    // Collect API registrations from domain modules. StartREST loads the
    // OHM9969_* configuration and installs the framework defaults.
    registrations := []*server.RestApiRegistration{
        items.NewItemAPIRegistration(),
    }

    if err := server.StartREST(registrations); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

---

## Development Workflows

### Makefile Targets

```bash
$ make help

Usage:
  make <target>

Targets:
  help                Display this help screen
  build               Verify compilation of all packages
  test                Run unit and integration tests with race detection
  test-coverage       Run tests with race detection and HTML coverage report
  vet                 Run go vet analysis
  lint                Run golangci-lint
  tidy                Tidy and verify Go module dependencies
  clean               Clean temporary test coverage and artifact files
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
