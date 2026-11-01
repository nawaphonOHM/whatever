# Storage Compatibility, Provider Detection & Security Sanitization Release Notes

## Executive Summary

This release introduces structured cloud storage provider detection, comprehensive S3 client compatibility diagnostics, safe URL credential sanitization, and extracted probe/lifecycle handling for cloud object storage connections. It makes existing checksum calculation policies and path-style addressing compatibility behaviors fully observable, helping troubleshoot signature mismatch errors (`SignatureDoesNotMatch`) and operational ambiguities when communicating across AWS S3, Google Cloud Storage (GCS) S3-compatible XML APIs, MinIO, and generic S3-compatible endpoints.

All changes maintain **100% backward compatibility** with existing public APIs in `pkg/storage`.

---

## 1. Before vs After Comparison

The following analysis compares behavior in `main` versus `HEAD`:

### Architectural & Operational Comparison

| Capability | Before (`main`) | After (`HEAD`) |
| :--- | :--- | :--- |
| **Provider Detection** | Returned only `contracts.Provider` enum via `DetectProvider`; classification reasoning, inspected hostname, and virtual-hosted bucket subdomains were discarded. | Returns rich `detector.Detection` metadata (`Provider`, `Host`, `Rule`, `Bucket`). `DetectProvider` is retained as a backward-compatible wrapper. |
| **Endpoint Logging** | Raw endpoints logged directly at `INFO`; embedded basic auth credentials (`user:pass@host`), SAS tokens, or signature query params risked credential leakage. | Strict URL sanitization via `urlutil.SanitizeURL`; redacts userinfo, queries, and fragments across all logs. Standard and schemeless URLs are supported without injecting undesired default schemes. |
| **Malformed URL Handling** | Malformed inputs (single-slash schemes like `https:/host` or path pseudo-schemes like `https://host/a/b://u:p@x`) were accepted and logged raw at `INFO`. Parser errors on whitespace or invalid ports could echo raw input (e.g. invalid port error wrapping `:SECRETPW`). | Malformed inputs fail closed to `<unparseable>` in logs and are rejected at detection/validation with generic error messages (`invalid URL format`, `invalid scheme format`, `invalid path credentials`) without echoing credentials. |
| **S3 Compatibility Diagnostics** | Checksum calculation policies and path-style addressing rules were applied internally, but runtime compatibility configurations were not observable via diagnostic logs. | Granular `DEBUG`-level diagnostics emit `evaluating storage provider detection` and `storage client compatibility configuration` detailing checksum calculation/validation modes, path-style routing, credential resolution, and driver overrides. |
| **Connection Pipeline** | Network probing lived in `probe.go`, with startup probe gating and abnormal exit handling inline in `connect.go` (along with `connect_validate.go` and `connect_client.go`), without dedicated diagnostic logging modules. | Startup probe gating and abnormal exit handling (`handleConnectionError`, `probeStartupClient`, `maybeProbeStartup`, `SetExitFunc`) were extracted verbatim into `connect_probe.go` (with network probing remaining in `probe.go`), and dedicated diagnostic logging modules were added (`connect_log_diagnostic.go`, `connect_log_compat.go`, `connect_log_compat_format.go`). |
| **Public API Signatures** | Exported API in `pkg/storage`. | **Zero breaking changes**: all exported functions, types, and option builders retain identical signatures. |
| **Test Matrix Verification** | Baseline unit tests. | Comprehensive matrix test suites covering multi-cloud provider classification, GCS trailer suppression, AWS vs custom path-style addressing, security redaction, and IPv6 endpoints. |

---

### Code Implementation Comparison

#### Provider Detection

**Before (`main`):**
```go
// Only returned provider enum; classification reasoning was lost
func DetectProvider(endpoint, bucketURL string) (contracts.Provider, error) {
    target := selectTargetURL(endpoint, bucketURL)
    if target == "" {
        return contracts.ProviderAWS, nil
    }
    // Heuristic returning only contracts.Provider
    ...
}
```

