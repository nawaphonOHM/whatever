//go:build testcontainers

package storage_test

import (
	"context"
	"fmt"
	"net"

	"github.com/testcontainers/testcontainers-go/modules/localstack"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func localStackStorageConfig(ctx context.Context, container *localstack.LocalStackContainer) (*storage.Config, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("get LocalStack host: %w", err)
	}
	port, err := container.MappedPort(ctx, localStackPort)
	if err != nil {
		return nil, fmt.Errorf("get LocalStack mapped port: %w", err)
	}
	return localStackConfigForEndpoint(host, port.Port()), nil
}

func localStackConfigForEndpoint(host, port string) *storage.Config {
	cfg := storage.DefaultConfig()
	cfg.Bucket = localStackBucket
	cfg.Region = "us-east-1"
	cfg.AccessKeyID = "test"
	cfg.SecretAccessKey = "test"
	cfg.Endpoint = "http://" + net.JoinHostPort(host, port)
	cfg.UsePathStyle = true
	cfg.Insecure = true
	return cfg
}
