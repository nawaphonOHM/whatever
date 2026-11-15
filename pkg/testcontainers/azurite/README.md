# Azurite Testcontainers (`pkg/testcontainers/azurite`)

The `pkg/testcontainers/azurite` package provides a reusable Testcontainers wrapper for starting the Azurite Azure Storage emulator in Go integration tests. It uses the official Testcontainers Azure module and exposes mapped Blob, Queue, and Table service endpoints.

Azurite implements Azure Storage APIs, not the AWS S3 API. An AWS S3 client cannot connect directly to Azurite or Azure Blob. The tagged integration tests use a separate S3Proxy gateway to translate S3 requests to Azure Blob requests; the public helper API and its Azure service URLs remain unchanged.

## Start and configure the emulator

`Run(ctx, opts...)` starts the pinned Azurite image (`DefaultImage`) with Blob, Queue, and Table services enabled. The functional options are `WithImage`, `WithEnv`, and `WithContainerOptions`; `DefaultOptions` and `NewOptions` are also available for inspecting configuration.

```go
ctx := context.Background()
container, err := azurite.Run(ctx, azurite.WithEnv("EXAMPLE_SETTING", "value"))
if err != nil {
	return err
}
defer container.Terminate(ctx)

blobEndpoint, err := container.BlobServiceURL(ctx)
if err != nil {
	return err
}
```

`BlobServiceURL`, `QueueServiceURL`, and `TableServiceURL` return service endpoints using Testcontainers' mapped host ports; they do not include an assumed localhost or fixed host port. `RawContainer` exposes the upstream Azurite container, and `Terminate` stops and removes it. Service URL calls return `ErrNilContainer` for a nil receiver and `ErrContainerUnavailable` when the container is stopped or terminated.

## S3 gateway integration fixture

The integration fixture starts two pinned Docker images on a private Testcontainers network:

- `mcr.microsoft.com/azure-storage/azurite:3.37.0` runs the Azure Blob emulator.
- `andrewgaul/s3proxy:3.2.0` runs S3Proxy with its `azureblob-sdk` backend.

The request path is `AWS S3 SDK or pkg/storage -> S3Proxy -> Azurite Blob`. S3Proxy uses the internal endpoint `http://azurite-blob:10000/devstoreaccount1`, the Azurite default account name and key (`DefaultAccountName` and `DefaultAccountKey`), and is configured to accept AWS Signature Version 2 or 4. The test client signs with SigV4 and uses path-style requests in region `us-east-1`, with dummy S3 credentials `azurite-test-access-key` / `azurite-test-secret-key`; these credentials are local test values only. The Go test client connects to S3Proxy's dynamically mapped host port, never to Azurite's host-mapped port.

The gateway credentials and endpoints are for the isolated local test network only. Do not use them with production storage, and do not interpret these tests as support for direct S3 connections to Azure Blob or Azurite. The S3Proxy backend is the protocol translation layer.

## Integration test

The `testcontainers`-tagged suite verifies authenticated bucket creation, `pkg/storage` object upload and byte-for-byte readback, object and bucket deletion, and failure when Azurite is unavailable. It requires Docker (or a compatible runtime) and the two images above; it does not require live AWS credentials or an Azure account:

```sh
go test -tags=testcontainers -race ./pkg/testcontainers/azurite/...
```