**After (`HEAD`):**
```go
// Returns structured Detection metadata
func Detect(endpoint, bucketURL string) (Detection, error) {
    target := selectTargetURL(endpoint, bucketURL)
    if target == "" {
        return Detection{Provider: contracts.ProviderAWS, Rule: RuleDefaultAWS}, nil
    }
    u, err := ParseTarget(target)
    if err != nil {
        return Detection{}, err
    }
    return classifyDetection(normalizeHost(u.Hostname())), nil
}

// Retained backward-compatible wrapper
func DetectProvider(endpoint, bucketURL string) (contracts.Provider, error) {
    d, err := Detect(endpoint, bucketURL)
    if err != nil {
        return "", err
    }
    return d.Provider, nil
}
```

---

#### Structured Diagnostic & Lifecycle Logging

The following comparisons show actual log records produced by `main` and `HEAD` for identical configurations (using default `Region: "us-east-1"`, timestamp omitted for clarity):

##### 1. MinIO / S3-Compatible Endpoint with Credentials (`Endpoint: "http://admin:secret@minio.internal:9000"`, `Bucket: "my-bucket"`)

**Before (`main`):**
```json
{"level":"INFO","msg":"initializing storage client","provider":"s3-compatible","bucket":"my-bucket","region":"us-east-1","endpoint":"http://admin:secret@minio.internal:9000","use_path_style":true}
{"level":"INFO","msg":"storage client initialization complete","provider":"s3-compatible","bucket":"my-bucket"}
```

**After (`HEAD`):**
```json
{"level":"DEBUG","msg":"evaluating storage provider detection","raw_endpoint":"http://minio.internal:9000","raw_bucket_url":"","inspected_host":"minio.internal","detected_provider":"s3-compatible","detection_rule":"s3_compatible_fallback","extracted_bucket":"","normalized_endpoint":"http://minio.internal:9000"}
{"level":"INFO","msg":"initializing storage client","provider":"s3-compatible","bucket":"my-bucket","region":"us-east-1","endpoint":"http://minio.internal:9000","use_path_style":true,"checksum_mode":"when_required"}
{"level":"DEBUG","msg":"storage client compatibility configuration","provider":"s3-compatible","request_checksum_calculation":"WHEN_REQUIRED","response_checksum_validation":"WHEN_REQUIRED","checksum_behavior_change":"enforce_when_required_for_compat","use_path_style":true,"path_style_behavior_change":"enforced_path_style_for_provider","disable_log_output_checksum_validation_skipped":true,"auth_mode":"default_credential_chain","has_driver_overrides":false}
{"level":"INFO","msg":"storage client initialization complete","provider":"s3-compatible","bucket":"my-bucket","region":"us-east-1","endpoint":"http://minio.internal:9000","use_path_style":true,"checksum_mode":"when_required"}
```

##### 2. Google Cloud Storage XML Endpoint (`Endpoint: "https://storage.googleapis.com"`, `Bucket: "my-bucket"`)

**Before (`main`):**
```json
{"level":"INFO","msg":"initializing storage client","provider":"gcs","bucket":"my-bucket","region":"us-east-1","endpoint":"https://storage.googleapis.com","use_path_style":true}
{"level":"INFO","msg":"storage client initialization complete","provider":"gcs","bucket":"my-bucket"}
```

**After (`HEAD`):**
```json
{"level":"DEBUG","msg":"evaluating storage provider detection","raw_endpoint":"https://storage.googleapis.com","raw_bucket_url":"","inspected_host":"storage.googleapis.com","detected_provider":"gcs","detection_rule":"gcs_hostname","extracted_bucket":"","normalized_endpoint":"https://storage.googleapis.com"}
{"level":"INFO","msg":"initializing storage client","provider":"gcs","bucket":"my-bucket","region":"us-east-1","endpoint":"https://storage.googleapis.com","use_path_style":true,"checksum_mode":"when_required"}
{"level":"DEBUG","msg":"storage client compatibility configuration","provider":"gcs","request_checksum_calculation":"WHEN_REQUIRED","response_checksum_validation":"WHEN_REQUIRED","checksum_behavior_change":"suppress_unsolicited_trailers_for_gcs","use_path_style":true,"path_style_behavior_change":"enforced_path_style_for_provider","disable_log_output_checksum_validation_skipped":true,"auth_mode":"default_credential_chain","has_driver_overrides":false}
{"level":"INFO","msg":"storage client initialization complete","provider":"gcs","bucket":"my-bucket","region":"us-east-1","endpoint":"https://storage.googleapis.com","use_path_style":true,"checksum_mode":"when_required"}
```

