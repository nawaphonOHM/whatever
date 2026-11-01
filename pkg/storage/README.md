# Cloud Object Storage (`pkg/storage`)

The `storage` package is the public entrypoint for interacting with cloud object storage across Amazon Web Services (AWS) S3, Google Cloud Storage (GCS via S3-compatible XML API), and generic S3-compatible storage backends (e.g., MinIO, LocalStack, Ceph).

It follows the same modular library pattern as `pkg/mongodb` and `pkg/rest`: consuming applications interact with clean public APIs (`pkg/storage`), while configuration, validation, provider detection, and driver options construction remain encapsulated in `internal/storage/*`.

---

## Key Features

- **Unified Storage Contract**: Clean interface (`Storage`) for object lifecycle operations (`PutObject`, `GetObject`, `DeleteObject`, `HeadObject`, `ListObjects`, `Exists`, `PresignGetObject`, `PresignPutObject`, and `Ping`).
- **Zero-Boilerplate Configuration**: Automated configuration parsing and strict validation from `OHM9996_STORAGE_*` environment variables (with fallback to `.env` / `configs/.env`).
- **Structured Cloud Provider Detection**: Detects Google Cloud Storage XML API, AWS S3, or S3-compatible providers from endpoint/bucket URL patterns with structured detection metadata (`Provider`, `Host`, `Rule`, `Bucket`).
- **Diagnostic URL & Credential Sanitization**: Best-effort redaction of sensitive basic auth credentials, session tokens, query strings, and URL fragments from diagnostic log records.
- **Pre-signed URLs**: Generates pre-signed GET and PUT URLs with configurable expiration durations (up to 7 days, falling back to default if 0).
- **Stream Handling & Multipart Uploads**: Efficiently handles seekable streams, buffers non-seekable streams, and automatically coordinates multipart uploads for large payloads.
- **Startup Connectivity Probing**: Non-destructive bucket reachability probe (`HeadBucket` with fallback to single-key `ListObjectsV2` on `AccessDenied`/`Forbidden`/HTTP 403) with diagnostic logging and process exit hooks.
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

## Provider Detection & Endpoint Routing (Internal Architecture)

When resolving target endpoints during client connection, `OHM9996_STORAGE_ENDPOINT` takes precedence over `OHM9996_STORAGE_BUCKET_URL`. URLs without an explicit scheme default to `https://`. Hostnames are normalized (lowercased, trailing dots trimmed), and provider classification is evaluated internally via `detector.Detect(endpoint, bucketURL)` returning structured `detector.Detection` metadata:

```go
type Detection struct {
	Provider contracts.Provider // contracts.ProviderAWS, ProviderGCS, or ProviderS3Compatible
	Host     string             // Normalized inspected hostname
	Rule     string             // Deterministic classification rule matched
	Bucket   string             // Extracted virtual-hosted bucket subdomain (if present)
}
```

### Detection Rules & Provider Classification

Endpoints are evaluated deterministically in the following priority order:

