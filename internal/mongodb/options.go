package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Option defines a functional option for customizing internal MongoDB client
// options.
type Option func(*Options)

// Options holds internal driver configuration options.
type Options struct {
	DriverOptions []*options.ClientOptions
}

// NewOptions creates an Options struct by evaluating internal functional
// options.
func NewOptions(opts ...Option) *Options {
	o := &Options{}
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	return o
}

// WithDriverOptions appends raw mongo-driver ClientOptions internally.
func WithDriverOptions(opts ...*options.ClientOptions) Option {
	return func(o *Options) {
		for _, opt := range opts {
			if opt != nil {
				o.DriverOptions = append(o.DriverOptions, opt)
			}
		}
	}
}