**Key Differences:**
1. Embedded basic auth credentials (`admin:secret@`) are redacted from `endpoint` and `raw_endpoint`.
2. The `storage client initialization complete` record now consistently includes `endpoint`, `region`, `use_path_style`, and `checksum_mode`.
3. Added `DEBUG`-level `evaluating storage provider detection` log showing inspected host, matched rule, and normalized endpoint.
4. Added `DEBUG`-level `storage client compatibility configuration` log exposing request/response checksum calculation modes, trailer suppression, path-style behavior change, and credential auth mode.

---

## 2. Key Modules & Technical Design

### Deterministic Provider Detection (`internal/storage/detector`)
Endpoints and bucket URLs are evaluated in deterministic priority order (`Endpoint` takes precedence over `BucketURL`):

```go
type Detection struct {
    Provider contracts.Provider // contracts.ProviderAWS, ProviderGCS, or ProviderS3Compatible
    Host     string             // Normalized inspected hostname
    Rule     string             // Deterministic classification rule matched
    Bucket   string             // Extracted virtual-hosted bucket subdomain (if present)
}
```

- **`RuleDefaultAWS` (`"default_aws"`):** Triggered when both `Endpoint` and `BucketURL` are empty strings. Uses AWS SDK v2 default regional endpoints.
- **`RuleGCSHostname` (`"gcs_hostname"`):** Matches Google Cloud Storage XML API endpoints (`storage.googleapis.com`, `storage.<region>.rep.googleapis.com`, `storage[-<name>].p.googleapis.com`). Identifies GCS XML interoperability mode, path-style addressing, and unsolicited checksum trailer suppression.
- **`RuleAWSHostname` (`"aws_hostname"`):** Matches AWS S3 hostnames with S3 service labels (`s3`, `s3-*`, `s3express-*`) preceding AWS domain suffixes (`.amazonaws.com`, `.amazonaws.com.cn`, `.api.aws`). Uses standard AWS SDK virtual-hosted addressing.
- **`RuleS3CompatibleFallback` (`"s3_compatible_fallback"`):** Matched for custom endpoints (MinIO, LocalStack, Cloudflare R2, Wasabi, Ceph, custom internal hosts). Identifies path-style addressing and custom endpoint routing.
- **Subdomain Extraction:** `detector.BucketSubdomain(host)` extracts virtual-hosted bucket prefixes for multi-tenant providers (`mybucket.s3.us-east-1.amazonaws.com` -> `"mybucket"`, `mybucket.storage.googleapis.com` -> `"mybucket"`).

---

### URL Sanitization & Credential Protection (`internal/storage/urlutil`)
Implemented in `internal/storage/urlutil` and wrapped by `internal/storage/client/url_sanitize.go`:
- **Userinfo Redaction:** Strips basic auth credentials (`http://user:secret@host:9000` -> `http://host:9000`).
- **Query & Fragment Stripping:** Strips query parameters and fragments (`https://host/bucket?token=xyz#frag` -> `https://host/bucket`).
- **Schemeless URL Support:** Safely scrubs standard and complex schemeless URLs (e.g. `user:pass@host:9000/x?r=a://b` -> `host:9000/x`) without injecting default schemes.
- **Fail-Closed Malformed Redaction:** Protocol-relative URLs (`//admin:pass@host`), single-slash schemes (`https:/user:pass@host`), empty-authority strings (`https:////user:pass@host`), and path pseudo-schemes (`a/b://user:pass@host`) fail closed to `<unparseable>`.
- **Sanitized Error Messages:** Validation and detector errors omit raw input containing userinfo, query strings, fragments, or non-numeric port strings.

---