| Rule Constant | Rule Identifier | Matched Provider | Target Pattern & Behavior |
| --- | --- | --- | --- |
| `RuleDefaultAWS` | `"default_aws"` | `ProviderAWS` (`"aws"`) | Triggered when both `Endpoint` and `BucketURL` are empty strings. Uses AWS SDK v2 default regional endpoints. |
| `RuleGCSHostname` | `"gcs_hostname"` | `ProviderGCS` (`"gcs"`) | Matches Google Cloud Storage XML API hostnames ending with or matching standard suffixes (`storage.googleapis.com`, `storage.googleapis`, `stroage.googleapis.com`, `stroage.googleapis`, `commondatastorage.googleapis.com`, `storage.cloud.google.com`), regional pattern `[<bucket>.]storage.<region>.rep.googleapis.com`, or Private Service Connect pattern `[<bucket>.]storage[-<name>].p.googleapis.com` (including `stroage` typo variants). Configures S3 XML interoperability mode with HMAC credentials and path-style addressing. Unrelated Google Cloud hostnames (e.g., `logging.rep.googleapis.com`, `compute.p.googleapis.com`, `bigquery.googleapis.com`) fall back to `s3_compatible_fallback`. |
| `RuleAWSHostname` | `"aws_hostname"` | `ProviderAWS` (`"aws"`) | Matches native AWS S3 endpoints where **both** criteria are met: (1) hostname ends with an AWS domain suffix (`.amazonaws.com`, `.amazonaws.com.cn`, or `.api.aws`), and (2) contains an S3 service label (`s3`, `s3-*`, or `s3express-*`) preceding the AWS suffix, with only permitted tail labels (`dualstack`, `fips`, `vpce`, or a standard AWS region pattern like `us-east-1`) allowed between them. Uses default AWS SDK virtual-hosted addressing (`OHM9996_STORAGE_USE_PATH_STYLE=true` forces path-style). Non-S3 AWS service hostnames (e.g., EC2, ELB, API Gateway) fall back to `s3_compatible_fallback`. |
| `RuleS3CompatibleFallback` | `"s3_compatible_fallback"` | `ProviderS3Compatible` (`"s3-compatible"`) | Fallback rule matched for all other custom endpoints and bucket URLs (e.g., MinIO, LocalStack, Wasabi, Cloudflare R2, Ceph, custom internal endpoints, or non-S3 AWS/GCS service endpoints). Automatically enforces path-style addressing and custom endpoint routing. |

### Virtual-Hosted Bucket Subdomain Extraction

When connecting to multi-tenant cloud providers (AWS S3 and GCS), the internal helper `detector.BucketSubdomain(host)` extracts the virtual-hosted bucket prefix from the hostname if present:
- **AWS S3**: Extracts bucket prefix preceding the S3 service label (e.g., `https://mybucket.s3.us-east-1.amazonaws.com` -> bucket `"mybucket"`, `https://mybucket.s3express-use1-az4.us-east-1.amazonaws.com` -> bucket `"mybucket"`, `https://mybucket.bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com` -> bucket `"mybucket"`). Base service endpoints (e.g., `s3.amazonaws.com`, `vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com`) return an empty bucket string.
- **Google Cloud Storage**: Extracts bucket prefix preceding standard suffixes, regional suffixes, or PSC suffixes (e.g., `https://mybucket.storage.googleapis.com` -> bucket `"mybucket"`, `https://mybucket.storage.us-central1.rep.googleapis.com` -> bucket `"mybucket"`, `https://mybucket.storage-psc.p.googleapis.com` -> bucket `"mybucket"`). Base endpoints (e.g., `storage.googleapis.com`, `storage.us-central1.rep.googleapis.com`) return an empty bucket string.

### Internal Detector Functions

The internal package `internal/storage/detector` provides:
- `detector.Detect(endpoint, bucketURL) (Detection, error)`: Computes full structured detection metadata.
- `detector.DetectProvider(endpoint, bucketURL) (contracts.Provider, error)`: Backward-compatibility wrapper returning the provider enum.

---

## URL Sanitization & Credential Redaction (Internal Helpers)

To safeguard against accidental credential leakage in diagnostic logs and observability pipelines, `internal/storage/client/url_sanitize.go` provides best-effort URL sanitization routines used during client initialization and diagnostic logging:

- **`sanitizeLoggedURL(raw string) string`**:
  - **Best-Effort Userinfo Redaction**: Redacts embedded basic authentication credentials (`user:secretpass@` -> stripped) on valid absolute URLs and standard schemeless hosts.
  - **Strips Query Parameters**: Removes all URL query parameters (`?token=secret...` -> stripped) to prevent leaking temporary SAS tokens, signature parameters, or auth keys.
  - **Strips URL Fragments**: Removes fragment identifiers (`#section` -> stripped).
  - **Preserves URL Structure**: Maintains protocol scheme, hostname, port, and path hierarchy (`https://user:pass@minio.internal:9000/bucket?token=123#frag` -> `https://minio.internal:9000/bucket`).
  - **Schemeless URL Support**: Safely scrubs standard schemeless endpoint strings without prepending unwanted schemes (`admin:secret@minio.corp:9000/data?key=val` -> `minio.corp:9000/data`).
  - **Defensive Error Handling**: Rejects malformed strings or inputs containing invalid whitespace characters (`" \t\r\n"`) by returning `<unparseable>` placeholder without throwing panics. Empty inputs return `""`.
