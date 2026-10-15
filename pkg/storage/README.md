# Cloud Object Storage (`pkg/storage`)

The `storage` package is the public entrypoint for interacting with cloud object storage across Amazon Web Services (AWS) S3, Google Cloud Storage (GCS via S3-compatible XML API), and generic S3-compatible storage backends (e.g., MinIO, LocalStack, Ceph).

It follows the same modular library pattern as `pkg/mongodb` and `pkg/rest`: consuming applications interact with clean public APIs (`pkg/storage`), while configuration, validation, provider detection, and driver options construction remain encapsulated in `internal/storage/*`.

---

## Key Features

- **Unified Storage Contract**: Clean interface (`Storage`) for object lifecycle operations (`PutObject`, `GetObject`, `DeleteObject`, `HeadObject`, `ListObjects`, `Exists`, `PresignGetObject`, `PresignPutObject`, and `Ping`).
- **Zero-Boilerplate Configuration**: Automated configuration parsing and strict validation from `OHM9996_STORAGE_*` environment variables (with fallback to `.env` / `configs/.env`).
- **Automatic Cloud Provider Detection**: Detects Google Cloud Storage XML API, AWS S3, or S3-compatible providers from endpoint/bucket URL patterns with endpoint precedence.
- **Pre-signed URLs**: Generates pre-signed GET and PUT URLs with configurable expiration durations (up to 7 days, falling back to default if 0).
- **Stream Handling & Multipart Uploads**: Efficiently handles seekable streams, buffers non-seekable streams, and automatically coordinates multipart uploads for large payloads.
- **Startup Connectivity Probing**: Non-destructive bucket reachability probe (`HeadBucket` with fallback to single-key `ListObjectsV2`) with diagnostic logging and process exit hooks.
- **Defensive Error Handling**: Safe receiver checks on all operations returning standardized sentinel errors (`ErrObjectNotFound`, `ErrBucketNotFound`, `ErrAccessDenied`, `ErrNilClient`, `ErrNilConfig`, `ErrInvalidExpiry`).

---

## Configuration

The client loads and validates configuration from environment variables when initialized. Importing applications configure storage connectivity via standardized `OHM9996_STORAGE_*` environment variables (with fallback to `.env` or `configs/.env` via `rest/config.Load`).

| Environment Variable | Description | Required / Default |
| --- | --- | --- |
| `OHM9996_STORAGE_BUCKET` | Default bucket name (3 to 222 characters) | **Required** |
| `OHM9996_STORAGE_REGION` | AWS region or S3-compatible region identifier | `us-east-1` |
| `OHM9996_STORAGE_ACCESS_KEY_ID` | Access key ID for static authentication | `""` (empty) |
| `OHM9996_STORAGE_SECRET_ACCESS_KEY` | Secret access key for static authentication | `""` (empty) |
| `OHM9996_STORAGE_SESSION_TOKEN` | AWS session token for temporary credentials | `""` (empty) |
| `OHM9996_STORAGE_ENDPOINT` | Custom service endpoint URL or hostname (takes precedence over `BUCKET_URL`) | `""` (empty) |
| `OHM9996_STORAGE_BUCKET_URL` | Custom bucket URL (used if `ENDPOINT` is unset) | `""` (empty) |
| `OHM9996_STORAGE_PORT` | Custom port number (1–65535), applied only when endpoint/bucket URL lacks an explicit port | `0` (default port for protocol) |
| `OHM9996_STORAGE_USE_PATH_STYLE` | Force path-style addressing (`true`/`false`) | `false` |
| `OHM9996_STORAGE_INSECURE` | Force endpoint/bucket URL scheme to `http` (`true`/`false`) | `false` |
| `OHM9996_STORAGE_DISABLE_SSL` | Force endpoint/bucket URL scheme to `http` (`true`/`false`) | `false` |
| `OHM9996_STORAGE_CONNECT_TIMEOUT` | Maximum duration bounding TCP dial and TLS handshake on transport client (ignored if `WithHTTPClient` is provided) | `10s` |
| `OHM9996_STORAGE_REQUEST_TIMEOUT` | Per-request timeout bounding header arrival for `GetObject` and whole transfer duration for `PutObject` (including multipart) | `30s` |
| `OHM9996_STORAGE_PING_TIMEOUT` | Maximum duration for bucket connectivity probes | `5s` |
| `OHM9996_STORAGE_PRESIGN_EXPIRY` | Default pre-signed URL expiration duration (maximum 7 days; used when `ExpiresIn` is `0`) | `15m` |

> **Note on TLS & Insecure Options**: `OHM9996_STORAGE_INSECURE` and `OHM9996_STORAGE_DISABLE_SSL` explicitly force the target endpoint or bucket URL scheme to `http://`. They do not bypass TLS certificate verification on `https://` endpoints.

---

## Provider Detection & Endpoint Routing

When resolving target endpoints, `OHM9996_STORAGE_ENDPOINT` takes precedence over `OHM9996_STORAGE_BUCKET_URL`. URLs without an explicit scheme default to `https://`. Hostnames are normalized (lowercased, trailing dots trimmed), and provider classification is evaluated in the following order:

1. **Google Cloud Storage (`ProviderGCS` / `"gcs"`)**:
   - Matches when the hostname equals or ends with `.storage.googleapis.com`, `.storage.googleapis`, `.stroage.googleapis.com`, or `.stroage.googleapis` (e.g., `https://storage.googleapis.com`, `https://mybucket.storage.googleapis.com`).
   - Connects to Google Cloud Storage via the S3-compatible XML API using HMAC credentials and path-style addressing.

2. **AWS S3 (`ProviderAWS` / `"aws"`)**:
   - Matches when any dot-separated hostname label equals `s3` or starts with `s3-` (e.g., `s3.amazonaws.com`, `bucket.s3.us-east-1.amazonaws.com`, `s3.local:9000`), or when the host equals or ends with `.amazonaws.com` / `amazonaws.com`.
   - Also used as the default provider when neither endpoint nor bucket URL is specified.
   - Uses native AWS SDK v2 defaults. By default, native `amazonaws.com` / `amazonaws.com.cn` endpoints use the SDK's virtual-hosted routing; non-Amazon target hosts (including AWS-classified `s3.local:9000`) use path-style addressing. `OHM9996_STORAGE_USE_PATH_STYLE=true` forces path-style addressing for any provider.

3. **S3-Compatible Storage (`ProviderS3Compatible` / `"s3-compatible"`)**:
   - Matches any other custom endpoint or bucket URL (e.g., `http://localhost:9000` for MinIO, `http://localhost:4566` for LocalStack, `https://minio.internal.net`, `https://mys3bucket.example.com`).
   - Automatically enables path-style bucket addressing and custom base endpoint routing.

---

## Connection Lifecycle & Startup Probing

Calling `storage.Connect(ctx, opts...)` executes the following lifecycle:

1. **Configuration Resolution**: Loads environment variables with `OHM9996_STORAGE_*` prefix (with fallback to `.env`/`configs/.env` via `rest/config.Load`) and applies defaults.
2. **Provider Detection**: Classifies provider type (`gcs`, `aws`, `s3-compatible`) from endpoint/bucket URL patterns (with `Endpoint` taking precedence over `BucketURL`).
3. **Client Initialization**: Builds the underlying AWS SDK v2 `s3.Client` and `s3.PresignClient`.
4. **Connectivity Probe**: Unless disabled via `storage.WithSkipProbe(true)`, executes a non-mutating probe against the configured default bucket:
   - Issues a `HeadBucket` request bounded by `PingTimeout`.
   - If `HeadBucket` returns `ErrAccessDenied` (HTTP 403), transparently falls back to `ListObjectsV2` with `MaxKeys: 1` to accommodate restrictive IAM policies.
5. **Error & Exit Handling**:
   - **Missing mandatory configuration** (`OHM9996_STORAGE_BUCKET` when loading via `storage.Connect` / `config.LoadConfig`) logs a diagnostic message and triggers **graceful termination (`exitFunc(0)`)**.
   - **Startup connectivity probe failure** (when probing is enabled) logs diagnostic details with a call stack snapshot and triggers **abnormal termination (`exitFunc(1)`)**.
   - **Configuration validation errors** (e.g., invalid URLs, invalid timeouts, port out of bounds) and client construction errors passed to `storage.ConnectWithConfig` return an error directly and **do not trigger any exit hook**.

---

## Caller Contracts & Streaming Semantics

- **Stream Ownership on `GetObject`**: The caller owns `GetObjectOutput.Body` and **must close it** (`defer getOutput.Body.Close()`) when reading completes to release the underlying network connection and context. Alternatively, use `storage.ReadAllBody(getOutput)` which automatically reads and closes the stream.
- **Request Timeout Scope**:
  - For `GetObject`: `RequestTimeout` only bounds the wait for metadata and initial response headers arrival; the caller's context controls the full body streaming duration.
  - For `PutObject` (including multipart uploads): `RequestTimeout` bounds the **entire upload operation**. For large objects, configure `OHM9996_STORAGE_REQUEST_TIMEOUT` generously or use pre-signed URLs for direct client transfers.
- **Pre-signed PUT Headers**: When `ContentType` is provided in `PresignPutObjectInput`, it is included in the signed request headers. The uploading client **must send the returned `Headers` (notably `Content-Type`) unchanged** in the HTTP PUT request.
- **Pre-signed Expiration Fallback**: In `PresignGetObject` and `PresignPutObject`, passing `ExpiresIn: 0` uses the configured default expiration (`OHM9996_STORAGE_PRESIGN_EXPIRY`, default `15m`). Negative durations or durations exceeding 7 days return `ErrInvalidExpiry`.

---

## Sentinel Errors & Constants

### Supported Providers

| Constant | Value | Description |
| --- | --- | --- |
| `ProviderAWS` | `"aws"` | Native Amazon Web Services S3 |
| `ProviderGCS` | `"gcs"` | Google Cloud Storage XML API (S3-compatible) |
| `ProviderS3Compatible` | `"s3-compatible"` | Generic S3-compatible object storage (MinIO, LocalStack, etc.) |

### Sentinel Errors

