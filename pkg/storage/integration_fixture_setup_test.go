//go:build testcontainers

package storage_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/testcontainers/testcontainers-go/modules/localstack"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const (
	localStackImage  = "localstack/localstack:4.14.0"
	localStackPort   = "4566/tcp"
	localStackBucket = "whatever-integration-bucket"
)

func setupSharedLocalStack() error {
	ctx := context.Background()
	container, err := startLocalStackContainer(ctx)
	if err != nil {
		return err
	}
	sharedLocalStack = container
	return connectSharedStorage(ctx, container)
}

func startLocalStackContainer(ctx context.Context) (*localstack.LocalStackContainer, error) {
	container, err := localstack.Run(ctx, localStackImage)
	registerLocalStackCleanup(container)
	if err != nil {
		return nil, fmt.Errorf("start LocalStack container: %w", err)
	}
	if container == nil {
		return nil, errors.New("start LocalStack container: no container returned")
	}
	return container, nil
}

func registerLocalStackCleanup(container *localstack.LocalStackContainer) {
	if container != nil {
		sharedStorageCleanup = localStackCleanup(container)
	}
}

func localStackCleanup(container *localstack.LocalStackContainer) func() {
	return func() {
		if err := container.Terminate(context.Background()); err != nil {
			log.Printf("failed to terminate LocalStack container: %v", err)
		}
	}
}

func connectSharedStorage(ctx context.Context, container *localstack.LocalStackContainer) error {
	cfg, err := localStackStorageConfig(ctx, container)
	if err != nil {
		return err
	}
	err = bootstrapLocalStackBucket(ctx, cfg)
	if err != nil {
		return err
	}
	return connectProbedSharedStorage(ctx, cfg)
}

func connectProbedSharedStorage(ctx context.Context, cfg *storage.Config) error {
	var err error
	sharedStorage, err = storage.ConnectWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect storage client with startup probe: %w", err)
	}
	return nil
}

func bootstrapLocalStackBucket(ctx context.Context, cfg *storage.Config) error {
	client, err := storage.ConnectWithConfig(ctx, cfg, storage.WithSkipProbe(true))
	if err != nil {
		return fmt.Errorf("connect storage client for bucket bootstrap: %w", err)
	}
	_, err = client.RawClient().CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(localStackBucket)})
	if err != nil {
		return fmt.Errorf("create LocalStack test bucket: %w", err)
	}
	return nil
}