- **`formatLoggedEndpoint(ep *string) string`**:
  - Safely handles pointer endpoints, returning `""` when `ep` is `nil` and `sanitizeLoggedURL(*ep)` otherwise.

### Known Limitations & Security Considerations

- **Schemeless URLs Containing `://` in Query**: In the current implementation, if a schemeless input string contains `://` in its query string (e.g., `user:pass@host:9000/x?r=a://b`), Go's standard `url.Parse` misidentifies the username as a scheme and treats the rest of the URL as opaque data without stripping the `user:pass` credentials. This known edge case is scheduled for code remediation and test verification in step 3.
- **Path-Based Secrets**: Secrets or authentication tokens embedded directly within URL path segments or non-standard token formats outside standard URL userinfo and query parameters are not recognized by URL parsers and are not redacted.

---

## Testing & Verification

The internal storage provider detector and URL sanitization components are validated with unit and matrix test suites:

- **`detector_detect_test.go`**:
  - `TestDetect_RulesAndProviders`: Verifies deterministic rule classification, provider mapping, hostname inspection, and bucket extraction across default AWS, GCS base/virtual-hosted, AWS base, and custom S3-compatible endpoints.
- **`url_sanitize_test.go`**:
  - `TestSanitizeLoggedURL_RedactsCredentialsAndQueries`: Validates credential removal, query string stripping, fragment removal, clean URL preservation, schemeless URL preservation, whitespace error handling (`<unparseable>`), and empty string handling.
  - `TestFormatLoggedEndpoint`: Validates nil-safety and pointer endpoint formatting.
- **Matrix Suites**:
  - `detector_aws_test.go`: Validates AWS S3 endpoint classification (standard endpoints, regional formats, dualstack, FIPS, VPC endpoints, S3 Express, and China regions) and virtual-hosted bucket subdomain extraction.
  - `detector_gcs_test.go`: Validates Google Cloud Storage XML API endpoint classification (standard domains, regional `rep.googleapis.com`, and PSC `p.googleapis.com`) and bucket subdomain extraction.
  - `detector_custom_test.go`: Validates fallback to `s3_compatible_fallback` for custom endpoints and ensures non-S3 Google Cloud and AWS service endpoints (e.g., BigQuery, Compute, EC2, ELB, API Gateway) correctly classify as S3-compatible rather than AWS/GCS.
  - `detector_precedence_test.go`: Validates that `Endpoint` takes precedence over `BucketURL`.
  - `detector_normalize_test.go`: Validates hostname case insensitivity and trailing root dot normalization.
  - `detector_invalid_test.go`: Validates rejection of unsupported schemes (`ftp://`, `postgres://`), out-of-range ports, and whitespace.

---

## Connection Lifecycle, Diagnostics & Startup Probing

### Modular Connection Pipeline (Internal Architecture)

Storage client initialization decomposes connection handling into specialized, single-purpose components in `internal/storage/client/`:

