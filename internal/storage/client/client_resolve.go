// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const errExpiryFmt = "%w: %v (allowed: >0 and <= %v)"

func (c *Client) isNil() bool {
	return c == nil || c.rawClient == nil
}

func (c *Client) isPresignNil() bool {
	return c == nil || c.rawClient == nil || c.presignClient == nil
}

func (c *Client) resolveBucket(bucket string) string {
	if bucket != "" {
		return bucket
	}
	return c.defaultBucket
}

func (c *Client) defaultExpiry() time.Duration {
	if d := time.Duration(c.presignExpiry.Load()); d > 0 {
		return d
	}
	return defaultPresignDuration
}

func validateExpiry(expiry time.Duration) error {
	if expiry < 0 || expiry > config.MaxPresignExpiry {
		return fmt.Errorf(errExpiryFmt, contracts.ErrInvalidExpiry, expiry, config.MaxPresignExpiry)
	}
	return nil
}

func (c *Client) resolvePresignExpiry(d time.Duration) (time.Duration, error) {
	expiry := d
	if d == 0 {
		expiry = c.defaultExpiry()
	}
	return expiry, validateExpiry(expiry)
}

func withTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d > 0 {
		return context.WithTimeout(ctx, d)
	}
	return ctx, func() {}
}

func (c *Client) requestContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return withTimeout(ctx, c.requestTimeout)
}