| Error | Description |
| --- | --- |
| `ErrObjectNotFound` | The requested object key does not exist |
| `ErrBucketNotFound` | The target bucket does not exist |
| `ErrNilClient` | Operation attempted on an uninitialized / `nil` client |
| `ErrNilConfig` | Configuration passed is `nil` |
| `ErrAccessDenied` | Storage service rejected access (HTTP 403) |
| `ErrInvalidExpiry` | Pre-signed URL expiration is negative or exceeds 7 days |

---

## Functional Options

The package provides functional options (`Option`) to customize client behavior:

| Option | Signature | Description |
| --- | --- | --- |
| `WithSkipProbe` | `WithSkipProbe(skip bool) Option` | Controls whether the startup bucket connectivity probe is skipped |
| `WithHTTPClient` | `WithHTTPClient(httpClient *http.Client) Option` | Overrides the HTTP client used by the underlying AWS SDK client |
| `WithDriverOptions` | `WithDriverOptions(opts ...func(*s3.Options)) Option` | Appends raw AWS SDK v2 `s3.Options` mutators |
| `WithPresignOptions` | `WithPresignOptions(opts ...func(*s3.PresignOptions)) Option` | Appends raw AWS SDK v2 `PresignOptions` mutators |

---

## Usage Examples

### 1. Basic Connection & Upload / Download

```go
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func main() {
	ctx := context.Background()

	// Initialize storage client from environment variables
	client, err := storage.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Upload an object
	content := []byte("Hello, cloud object storage!")
	putOutput, err := client.PutObject(ctx, storage.PutObjectInput{
		Key:         "documents/hello.txt",
		Body:        bytes.NewReader(content),
		ContentType: "text/plain",
		Metadata: map[string]string{
			"author": "ohm",
		},
	})
	if err != nil {
		log.Fatalf("Failed to upload object: %v", err)
	}
	fmt.Printf("Uploaded object, ETag: %s\n", putOutput.ETag)

	// Download the object
	getOutput, err := client.GetObject(ctx, storage.GetObjectInput{
		Key: "documents/hello.txt",
	})
	if err != nil {
		log.Fatalf("Failed to download object: %v", err)
	}

	// Read and close body using the helper utility (or defer getOutput.Body.Close())
	data, err := storage.ReadAllBody(getOutput)
	if err != nil {
		log.Fatalf("Failed to read object body: %v", err)
	}
	fmt.Printf("Downloaded content: %s\n", string(data))
}
```

### 2. Convenience Helpers (`NewPutObjectBytes`)

```go
input := storage.NewPutObjectBytes("my-bucket", "images/logo.png", imageData, "image/png")
_, err := client.PutObject(ctx, input)
```

### 3. Pre-signed URLs for Upload & Download

```go
// Pre-signed download URL (valid for 30 minutes; passing 0 uses configured default)
getPresigned, err := client.PresignGetObject(ctx, storage.PresignGetObjectInput{
	Key:       "reports/annual-2026.pdf",
	ExpiresIn: 30 * time.Minute,
})
if err != nil {
	log.Fatalf("Failed to presign GET: %v", err)
}
fmt.Printf("Download URL: %s\n", getPresigned.URL)

// Pre-signed upload URL (valid for 15 minutes; uploader must send returned Headers unchanged)
putPresigned, err := client.PresignPutObject(ctx, storage.PresignPutObjectInput{
	Key:         "uploads/avatar.jpg",
	ContentType: "image/jpeg",
	ExpiresIn:   15 * time.Minute,
})
if err != nil {
	log.Fatalf("Failed to presign PUT: %v", err)
}
fmt.Printf("Upload URL: %s (Method: %s)\n", putPresigned.URL, putPresigned.Method)
```

### 4. Listing Objects & Checking Existence

```go
// Check existence
exists, err := client.Exists(ctx, "", "documents/hello.txt")
if err != nil {
	log.Fatalf("Error checking existence: %v", err)
}
fmt.Printf("File exists: %t\n", exists)

// List objects with prefix
listOutput, err := client.ListObjects(ctx, storage.ListObjectsInput{
	Prefix:    "documents/",
	Delimiter: "/",
	MaxKeys:   50,
})
if err != nil {
	log.Fatalf("Failed to list objects: %v", err)
}

for _, obj := range listOutput.Objects {
	fmt.Printf("Key: %s (Size: %d bytes, Modified: %s)\n", obj.Key, obj.Size, obj.LastModified)
}
for _, prefix := range listOutput.Prefixes {
	fmt.Printf("Subfolder prefix: %s\n", prefix)
}
```

### 5. Manual Config & Readiness Probes

```go
cfg := storage.DefaultConfig()
cfg.Bucket = "custom-bucket"
cfg.Region = "ap-southeast-1"
cfg.Endpoint = "http://localhost:9000"

client, err := storage.ConnectWithConfig(ctx, cfg, storage.WithSkipProbe(false))
if err != nil {
	log.Fatalf("Failed to connect: %v", err)
}

// On-demand health check / readiness probe
if err := client.Ping(ctx); err != nil {
	log.Printf("Storage readiness probe failed: %v", err)
}
```