- **`connect.go`**: Top-level coordinator for `storage.Connect` (which loads environment configuration via `config.LoadConfig`) and `storage.ConnectWithConfig`. Coordinates provider detection, configuration validation, SDK client construction, diagnostic/lifecycle logging, and startup probing.
- **`connect_validate.go`**: Validates configuration parameters (`cfg.Validate()`), resolves effective region and endpoint settings, and formats structured lifecycle logs (`initializing storage client` and `storage client initialization complete` at `Info` level).
- **`connect_log_diagnostic.go`**: Emits pre-connection detection diagnostics (`evaluating storage provider detection` at `Debug` level) containing sanitized raw and normalized endpoints, inspected hostname, matched rule, detected provider, and extracted bucket subdomains.
- **`connect_log_compat.go` & `connect_log_compat_format.go`**: Evaluates and logs S3 compatibility settings (`storage client compatibility configuration` at `Debug` level) detailing checksum calculation policies, path-style addressing behaviors, authentication modes, and driver overrides.
- **`connect_probe.go`**: Manages startup bucket reachability probing (`HeadBucket` falling back to `ListObjectsV2`) and process exit hooks (`SetExitFunc`). On probe failure, dispatches abnormal process termination via `central.ExitWithAbnormal`.
- **`connect_client.go`**: Constructs the underlying AWS SDK v2 `s3.Client` and `s3.PresignClient`, binding base endpoints, credentials (static vs. default chain), HTTP transports, and custom functional options.

```
storage.Connect / ConnectWithConfig
  │
  ├── 1. detector.Detect() ────────────────► connect_log_diagnostic.go (Debug: "evaluating storage provider detection")
  ├── 2. validateConfig() ─────────────────► connect_validate.go
  ├── 3. buildClientInstance() ────────────► connect_client.go
  ├── 4. logClientInitAndCompat() ─────────► connect_validate.go (Info: "initializing storage client")
  │                                           connect_log_compat.go (Debug: "storage client compatibility configuration")
  ├── 5. maybeProbeStartup() ──────────────► connect_probe.go (HeadBucket -> ListObjectsV2)
  │                                           (On failure: central.ExitWithAbnormal -> LevelError + exitFunc(1))
  └── 6. logStartupComplete() ─────────────► connect_validate.go (Info: "storage client initialization complete")
```

### Execution Lifecycle

Calling `storage.Connect(ctx, opts...)` executes the following lifecycle:

1. **Configuration Resolution**: Loads environment variables with `OHM9996_STORAGE_*` prefix (with fallback to `.env`/`configs/.env` via `rest/config.Load`) and validates parameter bounds.
2. **Provider Detection & Diagnostic Evaluation**: Classifies provider type (`aws`, `gcs`, `s3-compatible`) from endpoint/bucket URL patterns (with `Endpoint` taking precedence over `BucketURL`) and emits pre-connection debug diagnostics (`evaluating storage provider detection`).
3. **Configuration Validation**: Resolves effective region and endpoint URL, returning validation errors directly if parameters are invalid.
4. **Client Initialization**: Builds the underlying AWS SDK v2 `s3.Client` and `s3.PresignClient`, binding transport options, authentication mode, path style, and flexible checksums.
5. **Lifecycle & Compatibility Diagnostics**: Emits `initializing storage client` (`Info`) and `storage client compatibility configuration` (`Debug`).
6. **Connectivity Probe**: Unless disabled via `storage.WithSkipProbe(true)`, executes a non-destructive reachability probe against the configured default bucket:
   - Issues a `HeadBucket` request bounded by `PingTimeout` (default `5s`).
   - If `HeadBucket` returns `ErrAccessDenied` (mapped from Smithy/SDK API error codes `AccessDenied`, `Forbidden`, or `403`, as well as HTTP 403 status codes), transparently falls back to `ListObjectsV2` with `MaxKeys: 1` to accommodate restrictive IAM policies that grant object listing without bucket metadata read permissions.
7. **Initialization Complete**: Emits `storage client initialization complete` (`Info`) and returns the ready `*Client`.

### Structured Diagnostic & Compatibility Logging Schemas

When structured logging is enabled, storage client initialization emits structured log records at `Info` and `Debug` levels.

#### 1. Provider Detection Evaluation (`Debug` Level)

Emitted by `connect_log_diagnostic.go` with message `"evaluating storage provider detection"` when the logger has `Debug` level enabled:

