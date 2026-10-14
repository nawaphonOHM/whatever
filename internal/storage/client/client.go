// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const defaultPresignDuration = 15 * time.Minute

// Client wraps the AWS SDK v2 S3 client and provides unified object storage operations.
type Client struct {
	rawClient      *s3.Client
	presignClient  *s3.PresignClient
	defaultBucket  string
	provider       contracts.Provider
	presignExpiry  atomic.Int64
	requestTimeout time.Duration
	pingTimeout    time.Duration
}

// ClientParams contains initialization parameters for Client.
type ClientParams struct {
	Raw           *s3.Client
	Presign       *s3.PresignClient
	Provider      contracts.Provider
	DefaultBucket string
	PresignExpiry time.Duration
	// RequestTimeout bounds each storage request (0 disables the client-side bound).
	RequestTimeout time.Duration
	// PingTimeout bounds every connectivity probe (0 disables the client-side bound).
	PingTimeout time.Duration
}

// NewClient initializes a new managed storage Client instance.
func NewClient(params ClientParams) *Client {
	expiry := params.PresignExpiry
	if expiry <= 0 {
		expiry = defaultPresignDuration
	}
	c := &Client{
		rawClient:      params.Raw,
		presignClient:  params.Presign,
		provider:       params.Provider,
		defaultBucket:  params.DefaultBucket,
		requestTimeout: params.RequestTimeout,
		pingTimeout:    params.PingTimeout,
	}
	c.presignExpiry.Store(int64(expiry))
	return c
}

// SetPresignExpiry updates the default pre-signed URL expiration duration concurrently.
func (c *Client) SetPresignExpiry(d time.Duration) {
	if c != nil && d > 0 {
		c.presignExpiry.Store(int64(d))
	}
}

// Provider returns the cloud object storage provider type.
func (c *Client) Provider() contracts.Provider {
	if c == nil {
		return ""
	}
	return c.provider
}

// DefaultBucket returns the default bucket configured for this client.
func (c *Client) DefaultBucket() string {
	if c == nil {
		return ""
	}
	return c.defaultBucket
}

// RawClient returns the underlying AWS SDK v2 S3 client handle.
func (c *Client) RawClient() *s3.Client {
	if c == nil {
		return nil
	}
	return c.rawClient
}

// PresignClient returns the underlying AWS SDK v2 presign client handle.
func (c *Client) PresignClient() *s3.PresignClient {
	if c == nil {
		return nil
	}
	return c.presignClient
}
