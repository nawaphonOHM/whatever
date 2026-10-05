# Firestore Testcontainers (`pkg/testcontainers/firestore`)

The `pkg/testcontainers/firestore` package provides a lightweight, ergonomic Testcontainers wrapper for running Google Cloud Firestore emulator containers in automated test suites and integration workflows.

Built on top of the official Testcontainers Google Cloud module (`github.com/testcontainers/testcontainers-go/modules/gcloud/firestore`), it encapsulates container lifecycle management, port mapping, and configuration options behind clean public APIs and type aliases.

---

## Key Features

- **Automated Lifecycle Management**: Spin up and tear down isolated Firestore emulator instances via `Run` and `Terminate`.
- **Flexible Functional Options**: Easily configure container images, Google Cloud project IDs, datastore compatibility mode, environment variables, and raw Testcontainers customizers.
- **Connection & Endpoint Helpers**: Convenient accessors for resolved container connection endpoints (`URI`, `Host`, `Port`), project metadata (`ProjectID`, `DatastoreMode`), and underlying raw container handles (`RawContainer`).
- **Defensive Error Handling**: Safe receiver checks on all operations returning descriptive sentinel errors (`ErrNilContainer`, `ErrContainerNotRunning`).

---

## Constants & Sentinel Errors

### Default Configuration Constants

| Constant | Value | Description |
|---|---|---|
| `DefaultImage` | `"gcr.io/google.com/cloudsdktool/google-cloud-cli:emulators"` | Default container image for the Google Cloud Firestore emulator |
| `DefaultPort` | `"8080/tcp"` | Default internal TCP port exposed by the Firestore emulator |
| `DefaultProjectID` | `"test-project"` | Default Google Cloud project ID assigned to the emulator container |

### Sentinel Errors

| Error | Description |
|---|---|
| `ErrNilContainer` | Returned when attempting to perform operations on a `nil` container instance |
| `ErrContainerNotRunning` | Returned when attempting to access endpoints on an unstarted or terminated container |

---

## Functional Options

The package provides functional options (`Option`) to customize container image, project ID, datastore mode, environment variables, and raw Testcontainers options:

| Option | Signature | Description | Default |
|---|---|---|---|
| `WithImage` | `WithImage(image string) Option` | Overrides the default emulator Docker image | `firestore.DefaultImage` (`"gcr.io/google.com/cloudsdktool/google-cloud-cli:emulators"`) |
| `WithProjectID` | `WithProjectID(projectID string) Option` | Sets the Google Cloud project ID for the emulator instance | `firestore.DefaultProjectID` (`"test-project"`) |
| `WithDatastoreMode` | `WithDatastoreMode(enabled ...bool) Option` | Runs the Firestore emulator in Datastore compatibility mode | `false` |
| `WithEnv` | `WithEnv(key, value string) Option` | Injects custom environment variables into the emulator container | `nil` |
| `WithContainerOptions` | `WithContainerOptions(opts ...testcontainers.ContainerCustomizer) Option` | Appends underlying Testcontainers customizers | `nil` |

Helper functions `firestore.DefaultOptions()` and `firestore.NewOptions(opts...)` are also available for programmatic options inspection.

---

## Connection Getters & Container Methods

The `*firestore.Container` handle provides methods to inspect connection endpoints and manage lifecycle:

| Method | Return Type | Description |
|---|---|---|
| `container.URI(ctx)` | `(string, error)` | Returns the connection URI (`host:port`) for the running Firestore container |
| `container.Host(ctx)` | `(string, error)` | Returns the host IP or hostname where the Firestore emulator is accessible |
| `container.Port(ctx)` | `(int, error)` | Returns the mapped external TCP port as an `int` |
| `container.ProjectID()` | `string` | Returns the configured Google Cloud Project ID |
| `container.DatastoreMode()` | `bool` | Returns whether the emulator is running in Datastore compatibility mode |
| `container.RawContainer()` | `*tcfirestore.Container` | Returns the underlying Testcontainers Firestore container handle |
| `container.Terminate(ctx)` | `error` | Stops and deletes the running Firestore container |

---

## Usage Examples

### 1. Basic Ephemeral Container

```go
package integration_test

import (
	"context"
	"log"

	tcfirestore "github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/firestore"
)

func main() {
	ctx := context.Background()

	// Spin up Firestore emulator with default configuration
	container, err := tcfirestore.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to start Firestore container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}()

	// Retrieve emulator endpoint (e.g., "127.0.0.1:49152")
	uri, err := container.URI(ctx)
	if err != nil {
		log.Fatalf("Failed to get container URI: %v", err)
	}

	log.Printf("Firestore emulator running at: %s (Project: %s)", uri, container.ProjectID())
}
```

### 2. Integration Testing with `testing.TB`

```go
//go:build testcontainers

package mytest

import (
	"context"
	"testing"

	tcfirestore "github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/firestore"
)

func TestFirestoreIntegration(t *testing.T) {
	ctx := context.Background()

	// Launch emulator with custom project ID and datastore mode
	container, err := tcfirestore.Run(ctx,
		tcfirestore.WithProjectID("my-integration-project"),
		tcfirestore.WithDatastoreMode(false),
	)
	if err != nil {
		t.Fatalf("Failed to start Firestore emulator: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("Failed to terminate container: %v", err)
		}
	})

	uri, err := container.URI(ctx)
	if err != nil {
		t.Fatalf("Failed to resolve container URI: %v", err)
	}

	t.Logf("Running tests against Firestore emulator at %s", uri)
}
```