| Field | Type | Description | Example |
| --- | --- | --- | --- |
| `raw_endpoint` | `string` | Configured endpoint sanitized via `sanitizeLoggedURL` | `"https://storage.googleapis.com"` |
| `raw_bucket_url` | `string` | Configured bucket URL sanitized via `sanitizeLoggedURL` | `""` |
| `inspected_host` | `string` | Normalized hostname inspected for classification | `"storage.googleapis.com"` |
| `detected_provider` | `string` | Classified provider identifier (`"aws"`, `"gcs"`, `"s3-compatible"`) | `"gcs"` |
| `detection_rule` | `string` | Classification rule matched (`"default_aws"`, `"aws_hostname"`, `"gcs_hostname"`, `"s3_compatible_fallback"`) | `"gcs_hostname"` |
| `extracted_bucket` | `string` | Virtual-hosted bucket subdomain (if present in multi-tenant hostnames) | `""` |
| `normalized_endpoint` | `string` | Fully resolved and sanitized base endpoint URL | `"https://storage.googleapis.com"` |

#### 2. Lifecycle Progress (`Info` Level)

Emitted by `connect_validate.go` with messages `"initializing storage client"` and `"storage client initialization complete"`:

| Field | Type | Description | Example |
| --- | --- | --- | --- |
| `provider` | `string` | Target cloud provider (`"aws"`, `"gcs"`, `"s3-compatible"`) | `"gcs"` |
| `bucket` | `string` | Target default bucket name | `"my-app-data"` |
| `region` | `string` | AWS or S3 region identifier (defaults to `us-east-1` if unset) | `"auto"` (when configured) |
| `endpoint` | `string` | Sanitized base endpoint URL (empty string `""` for default AWS) | `"https://storage.googleapis.com"` |
| `use_path_style` | `bool` | Whether path-style addressing is active | `true` |
| `checksum_mode` | `string` | Formatted checksum mode (`"when_required"` or `"when_supported"`) | `"when_required"` |

#### 3. Compatibility Configuration Diagnostics (`Debug` Level)

Emitted by `connect_log_compat.go` with message `"storage client compatibility configuration"` when `Debug` level is active:

| Field | Type | Description |
| --- | --- | --- |
| `provider` | `string` | Target cloud provider (`"aws"`, `"gcs"`, `"s3-compatible"`). |
| `request_checksum_calculation` | `string` | AWS SDK request checksum calculation mode (`"WHEN_REQUIRED"` or `"WHEN_SUPPORTED"`). |
| `response_checksum_validation` | `string` | AWS SDK response checksum validation mode (`"WHEN_REQUIRED"` or `"WHEN_SUPPORTED"`). |
| `checksum_behavior_change` | `string` | Behavioral reason comparing effective mode with provider defaults (see constants below). |
| `use_path_style` | `bool` | Whether path-style addressing is enabled after all options/overrides are evaluated. |
| `path_style_behavior_change` | `string` | Behavioral reason reflecting effective addressing style and provider (see constants below). |
| `disable_log_output_checksum_validation_skipped` | `bool` | Internal flag set to `true` to suppress noisy SDK internal checksum skip warnings. |
| `auth_mode` | `string` | Authentication mode: `"static_credentials"` (when access/secret keys are provided) or `"default_credential_chain"`. |
| `has_driver_overrides` | `bool` | Boolean flag indicating whether caller supplied custom SDK options via `WithDriverOptions`. |

##### Compatibility Behavior Constants

