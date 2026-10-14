// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Option defines a functional option for customizing internal storage client options.
type Option func(*Options)

// Options holds internal driver configuration options.
type Options struct {
	HTTPClient     *http.Client
	DriverOptions  []func(*s3.Options)
	PresignOptions []func(*s3.PresignOptions)
	SkipProbe      bool
}

// NewOptions creates an Options struct by evaluating functional options.
func NewOptions(opts ...Option) *Options {
	o := &Options{}
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	return o
}

// WithDriverOptions appends raw AWS SDK s3.Options modifiers internally.
func WithDriverOptions(opts ...func(*s3.Options)) Option {
	return func(o *Options) {
		for _, opt := range opts {
			if opt != nil {
				o.DriverOptions = append(o.DriverOptions, opt)
			}
		}
	}
}

// WithPresignOptions appends raw AWS SDK PresignOptions modifiers internally.
func WithPresignOptions(opts ...func(*s3.PresignOptions)) Option {
	return func(o *Options) {
		for _, opt := range opts {
			if opt != nil {
				o.PresignOptions = append(o.PresignOptions, opt)
			}
		}
	}
}

// WithHTTPClient overrides the HTTP client used by the AWS SDK client.
func WithHTTPClient(client *http.Client) Option {
	return func(o *Options) {
		o.HTTPClient = client
	}
}

// WithSkipProbe controls whether startup connectivity probing is skipped.
func WithSkipProbe(skip bool) Option {
	return func(o *Options) {
		o.SkipProbe = skip
	}
}

func applyDriverFn(opts *s3.Options, fn func(*s3.Options)) {
	if fn != nil {
		fn(opts)
	}
}

func (o *Options) applyDriver(opts *s3.Options) {
	if o == nil {
		return
	}
	for _, fn := range o.DriverOptions {
		applyDriverFn(opts, fn)
	}
}
