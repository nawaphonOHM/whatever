//go:build testcontainers

package azurite_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/nawaphonOHM/whatever/v3/pkg/testcontainers/azurite"
)

func startS3Proxy(t *testing.T, ctx context.Context, networkName string) testcontainers.Container {
	t.Helper()
	gateway, err := testcontainers.GenericContainer(ctx, s3ProxyContainerRequest(networkName))
	registerS3ProxyCleanup(t, gateway)
	require.NoError(t, err)
	require.NotNil(t, gateway)
	return gateway
}

func s3ProxyContainerRequest(networkName string) testcontainers.GenericContainerRequest {
	backendEndpoint := "http://" + azuriteNetworkAlias + ":10000/" + azurite.DefaultAccountName
	containerRequest := testcontainers.ContainerRequest{
		Image:        s3ProxyImage,
		ExposedPorts: []string{s3ProxyPort},
		Env: map[string]string{
			"S3PROXY_AUTHORIZATION": "aws-v2-or-v4",
			"S3PROXY_IDENTITY":      s3ProxyIdentity,
			"S3PROXY_CREDENTIAL":    s3ProxyCredential,
			"JCLOUDS_PROVIDER":      "azureblob-sdk",
			"JCLOUDS_ENDPOINT":      backendEndpoint,
			"JCLOUDS_IDENTITY":      azurite.DefaultAccountName,
			"JCLOUDS_CREDENTIAL":    azurite.DefaultAccountKey,
		},
		Networks:   []string{networkName},
		WaitingFor: wait.ForHTTP("/healthz").WithPort(s3ProxyPort).WithStartupTimeout(time.Minute),
	}
	return testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	}
}

func s3ProxyEndpoint(t *testing.T, ctx context.Context, gateway testcontainers.Container) string {
	t.Helper()
	host, err := gateway.Host(ctx)
	require.NoError(t, err)
	port, err := gateway.MappedPort(ctx, s3ProxyPort)
	require.NoError(t, err)
	return "http://" + net.JoinHostPort(host, port.Port())
}

func registerS3NetworkCleanup(t *testing.T, tcNetwork *testcontainers.DockerNetwork) {
	t.Helper()
	if tcNetwork == nil {
		return
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), s3GatewayCleanupLimit)
		defer cancel()
		if err := tcNetwork.Remove(ctx); err != nil {
			t.Errorf("terminate Azurite S3 test network: %v", err)
		}
	})
}

func registerS3ProxyCleanup(t *testing.T, gateway testcontainers.Container) {
	t.Helper()
	if gateway == nil {
		return
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), s3GatewayCleanupLimit)
		defer cancel()
		if err := gateway.Terminate(ctx); err != nil {
			t.Errorf("terminate S3Proxy gateway: %v", err)
		}
	})
}