- **`checksum_behavior_change`**:
  - `BehaviorChecksumSuppressed` (`"suppress_unsolicited_trailers_for_gcs"`): Default for GCS (`WHEN_REQUIRED`), suppressing unsolicited `aws-chunked` trailing checksum headers that cause GCS S3 XML API signature mismatch errors (`SignatureDoesNotMatch`).
  - `BehaviorChecksumEnforced` (`"enforce_when_required_for_compat"`): Default for S3-compatible providers (`WHEN_REQUIRED`), preventing trailer checksum mismatches on non-AWS backends.
  - `BehaviorChecksumStandard` (`"standard_aws_defaults"`): Default for native AWS (`WHEN_SUPPORTED`).
  - `BehaviorChecksumExplicit` (`"explicit_checksum_configuration"`): Emitted when the effective request checksum calculation mode differs from the provider default (e.g. `WHEN_REQUIRED` on AWS, or any mode other than `WHEN_REQUIRED` on non-AWS backends). Note that `WithDriverOptions` mutators are applied last across all providers and can set arbitrary checksum modes, while the `AWS_REQUEST_CHECKSUM_CALCULATION` environment variable only takes effect on AWS when using the default credential chain (the static credentials path bypasses AWS environment loading, and non-AWS providers unconditionally enforce `WHEN_REQUIRED` during provider compatibility setup). An override that matches the provider default retains the default label.
- **`path_style_behavior_change`**:
  - `BehaviorPathStyleStandard` (`"standard_virtual_hosted"`): Effective `UsePathStyle` is `false` (standard AWS virtual-hosted addressing, or non-AWS endpoints where driver options explicitly set `UsePathStyle = false`).
  - `BehaviorPathStyleConfigured` (`"configured_path_style"`): Effective `UsePathStyle` is `true` on AWS (via `OHM9996_STORAGE_USE_PATH_STYLE=true` or via driver options `WithDriverOptions`).
  - `BehaviorPathStyleEnforced` (`"enforced_path_style_for_provider"`): Effective `UsePathStyle` is `true` on GCS or S3-compatible backends (standard default for non-AWS backends, unless overridden via driver options).

### Error Handling & Process Exit Semantics

- **Missing Mandatory Configuration**: If required connection settings (e.g. `OHM9996_STORAGE_BUCKET`) are missing when loading via `storage.Connect` / `config.LoadConfig`, the client logs a diagnostic message and dispatches graceful termination (`central.ExitWithGraceful`), invoking `exitFunc(0)` (exit code 0).
- **Configuration Validation & Build Errors**: Syntax errors, out-of-range ports, invalid timeouts, or unparseable URLs passed to `storage.ConnectWithConfig` return an error directly to the caller and **never trigger any exit hooks**.
- **Startup Connectivity Probe Failure**: When startup probing is enabled and the probe fails (e.g., `ErrBucketNotFound`, `context.DeadlineExceeded`, or network unreachable), `handleConnectionError` dispatches an abnormal exit event via `central.ExitWithAbnormal(ctx, "failed to connect to storage; exiting abnormally", err, nil)`:
  - The central logging worker logs the error with severity `slog.LevelError` (exit flag `ExitAbnormal`). *Note: URL sanitization applies specifically to endpoint and bucket URL log attributes; raw error text attached to the abnormal exit record is NOT sanitized and may contain hostnames or URLs embedded by the SDK/network layer.*
  - When the context carries a call stack (as decorated during `storage.Connect` / `ConnectWithConfig`), the worker logs an additional `"Call stack snapshot"` record at `Error` level containing the stack trace.
  - The central logging worker flushes all asynchronous log queues.
  - The worker invokes the configured exit hook `exitFunc(1)` (default `os.Exit(1)`), which terminates the process synchronously.
  - The underlying error is returned to the caller only if the exit hook is overridden via `SetExitFunc` (or if a custom exit hook returns without terminating the process).
- **Exit Hook Interception**: Process exit behavior can be intercepted during testing or within supervisor processes using `storage.SetExitFunc`:

```go
prevExit := storage.SetExitFunc(func(code int) {
	// Intercept exit code (e.g., 0 for graceful, 1 for abnormal probe failure)
})
defer storage.SetExitFunc(prevExit)
```

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
| `ErrAccessDenied` | Storage service rejected access (mapped from API error codes `AccessDenied`, `Forbidden`, `403` or HTTP 403) |
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
