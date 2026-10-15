package client_test

import (
	"context"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

type opCall func(context.Context, *client.Client) error

func callGet(ctx context.Context, c *client.Client) error {
	_, err := c.GetObject(ctx, contracts.GetObjectInput{Key: testObjectKey})
	return err
}

func callHead(ctx context.Context, c *client.Client) error {
	_, err := c.HeadObject(ctx, contracts.HeadObjectInput{Key: testObjectKey})
	return err
}

func callDelete(ctx context.Context, c *client.Client) error {
	_, err := c.DeleteObject(ctx, contracts.DeleteObjectInput{Key: testObjectKey})
	return err
}

func callList(ctx context.Context, c *client.Client) error {
	_, err := c.ListObjects(ctx, contracts.ListObjectsInput{})
	return err
}

func callExists(ctx context.Context, c *client.Client) error {
	_, err := c.Exists(ctx, "", testObjectKey)
	return err
}
