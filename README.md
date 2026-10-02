# whatever

[![CI](https://github.com/nawaphonOHM/whatever/actions/workflows/ci.yml/badge.svg)](https://github.com/nawaphonOHM/whatever/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.27-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/nawaphonOHM/whatever.svg)](https://pkg.go.dev/github.com/nawaphonOHM/whatever)
[![Go Report Card](https://goreportcard.com/badge/github.com/nawaphonOHM/whatever)](https://goreportcard.com/report/github.com/nawaphonOHM/whatever)

A production-ready, modular Go library designed to bootstrap high-performance microservices and RESTful API applications. It provides declarative Blueprint routing with preflight collision validation, encapsulated Gin HTTP server lifecycle management with signal-driven graceful shutdown, standardized RFC 9457 Problem Details error responses, uniform success envelopes, structured logging with `log/slog`, native OpenTelemetry distributed tracing correlation, zero-boilerplate managed MongoDB client connectivity, and isolated Testcontainers-based MongoDB integration testing.

---

## Table of Contents

- [Prerequisites & Installation](#prerequisites--installation)
- [Architecture & Directory Layout](#architecture--directory-layout)
- [Package Catalog & Feature Matrix](#package-catalog--feature-matrix)
  - [`pkg/rest`](#pkgrest)
  - [`pkg/logging`](#pkglogging)
  - [`pkg/logger`](#pkglogger)
  - [`pkg/mongodb`](#pkgmongodb)
  - [`pkg/testcontainers/mongodb`](#pkgtestcontainersmongodb)
- [Quick Start](#quick-start)
- [REST Routing & Blueprint API](#rest-routing--blueprint-api)
  - [Blueprint & Server Middleware Pipeline](#blueprint--server-middleware-pipeline)
  - [Declarative Registration Model](#declarative-registration-model)
  - [HTTP Methods](#http-methods)
  - [API Versioning & Path Resolution](#api-versioning--path-resolution)
  - [CORS & Blueprint Metadata](#cors--blueprint-metadata)
  - [Middleware Pipeline](#middleware-pipeline)
  - [Preflight Route Validation & Safeguards](#preflight-route-validation--safeguards)
- [Request Context (`rest.Context`)](#request-context-restcontext)
  - [URL Parameters & Request Metadata](#url-parameters--request-metadata)
  - [Query Parameters](#query-parameters)
  - [Form Parameters](#form-parameters)
  - [Headers & Cookies](#headers--cookies)
  - [Request Binding](#request-binding)
  - [Typed Context Values](#typed-context-values)
  - [Complete Handler Example](#complete-handler-example)
- [Standardized Response Envelopes](#standardized-response-envelopes)
  - [Response Interface](#response-interface)
  - [Success Responses](#success-responses)
  - [RFC 9457 Problem Details](#rfc-9457-problem-details)
- [Structured Logging (`pkg/logging` & `pkg/logger`)](#structured-logging-pkglogging--pkglogger)
  - [Core Logger (`pkg/logging`)](#core-logger-pkglogging)
  - [Log Severity Levels](#log-severity-levels)
  - [Output Formats](#output-formats)
  - [Global Logger & Package Helpers](#global-logger--package-helpers)
  - [OpenTelemetry Trace Correlation (`TraceHandler`)](#opentelemetry-trace-correlation-tracehandler)
  - [HTTP Access Logging Middleware (`pkg/logger`)](#http-access-logging-middleware-pkglogger)
- [Distributed Tracing & Observability](#distributed-tracing--observability)
  - [OpenTelemetry Architecture](#opentelemetry-architecture)
  - [Tracer Provider Initialization](#tracer-provider-initialization)
  - [Exporter Protocols & Transports](#exporter-protocols--transports)
  - [Sampling Strategies](#sampling-strategies)
  - [Resource Detection & Attributes](#resource-detection--attributes)
  - [W3C Distributed Context Propagation](#w3c-distributed-context-propagation)
  - [HTTP Tracing Middleware](#http-tracing-middleware)
- [Reserved Framework Endpoints](#reserved-framework-endpoints)
  - [Liveness Probe (`GET /health`)](#liveness-probe-get-health)
  - [Readiness Probe (`GET /ready`)](#readiness-probe-get-ready)
  - [Runtime Profiling & Metrics](#runtime-profiling--metrics)
  - [Collision Protection Guarantees](#collision-protection-guarantees)
- [MongoDB Client (`pkg/mongodb`)](#mongodb-client-pkgmongodb)
  - [Connecting & Lifecycle](#connecting--lifecycle)
  - [Two-Phase Automatic TLS Fallback](#two-phase-automatic-tls-fallback)
  - [Startup Ping Verification](#startup-ping-verification)
  - [Database & Collection Handles](#database--collection-handles)
  - [Readiness & Health Verification](#readiness--health-verification)
  - [Raw Driver Access](#raw-driver-access)
  - [Termination Hooks & Sentinel Errors](#termination-hooks--sentinel-errors)
- [MongoDB Testcontainers (`pkg/testcontainers/mongodb`)](#mongodb-testcontainers-pkgtestcontainersmongodb)
  - [Container Lifecycle & Startup](#container-lifecycle--startup)
  - [Functional Options & Customization](#functional-options--customization)
  - [Connection Getters & Managed Client](#connection-getters--managed-client)
  - [Integration Testing Patterns](#integration-testing-patterns)
  - [Constants & Sentinel Errors](#constants--sentinel-errors)
- [Configuration](#configuration)
  - [Server General Settings](#server-general-settings)
  - [Server Timeouts](#server-timeouts)
  - [Server Resource Limits](#server-resource-limits)
  - [HTTP Path & Routing Handling](#http-path--routing-handling)
  - [Upstream Proxy & Forwarded Headers](#upstream-proxy--forwarded-headers)
  - [CORS Policy Configuration](#cors-policy-configuration)
  - [OpenTelemetry Tracing Configuration](#opentelemetry-tracing-configuration)
  - [MongoDB Configuration](#mongodb-configuration)
  - [Structured Logger Configuration](#structured-logger-configuration)
- [Consumer Bootstrap](#consumer-bootstrap)
- [Development Workflows](#development-workflows)
  - [Makefile Targets](#makefile-targets)
  - [Testing & Code Coverage](#testing--code-coverage)
  - [Static Analysis & Linting](#static-analysis--linting)
  - [Code Formatting](#code-formatting)
  - [Dependency Management & Verification](#dependency-management--verification)
  - [Full Verification Pipeline](#full-verification-pipeline)
- [Continuous Integration](#continuous-integration)

---

## Prerequisites & Installation

### Prerequisites

- **Go**: Version `1.27.1` or higher (tested with Go `1.27+`)

### Installation

Install the library in your Go module:

```bash
go get github.com/nawaphonOHM/whatever
```

---

## Architecture & Directory Layout

This repository is structured as a modular library. Consuming microservices import public packages under `pkg/` while internal engine wiring, telemetry providers, and route validation logic remain encapsulated under `internal/`.

```
.
├── internal/
│   ├── mongodb/
│   │   ├── client/            # Managed MongoDB v2 client wrapper, pooling, and TLS fallback
│   │   └── config/            # MongoDB environment configuration loader and validation
│   ├── opentelemetry/
│   │   ├── config/            # OTel exporter and sampler configuration
│   │   ├── middleware/        # OTel HTTP tracing middleware and W3C trace propagation
│   │   └── provider/          # Tracer provider initialization, samplers, and OTLP exporters
│   ├── rest/
│   │   ├── config/            # REST server environment configuration loader and validation
│   │   ├── contracts/         # Core API, Context, Blueprint, and Response interfaces
│   │   ├── health/            # Built-in liveness (/health) and readiness (/ready) probe handlers
│   │   ├── middleware/        # CORS, Panic Recovery, Request ID, and Access Log middlewares
│   │   ├── problem/           # RFC 9457 Problem Details error response implementation
│   │   └── server/            # Gin engine bootstrap, preflight route validation, and server lifecycle
│   └── testcontainers/
│       └── mongodb/           # Testcontainers MongoDB module integration, lifecycle, and client wiring
├── pkg/
│   ├── logger/                # Gin HTTP access logging middleware with OTel trace correlation
│   ├── logging/               # Structured slog-based logging utilities with TRACE/FATAL levels
│   ├── mongodb/               # Public MongoDB connection entrypoint and managed client
│   ├── rest/                  # Declarative Blueprint routing contracts, Context, Response, and StartREST
│   └── testcontainers/
│       └── mongodb/           # Public Testcontainers MongoDB testing runner, options, and container handle
├── .github/
│   └── workflows/
│       └── ci.yml             # Continuous integration pipeline
├── Makefile                   # Local development, test, and build targets
├── go.mod                     # Go module definition (Go 1.27.1+)
├── go.sum                     # Go module checksums
└── README.md                  # Project documentation
```

---

## Package Catalog & Feature Matrix

The toolkit is divided into focused public packages under `pkg/`:

| Package | Primary Role | Key Types & Functions | Underlying Technology |
|---|---|---|---|
| [`pkg/rest`](#pkgrest) | Declarative REST API framework, context, response envelopes, & server lifecycle | `rest.NewBluePrint()`, `rest.ExportableAPI`, `rest.StartREST()`, `rest.OK()`, `rest.BadRequest()` | `gin-gonic/gin`, RFC 9457 |
| [`pkg/logging`](#pkglogging) | Structured slog logging with extended levels & OTel trace injection | `logging.New()`, `logging.TraceContext()`, `logging.FatalContext()`, `logging.NewTraceHandler()` | `log/slog`, `go.opentelemetry.io/otel` |
| [`pkg/logger`](#pkglogger) | Gin HTTP access logging middleware with trace context correlation | `logger.Logger()`, `logger.WithLogger()`, `logger.WithConfig()`, `logger.GetRequestID()` | `gin-gonic/gin`, `log/slog`, OpenTelemetry |
| [`pkg/mongodb`](#pkgmongodb) | Managed MongoDB client with auto-TLS fallback, pooling & health verification | `mongodb.Connect()`, `client.Database()`, `client.Collection()`, `client.Ping()`, `client.RawClient()` | `go.mongodb.org/mongo-driver/v2` |
| [`pkg/testcontainers/mongodb`](#pkgtestcontainersmongodb) | Ephemeral MongoDB containers with automated lifecycle and connection wiring for integration tests | `mongodb.Run()`, `container.Client()`, `container.ConnectionString()`, `container.Terminate()` | `testcontainers-go`, Docker |

### `pkg/rest`

The primary REST framework package provides declarative route registration contracts, typed request context access, standardized JSON response envelopes, and complete server lifecycle management:

- **Declarative Blueprint Routing**: Define route groups using `rest.RRestAPIRegistration` and `rest.ExportableAPI` with semantic versioning (`/api/v1`, `/api/v2`) and automated canonical URL path calculation.
- **Preflight Route Validation**: Detects conflicting paths, duplicate route definitions, missing handlers, and reserved health endpoint collisions at bootstrap before binding sockets.
- **RFC 9457 Problem Details**: Compliant error representations (`rest.BadRequest`, `rest.Unauthorized`, `rest.Forbidden`, `rest.NotFound`, `rest.InternalServerError`, `rest.Error`).
- **Standardized Response Builders**: Consistent JSON success formatting (`rest.OK`, `rest.Created`, `rest.NoContent`, `rest.JSON`).
- **Context Abstraction**: Unified `rest.Context` interface for URL params, query parameters, headers, cookies, JSON payload binding (`c.ShouldBindJSON`), and typed context values.
- **CORS & Metadata Configuration**: Flexible CORS configuration via `rest.NewCorsSetting()` and `rest.NewMeta()`.
- **Built-in Probes & Lifecycle**: Mounts reserved `/health` (liveness) and `/ready` (readiness) probe endpoints and provides signal-driven graceful shutdown via `rest.StartREST(bp)`.

### `pkg/logging`

Structured logging built on Go standard library `log/slog` with extended severity levels, flexible formatters, and native OpenTelemetry correlation:

- **Extended Severity Levels**: Native support for `TRACE` (level -8) and `FATAL` (level 12) in addition to standard `DEBUG`, `INFO`, `WARN`, `ERROR`.
- **Custom Output Formatters**: Easy configuration for JSON (`FormatJSON`, `NewJSON`) and Text (`FormatText`, `NewText`) outputs.
- **Global & Instance Support**: Use package-level convenience functions (`logging.InfoContext`, `logging.ErrorContext`, `logging.FatalContext`) or isolated `*logging.Logger` instances.
- **OpenTelemetry Correlation**: Automatic `trace_id` and `span_id` attribute injection from active context spans via `TraceHandler`.
- **Injectable Process Control**: Configurable exit hook for fatal logging (`SetExitFunc`), ideal for testing.

### `pkg/logger`

Zero-allocation Gin HTTP access logging middleware designed for production observability:

- **Comprehensive Request Metrics**: Automatically logs request duration/latency, client IP, HTTP method, URL path, query string, response status code, and response payload size in bytes.
- **Trace Context Extraction**: Automatically extracts active span and trace IDs from the Gin request context or W3C distributed tracing headers.
- **Status-Based Log Leveling**: Assigns `INFO` for 2xx/3xx responses, `WARN` for 4xx client errors, and `ERROR` for 5xx server errors.
- **Path Filtering**: Exclude high-frequency health probes or metrics scrapers via configurable `SkipPaths`.

### `pkg/mongodb`

Managed client for MongoDB deployments using the official MongoDB Go driver v2 (`go.mongodb.org/mongo-driver/v2`):

- **Zero-Boilerplate Initialization**: Seamlessly reads and validates configuration from `OHM9996_MONGODB_*` environment variables.
- **Two-Phase Automatic TLS Fallback**: Attempts unencrypted connection first and automatically negotiates TLS if required by the remote cluster (e.g. MongoDB Atlas).
- **Managed Connection Pool**: Configurable connection limits (`MaxPoolSize`, `MinPoolSize`), socket timeouts, and connect timeouts.
- **Health Verification**: Built-in `Ping(ctx)` method to verify live cluster connectivity during readiness checks.
- **Direct Handle & Raw Driver Access**: Provides `client.Database(...)` and `client.Collection(...)` helpers with fallback to default database, and `client.RawClient()` for transactions and change streams.

### `pkg/testcontainers/mongodb`

Dedicated Testcontainers integration for spin-up and teardown of ephemeral MongoDB instances in automated test suites:

- **Ephemeral Container Lifecycle**: Spin up isolated, throwaway MongoDB containers in Go tests using `mongodb.Run(ctx, opts...)` and terminate them cleanly with `container.Terminate(ctx)`.
- **Preconfigured Managed Client**: Directly obtain a ready-to-use, ping-verified `*pkg/mongodb.Client` prewired with mapped host and port via `container.Client(ctx)`.
- **Dynamic Connection Introspection**: Resolve mapped external ports, host addresses, and full connection strings via `container.Port(ctx)`, `container.Host(ctx)`, and `container.ConnectionString(ctx)`.
- **Declarative Container Customization**: Configure custom Docker images (`WithImage`), replica set names (`WithReplicaSet`), default database initialization (`WithDatabase`), credentials (`WithUsername`, `WithPassword`), environment variables (`WithEnv`), and raw container customizers (`WithContainerOptions`).

---

## Quick Start

Get a high-performance HTTP REST microservice running with declarative routing, structured logging, and signal-driven graceful shutdown in just a few lines of code:

```go
package main

import (
	"os"

	"github.com/nawaphonOHM/whatever/pkg/logging"
	"github.com/nawaphonOHM/whatever/pkg/rest"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// 1. Initialize structured JSON logging
	logger := logging.New(logging.Config{
		Output: os.Stdout,
		Level:  logging.LevelInfo,
		Format: logging.FormatJSON,
	})
	logging.SetDefault(logger)

	// 2. Declare API route endpoints and versioned group
	itemRoutes := &rest.RRestAPIRegistration{
		Version: 1,       // Generates /api/v1 prefix (/api/v1/items)
		Prefix:  "/items",
		Apis: []*rest.ExportableAPI{
			{
				Path:   "",
				Method: rest.GET,
				Handler: func(c rest.Context) rest.Response {
					items := []Item{
						{ID: "1", Name: "Widget"},
						{ID: "2", Name: "Gadget"},
					}
					logging.Info("Fetched item catalog", "count", len(items))
					return rest.OK(items, "Items retrieved successfully")
				},
			},
			{
				Path:   "/:id",
				Method: rest.GET,
				Handler: func(c rest.Context) rest.Response {
					id := c.Param("id")
					if id == "" {
						return rest.BadRequest("INVALID_ID", "Item ID is required")
					}
					return rest.OK(Item{ID: id, Name: "Widget"})
				},
			},
		},
	}

	// 3. Configure CORS policy and assemble application blueprint
	cors := rest.NewCorsSetting().
		WithAllowOrigin("http://localhost:3000").
		WithAllowHTTPMethods(rest.GET, rest.POST, rest.OPTIONS)

	bp := rest.NewBluePrint().
		WithMeta(rest.NewMeta().WithCors(cors)).
		WithAPIs(itemRoutes)

	// 4. Start the REST server
	// StartREST initializes Gin, binds OHM9996_ configuration, validates routes,
	// mounts health probes (/health, /ready), and manages graceful shutdown on SIGINT/SIGTERM.
	logging.Info("Starting REST service...")
	if err := rest.StartREST(bp); err != nil {
		logging.Fatal("Server terminated with error", "error", err)
	}
	logging.Info("Server exited cleanly")
}
```

---

## REST Routing & Blueprint API

The library utilizes a declarative routing paradigm built on top of `rest.BluePrint`, `rest.RRestAPIRegistration`, and `rest.ExportableAPI`. Rather than registering routes imperatively on a mutable router instance, domain packages declare their routing contracts as pure data structures. The engine then validates, computes canonical URLs, and mounts the entire tree into the underlying Gin engine during `rest.StartREST(bp)`.

### Blueprint & Server Middleware Pipeline

Application bootstrap is centered around `rest.BluePrint`, constructed via `rest.NewBluePrint()`. A blueprint encapsulates server metadata (such as CORS policies) and registered API route groups:

| Blueprint Method | Description |
|---|---|
| `rest.NewBluePrint()` | Initializes an empty `*rest.BluePrint` builder |
| `bp.WithMeta(meta *rest.Meta)` | Assigns server metadata (e.g. CORS settings) to the blueprint |
| `bp.WithAPIs(apis ...*rest.RRestAPIRegistration)` | Replaces the registered API route groups with the provided slice |
| `bp.AddAPIs(apis ...*rest.RRestAPIRegistration)` | Appends additional API route groups to the blueprint |
| `bp.Meta()` | Returns the configured `*rest.Meta` or `nil` |
| `bp.Apis()` | Returns a copy slice of all registered API route groups |

When `rest.StartREST(bp)` is invoked, it configures the underlying HTTP engine and automatically provisions an enterprise-grade middleware and lifecycle pipeline:

1. **OpenTelemetry Distributed Tracing**: Automatic tracing middleware (`otelmw.Middleware`) extracting W3C `traceparent` headers and instrumenting HTTP spans.
2. **Request ID Tracking**: Generates or propagates `X-Request-ID` headers across all requests via `middleware.RequestID()`.
3. **Structured HTTP Access Logging**: High-throughput access logger (`pkg/logger`) recording method, path, status, latency, client IP, and trace IDs (configurable via `OHM9996_SERVER_ENABLE_ACCESS_LOG`).
4. **Panic Recovery**: Robust recovery middleware (`middleware.Recovery()`) capturing unhandled panics and writing RFC 9457 Problem Details error responses.
5. **CORS Handling**: Attaches CORS middleware (`middleware.CORS()`) using rules configured via `bp.WithMeta(rest.NewMeta().WithCors(...))` or permissive environment defaults.
6. **Liveness & Readiness Probes**: Mounts reserved health check endpoints (`GET /health` and `GET /ready`) exempted from both `/api` and version prefixes.
7. **Runtime Profiling (pprof)**: Exposes profiling endpoints at `OHM9996_SERVER_PROFILE_PATH` when enabled (`OHM9996_SERVER_ENABLE_PROFILING=true`).
8. **Signal-Driven Graceful Shutdown**: Automatically captures `SIGINT` and `SIGTERM` signals and allows inflight requests to complete within `OHM9996_SERVER_SHUTDOWN_TIMEOUT`.

### Declarative Registration Model

An API group is defined by creating an instance of `*rest.RRestAPIRegistration`, which encapsulates an API version, a route prefix, and a slice of `*rest.ExportableAPI` endpoints:

```go
package users

import (
	"github.com/nawaphonOHM/whatever/pkg/rest"
)

// NewUserAPIs returns declarative API registrations for user endpoints.
func NewUserAPIs() *rest.RRestAPIRegistration {
	return &rest.RRestAPIRegistration{
		Version: 1,       // Generates /api/v1 prefix
		Prefix:  "/users", // Route group prefix: /api/v1/users
		Apis: []*rest.ExportableAPI{
			{
				Path:       "",
				Method:     rest.GET,
				Middleware: []rest.Middleware{authMiddleware},
				Handler:    listUsersHandler,
			},
			{
				Path:       "/:id",
				Method:     rest.GET,
				Middleware: []rest.Middleware{authMiddleware},
				Handler:    getUserHandler,
			},
			{
				Path:       "",
				Method:     rest.POST,
				Middleware: []rest.Middleware{authMiddleware, validateBodyMiddleware},
				Handler:    createUserHandler,
			},
			{
				Path:       "/:id",
				Method:     rest.DELETE,
				Middleware: []rest.Middleware{authMiddleware, adminOnlyMiddleware},
				Handler:    deleteUserHandler,
			},
		},
	}
}
```

To mount API groups into the server, initialize a blueprint with `rest.NewBluePrint()` and attach them via `WithAPIs` (replaces list) or `AddAPIs` (appends to list):

```go
bp := rest.NewBluePrint().WithAPIs(
	users.NewUserAPIs(),
	orders.NewOrderAPIs(),
	products.NewProductAPIs(),
)
```

### HTTP Methods

The `pkg/rest` package exports typed integer constants for all standard HTTP methods:

| Method Constant | HTTP Verb | Description |
|---|---|---|
| `rest.GET` | `GET` | Retrieve resource representation |
| `rest.HEAD` | `HEAD` | Retrieve resource headers without body |
| `rest.POST` | `POST` | Create resource or execute operation |
| `rest.PUT` | `PUT` | Replace resource representation |
| `rest.PATCH` | `PATCH` | Apply partial modification to resource |
| `rest.DELETE` | `DELETE` | Remove resource |
| `rest.CONNECT` | `CONNECT` | Establish tunnel to server |
| `rest.OPTIONS` | `OPTIONS` | Describe communication options for resource |
| `rest.TRACE` | `TRACE` | Perform message loop-back test |

### API Versioning & Path Resolution

Full route paths are computed using a deterministic canonical resolution algorithm:

$$\text{FullPath} = \text{CalculateFullPath}(\text{Version}, \text{Prefix}, \text{Path})$$

#### Resolution Rules:
1. **API and Version Prefixes**: Every non-reserved route is resolved as `/api[/v<N>]<prefix><path>`. If `Version > 0`, the version segment `/v<Version>` follows the mandatory `/api` prefix (e.g. `Version: 1` becomes `/api/v1`). If the prefix already starts with `/api` or `/v<Version>`, that leading token is not duplicated.
2. **Path Normalization**: Leading slashes are added automatically if omitted, trailing slashes are trimmed (except for the root `/`), and redundant duplicate slashes (`//`) are normalized to a single slash (`/`).
3. **Reserved Route Exemption**: Internal framework health probes (`/health`, `/ready`) are never API- or version-prefixed and remain available at those exact paths.

#### Path Resolution Examples:

| Version | Prefix | Endpoint Path | Resolved Full Path |
|---|---|---|---|
| `1` | `"/items"` | `""` | `/api/v1/items` |
| `1` | `"/items"` | `"/:id"` | `/api/v1/items/:id` |
| `2` | `"orders"` | `"summary"` | `/api/v2/orders/summary` |
| `1` | `"/v1/billing"` | `"/invoices"` | `/api/v1/billing/invoices` |
| `0` | `"/webhooks"` | `"/stripe"` | `/api/webhooks/stripe` |
| `0` | `""` | `"/ping"` | `/api/ping` |

### CORS & Blueprint Metadata

Custom Cross-Origin Resource Sharing (CORS) policies can be configured on the blueprint via `rest.NewMeta()` and `rest.NewCorsSetting()`:

```go
cors := rest.NewCorsSetting().
	WithAllowOrigin("https://app.example.com", "https://admin.example.com").
	WithAllowHTTPMethods(rest.GET, rest.POST, rest.PUT, rest.DELETE, rest.OPTIONS)

meta := rest.NewMeta().WithCors(cors)

bp := rest.NewBluePrint().
	WithMeta(meta).
	WithAPIs(apiRegistrations...)
```

If no CORS configuration is provided, the server defaults to a permissive CORS configuration allowing all origins (`*`) and standard headers.

### Middleware Pipeline

Middlewares in `pkg/rest` have the signature:

```go
type Middleware = func(Context)
```

Middlewares are registered per endpoint in the `Middleware` slice of `rest.ExportableAPI`. When a request is received, the middleware pipeline executes in sequential order before the endpoint `Handler`:

```go
func authMiddleware(c rest.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		// Middleware can inspect headers or context state
		return
	}
}

api := &rest.ExportableAPI{
	Path:       "/profile",
	Method:     rest.GET,
	Middleware: []rest.Middleware{authMiddleware, auditLogMiddleware},
	Handler: func(c rest.Context) rest.Response {
		return rest.OK(map[string]string{"status": "authenticated"})
	},
}
```

### Preflight Route Validation & Safeguards

When `rest.StartREST(bp)` is invoked, it validates the entire route graph before starting the HTTP listener. If any rule is violated, the server fails fast and returns a sentinel error:

| Sentinel Error | Cause | Remediation |
|---|---|---|
| `rest.ErrNilBluePrint` | A `nil` blueprint was passed to `rest.StartREST` | Pass a valid `*rest.BluePrint` initialized via `rest.NewBluePrint()` |
| `rest.ErrNilRegistration` | An element of the registrations slice is `nil` | Ensure all elements in `bp.Apis()` are non-nil |
| `rest.ErrNilAPI` | An element of `reg.Apis` is `nil` | Ensure all `*rest.ExportableAPI` pointers in `reg.Apis` are non-nil |
| `rest.ErrNilHandler` | An API entry has a `nil` handler function | Supply a valid `func(rest.Context) rest.Response` handler |
| `rest.ErrInvalidMethod` | An unknown HTTP method constant was provided | Use a valid constant (`rest.GET`, `rest.POST`, etc.) |
| `rest.ErrReservedPath` | A route conflicts with `/health` or `/ready` | Change the route path; `/health` and `/ready` are reserved for framework health probes |
| `rest.ErrDuplicateRoute` | Two or more APIs share the same HTTP method and resolved path | Resolve overlapping route paths or methods across registration groups |

---

## Request Context (`rest.Context`)

Every route handler and middleware receives a `rest.Context` interface, which wraps the underlying HTTP request and provides safe, high-level accessors for URL parameters, query strings, headers, form values, payload binding, and context state.

### URL Parameters & Request Metadata

| Method | Return Type | Description |
|---|---|---|
| `c.Param(key string)` | `string` | Retrieves a named URL path parameter (e.g. `/:id` -> `c.Param("id")`) |
| `c.FullPath()` | `string` | Returns the matched route pattern template (e.g. `"/api/v1/users/:id"`) |
| `c.ClientIP()` | `string` | Resolves the client IP address considering proxy headers |
| `c.ContentType()` | `string` | Returns the `Content-Type` header of the incoming request |

### Query Parameters

| Method | Return Type | Description |
|---|---|---|
| `c.Query(key string)` | `string` | Returns the query parameter value, or `""` if absent |
| `c.DefaultQuery(key, defaultValue string)` | `string` | Returns the query value, or `defaultValue` if absent |
| `c.QueryArray(key string)` | `[]string` | Returns all values for a repeated query parameter (`?tag=a&tag=b`) |
| `c.QueryMap(key string)` | `map[string]string` | Returns map query parameters (`?filter[name]=john&filter[status]=active`) |

### Form Parameters

| Method | Return Type | Description |
|---|---|---|
| `c.PostForm(key string)` | `string` | Returns a form parameter from POST, PUT, or PATCH request bodies |
| `c.DefaultPostForm(key, defaultValue string)` | `string` | Returns the form value, or `defaultValue` if empty |
| `c.PostFormArray(key string)` | `[]string` | Returns a slice of form values for a repeated form key |
| `c.PostFormMap(key string)` | `map[string]string` | Returns form values as a map for nested form fields |

### Headers & Cookies

| Method | Return Type | Description |
|---|---|---|
| `c.GetHeader(key string)` | `string` | Retrieves the value of a request header (e.g. `c.GetHeader("Authorization")`) |
| `c.Cookie(name string)` | `(string, error)` | Retrieves the value of a named cookie, or returns `http.ErrNoCookie` |

### Request Binding

The `rest.Context` interface provides declarative model binding using standard struct tags (`json`, `form`, `uri`, `header`, `binding`):

| Method | Description |
|---|---|
| `c.ShouldBindJSON(obj any) error` | Unmarshals a JSON request body into a pointer struct |
| `c.ShouldBindQuery(obj any) error` | Binds URL query string parameters into a pointer struct (`form` tag) |
| `c.ShouldBindURI(obj any) error` | Binds URL path parameters into a pointer struct (`uri` tag) |
| `c.ShouldBindHeader(obj any) error` | Binds HTTP request headers into a pointer struct (`header` tag) |
| `c.ShouldBind(obj any) error` | Automatically chooses binding engine based on HTTP Method and Content-Type |
| `c.BindJSON(obj any) error` | Binds JSON and writes HTTP 400 Bad Request to Gin context on validation failure |
| `c.BindQuery(obj any) error` | Binds query parameters and sets HTTP 400 on error |
| `c.BindURI(obj any) error` | Binds URI parameters and sets HTTP 400 on error |
| `c.BindHeader(obj any) error` | Binds request headers and sets HTTP 400 on error |
| `c.Bind(obj any) error` | Binds request body and sets HTTP 400 on error |

> **Best Practice**: Prefer `ShouldBind*` methods in your handlers so you can return customized RFC 9457 Problem Details envelopes (via `rest.BadRequest`) rather than generic 400 responses.

### Typed Context Values

For values injected into the context by upstream middlewares, `rest.Context` provides type-safe retrieval helpers:

| Method | Return Type | Description |
|---|---|---|
| `c.Get(key string)` | `(any, bool)` | Retrieves a raw value by key and indicates if it exists |
| `c.MustGet(key string)` | `any` | Retrieves a raw value or panics if the key is missing |
| `c.GetString(key string)` | `string` | Retrieves value as `string` (or `""` if absent) |
| `c.GetBool(key string)` | `bool` | Retrieves value as `bool` (or `false` if absent) |
| `c.GetInt(key string)` | `int` | Retrieves value as `int` (or `0` if absent) |
| `c.GetInt64(key string)` | `int64` | Retrieves value as `int64` (or `0` if absent) |
| `c.GetUint(key string)` | `uint` | Retrieves value as `uint` (or `0` if absent) |
| `c.GetUint64(key string)` | `uint64` | Retrieves value as `uint64` (or `0` if absent) |
| `c.GetFloat64(key string)` | `float64` | Retrieves value as `float64` (or `0.0` if absent) |
| `c.GetTime(key string)` | `time.Time` | Retrieves value as `time.Time` |
| `c.GetDuration(key string)` | `time.Duration` | Retrieves value as `time.Duration` |
| `c.GetStringSlice(key string)` | `[]string` | Retrieves value as `[]string` |
| `c.GetStringMap(key string)` | `map[string]any` | Retrieves value as `map[string]any` |
| `c.GetStringMapString(key string)` | `map[string]string` | Retrieves value as `map[string]string` |
| `c.GetStringMapStringSlice(key string)` | `map[string][]string` | Retrieves value as `map[string][]string` |

### Complete Handler Example

```go
package users

import (
	"github.com/nawaphonOHM/whatever/pkg/rest"
)

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role"`
}

type UserURIParams struct {
	ID string `uri:"id" binding:"required"`
}

func CreateUserHandler(c rest.Context) rest.Response {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return rest.BadRequest("INVALID_PAYLOAD", "Request body validation failed", err.Error())
	}

	// Read optional query parameter with default fallback
	sendWelcome := c.DefaultQuery("notify", "true") == "true"

	// Read header
	tenantID := c.GetHeader("X-Tenant-ID")

	user := map[string]any{
		"id":        "usr_12345",
		"name":      req.Name,
		"email":     req.Email,
		"tenant_id": tenantID,
		"notified":  sendWelcome,
	}

	return rest.Created(user, "User created successfully")
}

func GetUserHandler(c rest.Context) rest.Response {
	var uri UserURIParams
	if err := c.ShouldBindURI(&uri); err != nil {
		return rest.BadRequest("INVALID_URI", "Invalid user ID in path")
	}

	if uri.ID != "123" {
		return rest.NotFound("USER_NOT_FOUND", "No user found with the provided ID", map[string]string{
			"user_id": uri.ID,
		})
	}

	return rest.OK(map[string]string{"id": uri.ID, "name": "Alice"})
}
```

---

## Standardized Response Envelopes

Every route handler returns a `rest.Response`. The library automatically serializes responses with appropriate HTTP headers, status codes, and uniform JSON structures conforming to industry standards.

### Response Interface

```go
type Response interface {
	StatusCode() int
	Write(c *gin.Context)
}
```

### Success Responses

The library provides constructor helpers for standard HTTP success statuses:

| Helper Function | HTTP Status | Description |
|---|---|---|
| `rest.OK(data any, message ...string)` | `200 OK` | Standard success response with data payload |
| `rest.Created(data any, message ...string)` | `201 Created` | Resource creation response |
| `rest.NoContent()` | `204 No Content` | Empty success response without body |
| `rest.JSON(statusCode int, data any, message ...string)` | `statusCode` | Generic JSON success response with custom status code |

#### Success Response JSON Schema

Success responses produce a JSON envelope with `data`, `timestamp` (RFC 3339 format), `message`, and `success` fields:

```json
{
  "success": true,
  "data": {
    "id": "1",
    "name": "Widget",
    "price": 29.99
  },
  "message": "Item retrieved successfully",
  "timestamp": "2026-09-28T12:00:00Z"
}
```

For client unmarshaling and type-safe testing, the generic `rest.SuccessResponse[T]` (and its alias `rest.Envelope[T]`) is exported:

```go
type SuccessResponse[T any] struct {
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message,omitempty"`
	Success   bool      `json:"success"`
}
```

### RFC 9457 Problem Details

For error handling, the framework implements [RFC 9457 (Problem Details for HTTP APIs)](https://datatracker.ietf.org/doc/html/rfc9457) with `Content-Type: application/problem+json`.

| Helper Function | HTTP Status | Default Title |
|---|---|---|
| `rest.BadRequest(code, detail string, details ...any)` | `400 Bad Request` | `"Bad Request"` |
| `rest.Unauthorized(code, detail string, details ...any)` | `401 Unauthorized` | `"Unauthorized"` |
| `rest.Forbidden(code, detail string, details ...any)` | `403 Forbidden` | `"Forbidden"` |
| `rest.NotFound(code, detail string, details ...any)` | `404 Not Found` | `"Not Found"` |
| `rest.InternalServerError(code, detail string, details ...any)` | `500 Internal Server Error` | `"Internal Server Error"` |
| `rest.Error(statusCode int, code, detail string, details ...any)` | `statusCode` | Standard HTTP text for `statusCode` |

#### Problem Details Schema & Fields

| Field | Type | Description |
|---|---|---|
| `type` | `string` | Problem type URI (defaults to `"about:blank"`) |
| `title` | `string` | Short, human-readable summary of the HTTP problem |
| `status` | `int` | HTTP status code integer |
| `code` | `string` | Machine-readable application error code (e.g. `"INVALID_PAYLOAD"`, `"USER_NOT_FOUND"`) |
| `detail` | `string` | Human-readable explanation specific to this occurrence of the problem |
| `instance` | `string` | Request URI path that originated the error (automatically extracted from request) |
| `details` | `any` | Optional structured error metadata (e.g. field validation errors, debug context) |

#### Problem Details Example

When returning an error from a handler:

```go
return rest.BadRequest(
	"VALIDATION_FAILED",
	"The submitted registration payload contains invalid fields",
	map[string]string{
		"email": "must be a valid corporate email address",
		"password": "must be at least 12 characters long",
	},
)
```

The response is rendered as:

```http
HTTP/1.1 400 Bad Request
Content-Type: application/problem+json

{
  "type": "about:blank",
  "title": "Bad Request",
  "status": 400,
  "code": "VALIDATION_FAILED",
  "detail": "The submitted registration payload contains invalid fields",
  "instance": "/api/v1/users/register",
  "details": {
    "email": "must be a valid corporate email address",
    "password": "must be at least 12 characters long"
  }
}
```

---

## Structured Logging (`pkg/logging` & `pkg/logger`)

The library provides a two-layer logging architecture: `pkg/logging` provides core general-purpose structured logging on top of standard library `log/slog` with extended severity levels (`TRACE`, `FATAL`) and automatic OpenTelemetry trace correlation, while `pkg/logger` provides pre-configured Gin HTTP access logging middleware with latency tracking, client IP resolution, and status-based log level escalation.

### Core Logger (`pkg/logging`)

The `pkg/logging` package encapsulates `*slog.Logger` and enhances it with custom handlers, pipeline formatting, and fatal termination hooks.

#### Configuration Options (`logging.Config`)

```go
type Config struct {
	// Output is the destination for log writes. Defaults to os.Stdout if nil.
	Output io.Writer

	// Level specifies the minimum severity level to log (TRACE, DEBUG, INFO, WARN, ERROR, FATAL).
	Level Level

	// Format specifies the output format (json or text). Defaults to FormatJSON.
	Format Format

	// AddSource attaches caller file and line numbers to records when true.
	AddSource bool

	// DisableTraceCorrelation disables automatic OpenTelemetry trace_id and span_id enrichment.
	DisableTraceCorrelation bool

	// ExitFunc is the function invoked on Fatal/FatalContext. Defaults to os.Exit.
	ExitFunc func(int)

	// ReplaceAttr allows customizing log attributes before writing.
	ReplaceAttr ReplaceAttrFunc

	// Handler allows supplying a pre-configured slog.Handler directly.
	Handler slog.Handler
}
```

#### Constructors & Factory Functions

| Constructor | Description |
|---|---|
| `logging.New(cfgs ...Config) *Logger` | Initializes a configured `*logging.Logger` instance (merges optional config) |
| `logging.NewJSON(w io.Writer, level Level) *Logger` | Convenience constructor for JSON-formatted logging to destination `w` |
| `logging.NewText(w io.Writer, level Level) *Logger` | Convenience constructor for human-readable text logging to destination `w` |
| `logging.NewWithHandler(h slog.Handler) *Logger` | Wraps a pre-existing `slog.Handler` directly |
| `logging.DefaultConfig() Config` | Returns recommended production defaults (`os.Stdout`, `LevelInfo`, `FormatJSON`) |

#### Instance Methods

An instance of `*logging.Logger` provides:
- **Context Logging**: `l.TraceContext(ctx, msg, args...)`, `l.DebugContext(ctx, msg, args...)`, `l.InfoContext(ctx, msg, args...)`, `l.WarnContext(ctx, msg, args...)`, `l.ErrorContext(ctx, msg, args...)`, `l.FatalContext(ctx, msg, args...)`
- **Context-Free Logging**: `l.Trace(msg, args...)`, `l.Debug(msg, args...)`, `l.Info(msg, args...)`, `l.Warn(msg, args...)`, `l.Error(msg, args...)`, `l.Fatal(msg, args...)`
- **Scoped Sub-Loggers**: `l.With(args ...any) *Logger` attaches persistent structured key-value pairs; `l.WithGroup(name string) *Logger` nests subsequent attributes within a named group.
- **Underlying Driver**: `l.Slog() *slog.Logger` returns the underlying standard library `*slog.Logger`.
- **Exit Hooks**: `l.ExitFunc()` and `l.SetExitFunc(fn func(int))` manage the exit handler invoked on fatal events.

```go
package main

import (
	"context"
	"os"

	"github.com/nawaphonOHM/whatever/pkg/logging"
)

func main() {
	// Initialize a structured JSON logger writing to stdout
	log := logging.New(logging.Config{
		Output:    os.Stdout,
		Level:     logging.LevelDebug,
		Format:    logging.FormatJSON,
		AddSource: true,
	})

	// Log with structured key-value arguments
	log.Info("Service started", "port", 8080, "environment", "production")

	// Create a sub-logger with scoped attributes
	userLogger := log.With("component", "user_service", "version", "1.0.0")
	userLogger.DebugContext(context.Background(), "Processing user request", "user_id", "usr_123")
}
```

### Log Severity Levels

The package defines strongly typed log severity levels extending standard `log/slog` levels:

| Level Constant | String Value | Underlying `slog.Level` Value | Description |
|---|---|---|---|
| `logging.LevelTrace` | `"TRACE"` | `-8` (`logging.SlogLevelTrace`) | High-volume granular diagnostic information |
| `logging.LevelDebug` | `"DEBUG"` | `-4` (`logging.SlogLevelDebug`) | Detailed debugging and troubleshooting messages |
| `logging.LevelInfo` | `"INFO"` | `0` (`logging.SlogLevelInfo`) | General informational operational events |
| `logging.LevelWarn` | `"WARN"` | `4` (`logging.SlogLevelWarn`) | Non-critical warnings and recovered conditions |
| `logging.LevelError` | `"ERROR"` | `8` (`logging.SlogLevelError`) | Actionable runtime failures and operational errors |
| `logging.LevelFatal` | `"FATAL"` | `12` (`logging.SlogLevelFatal`) | Critical unrecoverable failures; invokes configured `ExitFunc` |

#### Parsing Level Strings

The package provides case-insensitive level parsers that support standard string representations as well as common aliases:

```go
// Parse case-insensitive string into strongly typed logging.Level (supports aliases: warn/warning, fatal/critical/crit/panic)
lvl, err := logging.ParseLevel("debug")     // returns logging.LevelDebug
lvl, err := logging.ParseLevel("warning")   // returns logging.LevelWarn
lvl, err := logging.ParseLevel("critical")  // returns logging.LevelFatal

// Parse directly into standard slog.Level
slogLvl, err := logging.ParseSlogLevel("trace") // returns logging.SlogLevelTrace (-8)

// Convert standard slog.Level to canonical logging.Level
lvl = logging.LevelFromSlog(slog.LevelWarn) // returns logging.LevelWarn

// Convert logging.Level to standard slog.Level
slogLvl, err = logging.LevelFatal.SlogLevel() // returns logging.SlogLevelFatal (12)
```

### Output Formats

The package supports two primary output encodings controlled by `logging.Format`:

- **`logging.FormatJSON` (`"json"`)**: Renders machine-readable single-line JSON records formatted for log ingestion systems (Datadog, Elasticsearch, Grafana Loki, CloudWatch).
- **`logging.FormatText` (`"text"`)**: Renders human-readable key-value pairs formatted for local terminal development.

#### JSON Output Record Example:
```json
{
  "time": "2026-09-28T12:00:00.123456Z",
  "level": "INFO",
  "msg": "Order payment processed",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "order_id": "ord_9981",
  "amount": 49.99,
  "currency": "USD"
}
```

#### Text Output Record Example:
```text
time=2026-09-28T12:00:00.123456Z level=INFO msg="Order payment processed" trace_id=4bf92f3577b34da6a3ce929d0e0e4736 span_id=00f067aa0ba902b7 order_id=ord_9981 amount=49.99 currency=USD
```

### Global Logger & Package Helpers

`pkg/logging` maintains a thread-safe package-level default logger synchronized with standard library `slog.SetDefault`:

```go
// Retrieve the package-level default logger
defaultLogger := logging.Default()

// Replace the global default logger instance
customLogger := logging.NewJSON(os.Stdout, logging.LevelDebug)
logging.SetDefault(customLogger)

// Package-level logging helpers (use global default logger)
logging.Info("Application initialized", "version", "1.0.0")
logging.Warn("Rate limit threshold approached", "client_ip", "192.168.1.50")
logging.Error("Database query failed", "error", err)

// Context-aware package-level helpers (extracts active OTel trace_id & span_id)
logging.InfoContext(ctx, "User authenticated successfully", "user_id", "usr_123")
logging.ErrorContext(ctx, "Payment transaction declined", "transaction_id", "txn_554")

// Fatal logging: writes log record at FATAL level and triggers exitFunc (os.Exit(1))
logging.Fatal("Fatal startup error: unable to bind network port", "port", 8080)
```

### OpenTelemetry Trace Correlation (`TraceHandler`)

`pkg/logging` includes a native `slog.Handler` wrapper (`logging.TraceHandler`) that automatically extracts distributed tracing metadata from the active `context.Context` and appends them to every log record:

- **Active OTel Span Context**: Extracts `trace_id` and `span_id` from `trace.SpanFromContext(ctx)` when a valid OpenTelemetry span exists.
- **Context Fallback Keys**: If no active span is found, checks for typed context keys (`logging.ContextKeyTraceID`, `logging.ContextKeySpanID`) and common string fallback keys (`"trace_id"`, `"TraceID"`, `"traceId"`, `"span_id"`, `"SpanID"`, `"spanId"`).
- **Automatic Enablement**: Enabled by default in `logging.New()` unless explicitly disabled with `DisableTraceCorrelation: true`.
- **Attribute Normalization**: Automatically converts numeric `slog.Level` keys into canonical string names (`"TRACE"`, `"DEBUG"`, `"INFO"`, `"WARN"`, `"ERROR"`, `"FATAL"`) while preserving custom `ReplaceAttr` transformations.

```go
// Manually wrap any custom slog.Handler with OpenTelemetry trace correlation
baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
traceHandler := logging.NewTraceHandler(baseHandler)
logger := logging.NewWithHandler(traceHandler)

// Log with context to automatically correlate OTel trace_id & span_id
logger.InfoContext(ctx, "Processed transaction", "amount", 100.0)
```

### HTTP Access Logging Middleware (`pkg/logger`)

The `pkg/logger` package provides pre-built Gin middleware for HTTP access logging with automatic correlation to OpenTelemetry distributed traces:

```go
import (
	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/logger"
	"github.com/nawaphonOHM/whatever/pkg/logging"
)

// Attach default request logging middleware
engine.Use(logger.Logger())

// Or attach with a specific slog.Logger instance
engine.Use(logger.WithLogger(logging.Default().Slog()))

// Or attach with custom configuration (e.g. skipping noisy health checks)
engine.Use(logger.WithConfig(logger.Config{
	Logger:    logging.Default().Slog(),
	SkipPaths: []string{"/health", "/ready", "/metrics"},
}))

// Helper to extract request ID from context or X-Request-ID header
engine.GET("/ping", func(c *gin.Context) {
	reqID := logger.GetRequestID(c) // checks c.Get(logger.RequestIDKey) or c.GetHeader(logger.HeaderXRequestID)
	c.JSON(200, gin.H{"request_id": reqID})
})
```

#### Captured Access Log Attributes:

| Log Attribute | Type | Description |
|---|---|---|
| `request_id` | `string` | Unique request identifier extracted from `X-Request-ID` header or Gin context |
| `method` | `string` | HTTP request method (`GET`, `POST`, `PUT`, etc.) |
| `path` | `string` | URL request path |
| `query` | `string` | Raw URL query string |
| `status` | `int` | HTTP response status code |
| `latency_ms` | `int64` | Request execution duration in milliseconds |
| `latency` | `time.Duration` | Precise execution duration (e.g. `1.245ms`, `45.2µs`) |
| `client_ip` | `string` | Client IP address resolved through trusted proxy headers |
| `user_agent` | `string` | Client `User-Agent` header value |
| `bytes_out` | `int` | Size of written HTTP response body in bytes |
| `trace_id` | `string` | Active OpenTelemetry trace ID (32-hex character string) |
| `span_id` | `string` | Active OpenTelemetry span ID (16-hex character string) |
| `errors` | `string` | Detailed error string if Gin errors (`c.Errors`) occurred during execution |

#### Status-Based Severity Escalation:
- **`5xx` Server Errors** $\rightarrow$ Logged at **`slog.LevelError`**
- **`4xx` Client Errors** $\rightarrow$ Logged at **`slog.LevelWarn`**
- **`1xx` / `2xx` / `3xx` Success & Redirects** $\rightarrow$ Logged at **`slog.LevelInfo`**

---

## Distributed Tracing & Observability

The library provides native, zero-configuration OpenTelemetry distributed tracing and observability living within `internal/opentelemetry`. When enabled, incoming HTTP requests are automatically tracked as server spans, injected with semantic attributes, and propagated across downstream distributed services.

### OpenTelemetry Architecture

```
Incoming HTTP Request
       │
       ▼ (Extract W3C traceparent / baggage)
┌─────────────────────────────────────────────────────────┐
│  otelmw.Middleware                                      │
│  - Starts OpenTelemetry Span: "GET /api/v1/users/:id"   │
│  - Sets HTTP semantic attributes & Client IP            │
│  - Injects trace_id / span_id into context & response   │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│  Application Handler / Business Logic                   │
│  - logging.InfoContext(ctx, "Fetching user profile")    │
│    -> (Log record automatically includes trace_id)      │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼ (Inject W3C traceparent into response headers)
┌─────────────────────────────────────────────────────────┐
│  Span Finished (BatchSpanProcessor -> OTLP Exporter)    │
│  - gRPC / HTTP to OpenTelemetry Collector / Jaeger      │
└─────────────────────────────────────────────────────────┘
```

### Tracer Provider Initialization

During server bootstrap (`rest.StartREST`), the runtime initializes the global `sdktrace.TracerProvider` via `provider.InitTracerProvider(ctx, cfg)`:

- **Batch Processing**: Configures `sdktrace.NewBatchSpanProcessor` to buffer and asynchronously flush spans to the OTLP exporter.
- **Resource Composition**: Enriches all trace data with service name, application version, host identifier, and telemetry SDK metadata.
- **Graceful Flush**: Returns a `ShutdownFunc` callback registered into the server lifecycle to ensure all buffered spans are flushed during graceful shutdown.

### Exporter Protocols & Transports

The telemetry exporter sends traces to an OpenTelemetry Collector or compatible backend (Jaeger, Grafana Tempo, Datadog) via `OHM9996_OTEL_EXPORTER_OTLP_ENDPOINT` (defaults to `localhost:4317`) and `OHM9996_OTEL_EXPORTER_OTLP_PROTOCOL`:

| Protocol Key | Transport | Typical OTLP Port | Description |
|---|---|---|---|
| `grpc` | OTLP / gRPC | `4317` | High-performance binary transport over gRPC (supports `host:port` or `grpc://`) |
| `http` or `http/protobuf` | OTLP / HTTP (Protobuf) | `4318` | Standard OTLP over HTTP with Protobuf binary payload |
| `http/json` | OTLP / HTTP (JSON) | `4318` | OTLP over HTTP with JSON payload |

*Note: The environment variable `OHM9996_OTEL_EXPORTER_OTLP_ENDPOINT` defaults to `localhost:4317` regardless of the selected protocol.*

TLS security is configured via `OHM9996_OTEL_INSECURE` (defaults to `true` for local development and Kubernetes in-cluster mesh networks).

### Sampling Strategies

Trace sampling is configured via `OHM9996_OTEL_SAMPLE_RATE` (float value from `0.0` to `1.0`):

- **Parent-Based Sampling**: The tracer wraps the sampling algorithm in `sdktrace.ParentBased(...)`. If an upstream caller has already sampled a trace, downstream child spans in this microservice will always be sampled to maintain trace continuity.
- **Sampling Ratio Bounds**:
  - `1.0` (Default): 100% trace sampling (`sdktrace.AlwaysSample`).
  - `0.0`: Tracing disabled unless forced by upstream parent (`sdktrace.NeverSample`).
  - `0.0 < rate < 1.0`: Ratio-based sampling (`sdktrace.TraceIDRatioBased(rate)`).

### Resource Detection & Attributes

Every trace emitted contains standard OpenTelemetry Resource attributes identifying the service:

- `service.name`: Configured via `OHM9996_OTEL_SERVICE_NAME` (defaults to `whatever-service`).
- `service.version`: Set to application version (`1.0.0`).
- `host.name`: Automatically detected local host machine name (`resource.WithHost()`).
- `telemetry.sdk.*`: OpenTelemetry Go SDK version and language information.

### W3C Distributed Context Propagation

The provider registers a global composite text map propagator conforming to W3C specifications:
1. **W3C Trace Context (`traceparent`, `tracestate`)**: Transmits 128-bit trace IDs, 64-bit parent span IDs, and trace sampling flags across HTTP network boundaries.
2. **W3C Baggage (`baggage`)**: Propagates arbitrary contextual key-value pairs across service graphs.

The HTTP tracing middleware extracts incoming W3C headers on request entry and automatically injects them into outgoing response headers (`c.Writer.Header()`).

### HTTP Tracing Middleware

The framework engine installs `otelmw.Middleware(cfg)` on the Gin pipeline to record execution telemetry:

- **Span Naming**: Automatically names spans after HTTP method and matched route template (e.g. `GET /api/v1/users/:id`).
- **Semantic Attributes**:
  - `http.method`: HTTP method (`GET`, `POST`, `PUT`, `DELETE`, etc.)
  - `http.target`: Full request URI path with query string (e.g. `/api/v1/items?category=books`)
  - `client.address`: Client IP address
  - `http.route`: Matched parameterized route pattern
  - `http.status_code`: Response HTTP status code integer
- **Error Recording & Status**: If the handler encounters errors (`c.Errors`) or returns an HTTP status code $\ge 500$, the span status is set to `codes.Error` and the error description is attached to the span event log.
- **Skip Paths**: Framework health endpoints (`/health`, `/ready`) and any paths listed in `OHM9996_OTEL_SKIP_PATHS` bypass span creation to prevent collector flooding.

---

## Reserved Framework Endpoints

The framework pre-registers dedicated operational endpoints for container orchestration (Kubernetes, AWS ECS, Nomad), health monitoring, performance profiling, and Prometheus metrics scraping.

### Liveness Probe (`GET /health`)

The liveness probe verifies that the Go application process is alive and responsive to HTTP requests:

```http
GET /health HTTP/1.1
Host: localhost:8080
```

#### Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "status": "up",
    "timestamp": "2026-09-28T12:00:00Z",
    "version": "1.0.0"
  },
  "timestamp": "2026-09-28T12:00:00Z"
}
```

### Readiness Probe (`GET /ready`)

The readiness probe verifies that the HTTP server is ready to accept incoming production traffic:

```http
GET /ready HTTP/1.1
Host: localhost:8080
```

#### Response (`200 OK`):
```json
{
  "success": true,
  "data": {
    "status": "ready",
    "timestamp": "2026-09-28T12:00:00Z",
    "version": "1.0.0"
  },
  "timestamp": "2026-09-28T12:00:00Z"
}
```

### Runtime Profiling & Metrics

#### 1. Runtime Profiling (`GET /debug/pprof/*`)
Configured via `OHM9996_SERVER_ENABLE_PROFILING=true` at `OHM9996_SERVER_PROFILE_PATH` (defaults to `/debug/pprof`). When enabled, developers can inspect CPU profiles, heap allocations, goroutine stacks, and mutex contention:

```bash
# Analyze live heap memory profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Capture 30-second CPU profile
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Inspect active goroutine stack traces
curl http://localhost:8080/debug/pprof/goroutine?debug=1
```

#### 2. Prometheus Metrics (`GET /metrics`)
Configured via `OHM9996_SERVER_ENABLE_METRICS=true`. When enabled, standard runtime metrics (goroutines, GC pauses, memory allocations, thread counts) and HTTP metrics are exposed for scraping by Prometheus or OpenTelemetry Collector agents.

### Collision Protection Guarantees

The framework strictly protects reserved framework endpoints. During `rest.StartREST(bp)` startup, preflight validation scans all declared routes in the blueprint:
- If any application endpoint attempts to register on `/health` or `/ready`, startup aborts immediately and returns `rest.ErrReservedPath`.
- Probe paths are excluded from API versioning prefixes and bypass request access logging and distributed tracing by default.

---

## MongoDB Client (`pkg/mongodb`)

The `pkg/mongodb` package encapsulates MongoDB connection establishment, connection pooling, two-phase TLS negotiation, and lifecycle management while directly exposing official driver `*mongo.Database` and `*mongo.Collection` types for zero-overhead querying.

### Connecting & Lifecycle

Consuming applications connect to MongoDB using `mongodb.Connect(ctx)`. The library automatically reads and validates `OHM9996_MONGODB_*` environment variables, executes a two-phase connection flow, verifies connectivity via an initial ping (when enabled), and manages the underlying connection pool:

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

### Two-Phase Automatic TLS Fallback

To provide seamless connectivity across unencrypted local Docker instances, secured staging clusters, and cloud deployments (such as MongoDB Atlas), `mongodb.Connect(ctx)` implements an automated two-phase connection strategy:

1. **Unencrypted Connection Attempt**: First attempts connection without TLS (`tls=false`). If `OHM9996_MONGODB_ENABLE_PING=true` (the default), an immediate connectivity ping is performed.
2. **Automatic TLS Fallback**: If the server rejects the unencrypted connection with an error indicating TLS/SSL is required (such as `"server requires tls"`, `"server requires ssl"`, `"ssl handshake"`, `"tls handshake"`, or `"connection closed"`), the client transparently re-attempts connection with TLS enabled (`tls=true`).
3. **Graceful Failure Handling**: If connection fails after fallback or due to fatal configuration errors, the client outputs diagnostic details to `os.Stderr`, invokes the process exit hook (`exitFunc(0)`), and returns the underlying error.

### Startup Ping Verification

By default, `mongodb.Connect(ctx)` performs an active ping check during initialization to ensure that the remote MongoDB deployment is reachable before application routes begin serving traffic. If ping verification fails, the connection pool is immediately closed and a wrapped error is returned.

In serverless architectures, lazy initialization flows, or environments where MongoDB may boot concurrently after the microservice starts, startup ping verification can be toggled via environment variable:

```bash
# Disable startup ping check for lazy or deferred initialization
export OHM9996_MONGODB_ENABLE_PING=false
```

When startup ping verification is disabled, `mongodb.Connect(ctx)` initializes the client pool without blocking on network roundtrips. Applications can subsequently invoke `client.Ping(ctx)` on demand (e.g., inside readiness probes or background pollers).

### Database & Collection Handles

The managed client resolves default database names and exposes direct access to official MongoDB Go Driver v2 handles (`*mongo.Database` and `*mongo.Collection`):

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
	// Access collection using configured default database (OHM9996_MONGODB_DATABASE)
	coll := r.client.Collection("users")

	var user User
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) RecordAudit(ctx context.Context, entry bson.M) error {
	// Explicitly target a specific database override ("audit_db")
	coll := r.client.Collection("audit_logs", "audit_db")

	_, err := coll.InsertOne(ctx, entry)
	return err
}

func (r *UserRepository) CustomDatabase(ctx context.Context) error {
	// Obtain a direct *mongo.Database handle (uses default database if omitted)
	db := r.client.Database("reporting_db")
	_ = db
	return nil
}
```

- **Default Database Fallback**: Calling `client.Database()` or `client.Collection("users")` without specifying a database name automatically targets the default database configured via `OHM9996_MONGODB_DATABASE`.
- **Database Overrides**: Passing a database name parameter (`client.Database("custom_db")` or `client.Collection("users", "custom_db")`) explicitly targets the specified database.
- **Nil Safety**: Operations invoked on an uninitialized client return `nil` collections/databases without panicking.

### Readiness & Health Verification

Use `client.Ping(ctx)` to verify cluster connectivity with primary read preference (`readpref.Primary()`), ideal for Kubernetes readiness probes (`/ready`) and health checks:

```go
// Sends a primary read-preference ping command to the MongoDB cluster
if err := client.Ping(ctx); err != nil {
	log.Printf("MongoDB health check failed: %v", err)
}
```

### Raw Driver Access

When advanced MongoDB features are needed—such as multi-document transactions, client sessions, change streams, or GridFS—use `client.RawClient()` to obtain the underlying official `*mongo.Client`:

```go
rawClient := client.RawClient()
session, err := rawClient.StartSession()
if err != nil {
	return err
}
defer session.EndSession(ctx)
```

### Termination Hooks & Sentinel Errors

When required MongoDB configuration (such as `OHM9996_MONGODB_HOST`) is missing or connection fails during initialization, `mongodb.Connect` outputs diagnostic information and triggers a configurable exit hook (`os.Exit(0)` by default). In unit and integration test suites, this hook can be intercepted:

```go
// Override the process termination hook during tests (returns previous hook)
prevExit := mongodb.SetExitFunc(func(code int) {
	// Custom exit logic or test assertion
})
defer mongodb.SetExitFunc(prevExit)
```

The package exports standard sentinel errors:
- `mongodb.ErrNilClient`: Returned when attempting operations on an uninitialized client instance (`"mongodb client is not initialized"`).
- `mongodb.ErrNilConfig`: Returned when internal configuration resolving is nil (`"mongodb config cannot be nil"`).

---

## MongoDB Testcontainers (`pkg/testcontainers/mongodb`)

The `pkg/testcontainers/mongodb` package provides lightweight, isolated MongoDB container management for integration testing using [Testcontainers for Go](https://golang.testcontainers.org/). It spins up ephemeral MongoDB instances with zero host dependencies, automatically exposes connection details, and provides ready-to-use managed client instances.

### Container Lifecycle & Startup

Start an ephemeral MongoDB container using `mongodb.Run(ctx, opts...)`. The function provisions a container, assigns dynamic host port mappings, waits until MongoDB is ready to accept connections, and returns a `*mongodb.Container` handle:

```go
package repository_test

import (
	"context"
	"testing"

	tcmongo "github.com/nawaphonOHM/whatever/pkg/testcontainers/mongodb"
)

func TestContainerStartup(t *testing.T) {
	ctx := context.Background()

	// Spin up disposable MongoDB container
	container, err := tcmongo.Run(ctx)
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}
	defer func() {
		// Clean up and terminate container at test completion
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate container: %v", err)
		}
	}()

	// Introspect connection endpoints
	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}
	t.Logf("MongoDB container running at: %s", connStr)
}
```

### Functional Options & Customization

`mongodb.Run(ctx, opts...)` accepts functional options (`Option`) to customize container image, credentials, database initialization, replica sets, and environment variables:

| Option | Signature | Description | Default |
|---|---|---|---|
| `WithImage` | `WithImage(image string) Option` | Specifies the Docker image tag for the MongoDB container | `mongodb.DefaultImage` (`"mongo:6"`) |
| `WithUsername` | `WithUsername(username string) Option` | Sets root username for container authentication | `""` (no authentication) |
| `WithPassword` | `WithPassword(password string) Option` | Sets root password for container authentication | `""` (no authentication) |
| `WithDatabase` | `WithDatabase(database string) Option` | Configures default database created on startup (`MONGO_INITDB_DATABASE`) | `""` |
| `WithReplicaSet` | `WithReplicaSet(replicaSet string) Option` | Initializes a single-node replica set, required for transaction testing | `""` (standalone) |
| `WithEnv` | `WithEnv(key, value string) Option` | Sets custom environment variables inside the container | `nil` |
| `WithContainerOptions` | `WithContainerOptions(opts ...testcontainers.ContainerCustomizer) Option` | Appends raw `testcontainers-go` container customizers | `nil` |

Helper functions `mongodb.DefaultOptions()` and `mongodb.NewOptions(opts...)` are also exported for programmatic options inspection.

### Connection Getters & Managed Client

The `*mongodb.Container` instance provides helpers to retrieve dynamic connection details or immediately connect a managed `*pkg/mongodb.Client`:

| Method | Return Type | Description |
|---|---|---|
| `container.Client(ctx, opts...)` | `(*mongodb.Client, error)` | Returns a connected, ping-verified managed client prewired to container host, port, database, and auth |
| `container.ConnectionString(ctx)` | `(string, error)` | Returns the full connection URI (e.g., `mongodb://localhost:32768`) |
| `container.Host(ctx)` | `(string, error)` | Returns the host IP or hostname where the container is accessible |
| `container.Port(ctx)` | `(int, error)` | Returns the mapped external TCP port as an `int` |
| `container.Config(ctx)` | `(*config.Config, error)` | Returns a populated `Config` struct initialized with container connection details |
| `container.Database()` | `string` | Returns the configured default database name |
| `container.Username()` | `string` | Returns the configured root username |
| `container.Password()` | `string` | Returns the configured root password |
| `container.ReplicaSet()` | `string` | Returns the configured replica set name |
| `container.RawContainer()` | `*tcmongodb.MongoDBContainer` | Returns the underlying raw Testcontainers MongoDBContainer handle |
| `container.Terminate(ctx)` | `error` | Stops and deletes the running container |

### Integration Testing Patterns

#### 1. Repository Integration Test with Managed Client

Connect a managed client directly to an ephemeral container to test database operations against a real MongoDB engine:

```go
package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/nawaphonOHM/whatever/pkg/mongodb"
	tcmongo "github.com/nawaphonOHM/whatever/pkg/testcontainers/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	ID    string  `bson:"_id,omitempty"`
	SKU   string  `bson:"sku"`
	Price float64 `bson:"price"`
}

func TestProductRepository_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 1. Start MongoDB testcontainer with custom database
	container, err := tcmongo.Run(ctx,
		tcmongo.WithImage("mongo:7"),
		tcmongo.WithDatabase("test_store"),
	)
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("Failed to terminate container: %v", err)
		}
	}()

	// 2. Obtain preconfigured managed client
	client, err := container.Client(ctx)
	if err != nil {
		t.Fatalf("Failed to connect client: %v", err)
	}
	defer client.Disconnect(ctx)

	// 3. Execute database operations against default database
	coll := client.Collection("products")

	item := Product{SKU: "PROD-100", Price: 29.99}
	insertRes, err := coll.InsertOne(ctx, item)
	if err != nil {
		t.Fatalf("Failed to insert product: %v", err)
	}

	var found Product
	err = coll.FindOne(ctx, bson.M{"sku": "PROD-100"}).Decode(&found)
	if err != nil {
		t.Fatalf("Failed to query product: %v", err)
	}

	if found.Price != 29.99 {
		t.Errorf("Expected price 29.99, got %f", found.Price)
	}
	_ = insertRes
}
```

#### 2. Authenticated Container Setup

Configure username and password credentials for testing authenticated access and permission boundaries:

```go
func TestAuthenticatedMongoDB(t *testing.T) {
	ctx := context.Background()

	container, err := tcmongo.Run(ctx,
		tcmongo.WithImage("mongo:6"),
		tcmongo.WithUsername("admin"),
		tcmongo.WithPassword("supersecret"),
		tcmongo.WithDatabase("secure_db"),
	)
	if err != nil {
		t.Fatalf("Failed to start authenticated container: %v", err)
	}
	defer container.Terminate(ctx)

	// Client automatically connects using admin authSource and credentials
	client, err := container.Client(ctx)
	if err != nil {
		t.Fatalf("Failed to connect with credentials: %v", err)
	}
	defer client.Disconnect(ctx)

	// Verify health check succeeds under authenticated session
	if err := client.Ping(ctx); err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}
```

#### 3. Multi-Document Transactions with Replica Sets

MongoDB multi-document ACID transactions require a replica set topology. Use `WithReplicaSet` to spin up a single-node replica set for transaction testing:

```go
func TestMultiDocumentTransaction(t *testing.T) {
	ctx := context.Background()

	container, err := tcmongo.Run(ctx,
		tcmongo.WithReplicaSet("rs0"),
		tcmongo.WithDatabase("bank_db"),
	)
	if err != nil {
		t.Fatalf("Failed to start replica set container: %v", err)
	}
	defer container.Terminate(ctx)

	client, err := container.Client(ctx)
	if err != nil {
		t.Fatalf("Failed to connect client: %v", err)
	}
	defer client.Disconnect(ctx)

	// Access raw mongo driver client to initiate session & transaction
	rawClient := client.RawClient()
	session, err := rawClient.StartSession()
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	accounts := client.Collection("accounts")

	// Execute transactional multi-document balance transfer
	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		if _, err := accounts.InsertOne(sessCtx, bson.M{"account": "A", "balance": 100}); err != nil {
			return nil, err
		}
		if _, err := accounts.InsertOne(sessCtx, bson.M{"account": "B", "balance": 50}); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}
}
```

### Constants & Sentinel Errors

`pkg/testcontainers/mongodb` exports the following default constants and error sentinels:

- `mongodb.DefaultImage`: `"mongo:6"`
- `mongodb.DefaultPort`: `"27017/tcp"`
- `mongodb.ErrNilContainer`: Returned when calling methods on a nil container instance (`"mongodb container is nil"`).
- `mongodb.ErrContainerNotRunning`: Returned when calling methods on a container whose underlying Docker process is not running (`"mongodb container is not running"`).

---

## Configuration

All environment variables read by the library use the **`OHM9996_`** prefix.

Configuration values are resolved using the following order of precedence:
1. **System Environment Variables**: Explicit environment variables set on the host or container environment take top priority.
2. **Local `.env` File**: If a `.env` file exists in the working directory, it is automatically parsed using `godotenv`. An absent `.env` file is silently ignored.
3. **Library Defaults**: If a variable is unset across environment sources, built-in defaults defined on the configuration structs are applied.

### Server General Settings

| Variable | Type | Default | Description |
|---|---|---|---|
| `OHM9996_SERVER_HOST` | `string` | `""` (all interfaces) | Network interface address to bind HTTP listener |
| `OHM9996_SERVER_PORT` | `int` | `8080` | TCP port number for incoming HTTP requests |
| `OHM9996_GIN_MODE` | `string` | `release` | Gin engine mode (`debug`, `release`, `test`) |
| `OHM9996_APP_VERSION` | `string` | `""` | Application version reported in `/health` and `/ready` response metadata |
| `OHM9996_SERVER_DISPLAY_NAME` | `string` | `application` | Service display name used for startup banners and logging |
| `OHM9996_SERVER_ENABLE_ACCESS_LOG` | `bool` | `true` | Enables HTTP request access logging middleware |
| `OHM9996_SERVER_ENABLE_METRICS` | `bool` | `false` | Enables Prometheus `/metrics` scraping endpoint |
| `OHM9996_SERVER_ENABLE_PROFILING` | `bool` | `false` | Enables runtime pprof profiling endpoints |
| `OHM9996_SERVER_PROFILE_PATH` | `string` | `/debug/pprof` | Base URL path prefix for runtime pprof profiling endpoints |

### Server Timeouts

| Variable | Type | Default | Description |
|---|---|---|---|
| `OHM9996_SERVER_READ_TIMEOUT` | `duration` | `10s` | Maximum duration for reading the entire request (including body) |
| `OHM9996_SERVER_WRITE_TIMEOUT` | `duration` | `10s` | Maximum duration before timing out writes of the response |
| `OHM9996_SERVER_IDLE_TIMEOUT` | `duration` | `60s` | Maximum duration to wait for the next request on keep-alive connections |
| `OHM9996_SERVER_SHUTDOWN_TIMEOUT` | `duration` | `10s` | Graceful shutdown deadline before forcefully terminating active connections |

### Server Resource Limits

| Variable | Type | Default | Description |
|---|---|---|---|
| `OHM9996_SERVER_READ_HEADER_TIMEOUT` | `duration` | `5s` | Maximum duration allowed to read HTTP request headers |
| `OHM9996_SERVER_MAX_HEADER_BYTES` | `int` | `1048576` (1MB) | Maximum allowed HTTP request header size in bytes |
| `OHM9996_SERVER_MAX_BODY_SIZE` | `int64` | `33554432` (32MB) | Maximum allowed multipart memory and request body payload in bytes |

### HTTP Path & Routing Handling

| Variable | Type | Default | Description |
|---|---|---|---|
| `OHM9996_SERVER_REDIRECT_TRAILING_SLASH` | `bool` | `true` | Automatically redirects requests with or without trailing slashes |
| `OHM9996_SERVER_REDIRECT_FIXED_PATH` | `bool` | `false` | Automatically corrects URL casing and cleans redundant segments |
| `OHM9996_SERVER_HANDLE_METHOD_NOT_ALLOWED` | `bool` | `true` | Returns RFC 9457 405 Method Not Allowed when route exists for other methods |
| `OHM9996_SERVER_USE_RAW_PATH` | `bool` | `false` | Uses raw `req.URL.RawPath` instead of unescaped `req.URL.Path` for route matching |
| `OHM9996_SERVER_UNESCAPE_PATH_VALUES` | `bool` | `true` | Unescapes URI path parameters before passing to endpoint handlers |
| `OHM9996_SERVER_REMOVE_EXTRA_SLASH` | `bool` | `false` | Removes redundant consecutive slashes from URL path prior to matching |

### Upstream Proxy & Forwarded Headers

| Variable | Type | Default | Description |
|---|---|---|---|
| `OHM9996_SERVER_TRUSTED_PROXIES` | `[]string` | `""` (none) | Comma-separated list of trusted upstream proxy IP addresses or CIDR blocks |
| `OHM9996_SERVER_REMOTE_IP_HEADERS` | `[]string` | `X-Forwarded-For, X-Real-IP` | Comma-separated list of proxy headers inspected to determine real client IP |
| `OHM9996_SERVER_FORWARDED_BY_CLIENT_IP` | `bool` | `true` | Resolves client IP using closest proxy in header chain when trusted |

### CORS Policy Configuration

CORS policies are configured programmatically via blueprint metadata (`rest.NewCorsSetting()` and `meta.WithCors(cors)`). When CORS is enabled without overrides, the server applies the following built-in defaults:

| Setting / Property | Type | Default Value | Description |
|---|---|---|---|
| `AllowOrigins` | `[]string` | `["*"]` | List of allowed CORS origin URL patterns (supports wildcard `*`) |
| `AllowMethods` | `[]string` | `["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"]` | Allowed HTTP methods |
| `AllowHeaders` | `[]string` | `["Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"]` | Allowed request headers in CORS requests |
| `ExposeHeaders` | `[]string` | `["Content-Length", "X-Request-ID"]` | Response headers accessible to browser JavaScript |
| `AllowCredentials` | `bool` | `false` | Allows cookies and HTTP authentication credentials |
| `MaxAge` | `string` | `"86400"` (24 hours) | Preflight response cache time in seconds |

### OpenTelemetry Tracing Configuration

| Variable | Type | Default | Description |
|---|---|---|---|
| `OHM9996_OTEL_ENABLED` | `bool` | `true` | Enables OpenTelemetry distributed tracing and span instrumentation |
| `OHM9996_OTEL_SERVICE_NAME` | `string` | `whatever-service` | Service name reported to OpenTelemetry collectors |
| `OHM9996_OTEL_EXPORTER_OTLP_ENDPOINT` | `string` | `localhost:4317` | Collector host:port address or URL endpoint |
| `OHM9996_OTEL_EXPORTER_OTLP_PROTOCOL` | `string` | `grpc` | Transport protocol (`grpc`, `http`, `http/protobuf`, `http/json`) |
| `OHM9996_OTEL_INSECURE` | `bool` | `true` | Disables TLS verification for OTLP collector connections |
| `OHM9996_OTEL_SAMPLE_RATE` | `float64` | `1.0` | Trace sampling ratio between `0.0` (0%) and `1.0` (100%) |
| `OHM9996_OTEL_SKIP_PATHS` | `string` | `""` | Comma-separated list of URL paths excluded from tracing spans |

### MongoDB Configuration

| Variable | Type | Default | Description |
|---|---|---|---|
| `OHM9996_MONGODB_HOST` | `string` | `""` (Required) | Hostname or IP address of the MongoDB server |
| `OHM9996_MONGODB_PORT` | `int` | `27017` | Network port (1–65535, optional when using `mongodb+srv`) |
| `OHM9996_MONGODB_PROTOCOL` | `string` | `mongodb` | Connection protocol (`mongodb` or `mongodb+srv`) |
| `OHM9996_MONGODB_DATABASE` | `string` | `""` (empty) | Default application database for direct collection access |
| `OHM9996_MONGODB_USERNAME` | `string` | `""` (empty) | Username for authentication (optional, paired with password) |
| `OHM9996_MONGODB_PASSWORD` | `string` | `""` (empty) | Password for authentication (optional, paired with username) |
| `OHM9996_MONGODB_AUTH_SOURCE` | `string` | `""` (empty) | Authentication database name (e.g., `admin`) |
| `OHM9996_MONGODB_APP_NAME` | `string` | `""` (empty) | Application name for connection metadata and diagnostics |
| `OHM9996_MONGODB_UUID_REPRESENTATION` | `string` | `unspecified` | UUID binary representation (`unspecified`, `standard`, `csharpLegacy`, `javaLegacy`, `pythonLegacy`) |
| `OHM9996_MONGODB_ENABLE_PING` | `bool` | `true` | Enables startup connectivity ping verification |
| `OHM9996_MONGODB_CONNECT_TIMEOUT` | `duration` | `10s` | Maximum duration for initial TCP connection establishment |
| `OHM9996_MONGODB_SERVER_SELECTION_TIMEOUT` | `duration` | `5s` | Timeout for cluster server discovery and primary election |
| `OHM9996_MONGODB_SOCKET_TIMEOUT` | `duration` | `10s` | Socket read and write operation timeout |
| `OHM9996_MONGODB_MAX_CONN_IDLE_TIME` | `duration` | `10m` | Maximum idle duration before an unused connection is closed |
| `OHM9996_MONGODB_MAX_POOL_SIZE` | `uint64` | `100` | Maximum number of concurrent connections in the pool |
| `OHM9996_MONGODB_MIN_POOL_SIZE` | `uint64` | `5` | Minimum number of idle connections maintained in the pool |

If mandatory configuration (`HOST`) is missing, or if connection fails after TLS fallback, the client logs descriptive error details and triggers a peaceful termination hook (`exitFunc(0)`).

### Structured Logger Configuration

Structured logging via `pkg/logging` is configured programmatically via `logging.Config` or `logging.DefaultConfig()`. In addition, framework-level HTTP access logging in the REST server is toggled via the `OHM9996_SERVER_ENABLE_ACCESS_LOG` environment variable.

| Field / Option | Type | Default | Description |
|---|---|---|---|
| `Level` | `logging.Level` | `LevelInfo` (`"INFO"`) | Minimum log level threshold (`TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`) |
| `Format` | `logging.Format` | `FormatJSON` (`"json"`) | Log output formatting scheme (`json` or `text`) |
| `AddSource` | `bool` | `false` | When true, includes caller file path and line number in record attributes |
| `DisableTraceCorrelation` | `bool` | `false` | When false, automatically correlates OpenTelemetry `trace_id` and `span_id` |
| `Output` | `io.Writer` | `os.Stdout` | Target output stream for serialized log records |
| `ReplaceAttr` | `ReplaceAttrFunc` | `nil` | Custom attribute transformation and filtering callback |
| `ExitFunc` | `func(int)` | `os.Exit` | Custom termination function executed on `Fatal` and `FatalContext` calls |
| `OHM9996_SERVER_ENABLE_ACCESS_LOG` | `bool` | `true` | Environment variable controlling Gin HTTP access logging middleware |

---

## Consumer Bootstrap

The example below demonstrates how an importing application bootstraps MongoDB connectivity, configures structured logging and CORS policies, defines declarative REST API registrations, and starts the server with signal-driven graceful shutdown:

```go
package main

import (
	"context"
	"time"

	"github.com/nawaphonOHM/whatever/pkg/logging"
	"github.com/nawaphonOHM/whatever/pkg/mongodb"
	"github.com/nawaphonOHM/whatever/pkg/rest"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	ctx := context.Background()

	// 1. Initialize MongoDB connection with zero-boilerplate environment configuration
	mongoClient, err := mongodb.Connect(ctx)
	if err != nil {
		logging.FatalContext(ctx, "Failed to initialize MongoDB", "error", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(shutdownCtx); err != nil {
			logging.ErrorContext(shutdownCtx, "MongoDB disconnect error", "error", err)
		}
	}()

	// 2. Define domain API registration
	itemsColl := mongoClient.Collection("items", "app_db")
	_ = itemsColl

	itemAPI := &rest.RRestAPIRegistration{
		Version: 1,       // Mounts under /api/v1 prefix
		Prefix:  "/items", // Route prefix: /api/v1/items
		Apis: []*rest.ExportableAPI{
			{
				Path:   "",
				Method: rest.GET,
				Handler: func(c rest.Context) rest.Response {
					items := []Item{
						{ID: "1", Name: "Widget"},
						{ID: "2", Name: "Gadget"},
					}
					return rest.OK(items, "Items retrieved successfully")
				},
			},
			{
				Path:   "/:id",
				Method: rest.GET,
				Handler: func(c rest.Context) rest.Response {
					id := c.Param("id")
					if id == "" {
						return rest.BadRequest("INVALID_ID", "Item ID is required")
					}
					return rest.OK(Item{ID: id, Name: "Widget"})
				},
			},
		},
	}

	// 3. Configure CORS policy and assemble blueprint metadata
	cors := rest.NewCorsSetting().
		WithAllowOrigin("http://localhost:3000", "https://app.example.com").
		WithAllowHTTPMethods(rest.GET, rest.POST, rest.PUT, rest.DELETE, rest.OPTIONS)

	meta := rest.NewMeta().
		WithCors(cors)

	blueprint := rest.NewBluePrint().
		WithMeta(meta).
		WithAPIs(itemAPI)

	// 4. Start HTTP server (blocks until SIGINT / SIGTERM signal, then gracefully drains)
	if err := rest.StartREST(blueprint); err != nil {
		logging.FatalContext(ctx, "Server encountered fatal error", "error", err)
	}
}
```

---

## Development Workflows

The project provides standard local development, testing, linting, formatting, and dependency verification targets for contributors and CI pipelines.

### Makefile Targets

The repository includes a dedicated `Makefile` defining verification, test, and maintenance targets:

| Target | Command | Description |
|---|---|---|
| `make` / `make help` | `make help` | Displays list of available Makefile targets with descriptions |
| `make build` | `go build -v ./...` | Verifies compilation of all library packages |
| `make test` | `go test -race -v ./...` | Runs all unit and integration test suites with the race detector enabled |
| `make test-coverage` | `go test -race -coverprofile=coverage.out ./...` | Executes tests and generates an HTML code coverage report (`coverage.html`) |
| `make vet` | `go vet ./...` | Runs standard Go static analysis (`go vet`) across all packages |
| `make lint` | `golangci-lint run` | Executes `golangci-lint` (falls back to `go vet ./...` if `golangci-lint` is not installed) |
| `make tidy` | `go mod tidy && go mod verify` | Tidies and verifies Go module dependencies in `go.mod` and `go.sum` |
| `make all` | `make test vet build` | Runs the full verification pipeline (`test`, `vet`, and `build`) |
| `make clean` | `rm -rf bin tmp coverage.out coverage.html profile.out` | Cleans temporary test coverage, profiling, and build artifact files |

### Testing & Code Coverage

#### Full Test Suite (Makefile)

Execute the complete unit and integration test suite with Go's race detector enabled:
```bash
make test
```

#### Coverage Profiling & Visualization (Makefile)

Generate a code coverage profile and export an interactive HTML visualization:
```bash
# Generate coverage.out and coverage.html
make test-coverage

# View generated coverage report in browser (optional)
xdg-open coverage.html 2>/dev/null || open coverage.html 2>/dev/null || echo "Report generated: coverage.html"
```

#### Targeted & Direct Test Execution (Go Toolchain)

To execute specific packages or subsets of tests directly without running the entire suite (note: these are direct `go test` toolchain commands, not Makefile targets):

```bash
# Run unit tests only for public packages
go test -v ./pkg/...

# Run unit tests only for internal implementation packages
go test -v ./internal/...

# Run targeted package tests (e.g., REST routing or MongoDB client)
go test -v ./pkg/rest/...
go test -v ./pkg/mongodb/...

# Run MongoDB Testcontainers integration tests
go test -v ./pkg/testcontainers/mongodb/...
```

### Static Analysis & Linting

#### Linting & Vetting (Makefile)

Run automated linters and static analyzers to catch issues early:
```bash
# Run golangci-lint (with automatic fallback to go vet)
make lint

# Run go vet explicitly across all packages
make vet
```

#### Linter Auto-Fixing (Direct Tooling)

To automatically fix supported linter issues (direct CLI command, not a Makefile target):
```bash
golangci-lint run --fix
```

### Code Formatting

Ensure all Go source files adhere to standard formatting conventions. The Makefile does not define a `format` target; format code directly using the Go toolchain:

```bash
# Format and simplify Go source files
gofmt -s -w .

# Or format all packages via standard go fmt
go fmt ./...
```

### Dependency Management & Verification

Maintain clean `go.mod` and `go.sum` files and verify cryptographic checksums of dependencies:

```bash
# Tidy unused requirements and verify module checksums via Makefile
make tidy

# Or run individual Go toolchain commands directly (non-Makefile)
go mod tidy
go mod verify
```

### Full Verification Pipeline

Run the complete verification pipeline locally before submitting changes or creating pull requests:
```bash
make all
```

---

## Continuous Integration

The repository includes a GitHub Actions workflow (`.github/workflows/ci.yml`) configured to:
- Enforce Go module integrity (`go mod verify` and dirty check on `go.mod`/`go.sum`).
- Execute static analysis with `golangci-lint` (`--timeout=5m`).
- Run the full test suite across all packages with race detection and coverage collection (`go test -race -v -coverprofile=coverage.out ./...`).
- Verify compilation of all library packages with `go build -v ./...`.
