//go:build testcontainers

package azurite_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"

	"github.com/nawaphonOHM/whatever/v3/pkg/testcontainers/azurite"
)

const (
	s3ProxyImage          = "andrewgaul/s3proxy:3.2.0"
	s3ProxyPort           = "80/tcp"
	s3ProxyIdentity       = "azurite-test-access-key"
	s3ProxyCredential     = "azurite-test-secret-key"
	s3ProxyRegion         = "us-east-1"
	azuriteNetworkAlias   = "azurite-blob"
	s3GatewayCleanupLimit = 30 * time.Second
	s3GatewayRequestLimit = 20 * time.Second
)

func startS3GatewayFixture(
	t *testing.T,
	ctx context.Context,
) (string, *s3.Client, func(context.Context) error) {
	t.Helper()
	tcNetwork := startS3TestNetwork(t, ctx)
	emulator, markEmulatorRemoved := startS3Azurite(t, ctx, tcNetwork)
	gateway := startS3Proxy(t, ctx, tcNetwork.Name)
	endpoint := s3ProxyEndpoint(t, ctx, gateway)
	return endpoint, newS3GatewayClient(endpoint), azuriteStopper(emulator, markEmulatorRemoved)
}

func startS3TestNetwork(t *testing.T, ctx context.Context) *testcontainers.DockerNetwork {
	t.Helper()
	tcNetwork, err := network.New(ctx)
	registerS3NetworkCleanup(t, tcNetwork)
	require.NoError(t, err)
	require.NotNil(t, tcNetwork)
	return tcNetwork
}

func startS3Azurite(
	t *testing.T,
	ctx context.Context,
	tcNetwork *testcontainers.DockerNetwork,
) (*azurite.Container, func()) {
	t.Helper()
	emulator, err := azurite.Run(ctx, azurite.WithContainerOptions(
		network.WithNetworkName([]string{azuriteNetworkAlias}, tcNetwork.Name),
	))
	markEmulatorRemoved := registerEmulatorCleanup(t, emulator)
	require.NoError(t, err)
	require.NotNil(t, emulator)
	return emulator, markEmulatorRemoved
}