### Modular Connection Pipeline & Diagnostics (`internal/storage/client`)
Storage client initialization coordinates connection handling and observability across specialized components:
- **`connect.go`:** Top-level coordinator for `storage.Connect` and `storage.ConnectWithConfig`.
- **`connect_validate.go`:** Configuration validation and lifecycle progress logging (`initializing storage client`, `storage client initialization complete`).
- **`connect_log_diagnostic.go`:** Pre-connection detection diagnostics (`evaluating storage provider detection` at `DEBUG`).
- **`connect_log_compat.go` & `connect_log_compat_format.go`:** S3 compatibility diagnostic evaluation (`storage client compatibility configuration` at `DEBUG`).
- **`connect_probe.go`:** Startup probe gating (`maybeProbeStartup`, `probeStartupClient`) and centralized abnormal exit handling via `central.ExitWithAbnormal` (`handleConnectionError`, `SetExitFunc`).
- **`probe.go`:** Bucket connectivity and reachability probing (`Ping`, `probeBucket`, executing `HeadBucket` with fallback to `ListObjectsV2` on `AccessDenied`).
- **`connect_client.go`:** AWS SDK v2 client construction and option binding.

---

## 3. Public API & Backward Compatibility

- **Public API Surface:** Zero breaking changes to `pkg/storage` public signatures, interfaces, or option builders (`storage.Connect`, `storage.ConnectWithConfig`, `storage.Config`, `storage.WithSkipProbe`, `storage.WithDriverOptions`, etc.).
- **Contracts:** `contracts.Provider` enum constants (`ProviderAWS`, `ProviderGCS`, `ProviderS3Compatible`) remain fully backward-compatible.
- **Detector Shim:** `detector.DetectProvider(endpoint, bucketURL)` is preserved for legacy callers.

---

## 4. Behavior Caveats & Known Limitations

1. **Malformed URL Rejection:** `detector.ParseTarget`, `detector.Detect`, and `config.Config.Validate` reject malformed shapes (single-slash schemes like `https:/host`, bare `http`/`https` hosts, colon-terminated hosts, and nested `scheme://userinfo@host` inside paths). Previously on `main`, some of these shapes were accepted and logged raw at `INFO`. Legitimate endpoints (including IPv6 literals like `[::1]:9000` and paths with `@` or `%40`) remain fully supported.
2. **Log Sanitizer on Path `@`:** URLs containing literal `@` characters inside path segments (e.g. `https://proxy.example.com/s3/a@b/`) are accepted by detection and client validation, but sanitize to `<unparseable>` in log output as a defense-in-depth measure against ambiguous credential placement.
3. **Path-Based Secrets:** Secrets or tokens embedded directly inside URL path segments outside standard userinfo and query parameters are not recognized by standard URL parsers and are not redacted.
4. **Unsanitized Error Texts (Probe & Path Mismatch):** Error messages returned from underlying HTTP/network transports during bucket reachability probing (and logged on abnormal exit records via `central.ExitWithAbnormal`) are not passed through URL sanitization. Similarly, multi-tenant path-mismatch errors (`checkMultiTenantPathMismatch` in `options_endpoint_bucket_norm.go`) quote the configured bucket name and URL path (`u.Path`) verbatim.

---

## 5. Verification & Test Suite Summary

All test suites pass cleanly across all packages with zero static analysis issues and race detection enabled:

- **Static Analysis:** `go vet ./...` reported **0 issues**.
- **Race Detection:** `go test -race -count=1 ./...` passed across all 31 testable packages (`internal/json_schema` contains no test files).
- **Coverage Areas:**
  - Multi-cloud provider classification matrices (`detector_aws_test.go`, `detector_gcs_test.go`, `detector_custom_test.go`, `detector_detect_test.go`).
  - S3 compatibility diagnostics matrices (`connect_debug_compat_matrix_test.go`, `connect_debug_compat_test.go`, `connect_debug_driver_test.go`).
  - Detection debug logging suites (`connect_debug_detection_test.go`, `connect_debug_detection_nonaws_test.go`).
  - Credential and security redaction suites (`url_sanitize_test.go`, `urlutil_test.go`, `connect_debug_security_test.go`, `detector_invalid_test.go`).
  - Network and configuration validation suites (`config_validate_net_test.go`, `config_validate_bucket_test.go`, `config_defaults_test.go`).
  - Connection routing and schemeless endpoint suites (`connect_routing_test.go`, `connect_routing_schemeless_test.go`, `connect_routing_aws_test.go`, `connect_routing_custom_test.go`).
